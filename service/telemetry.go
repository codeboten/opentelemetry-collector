// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package service // import "go.opentelemetry.io/collector/service"

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	ocmetric "go.opencensus.io/metric"
	"go.opencensus.io/metric/metricproducer"
	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/contrib/config"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/internal/obsreportconfig"
	"go.opentelemetry.io/collector/service/internal/proctelemetry"
	"go.opentelemetry.io/collector/service/internal/servicetelemetry"
	"go.opentelemetry.io/collector/service/telemetry"
)

const (
	zapKeyTelemetryAddress = "address"
	zapKeyTelemetryLevel   = "level"

	// supported trace propagators
	traceContextPropagator = "tracecontext"
	b3Propagator           = "b3"
)

var (
	errUnsupportedPropagator = errors.New("unsupported trace propagator")
)

type telemetryInitializer struct {
	views      []*view.View
	ocRegistry *ocmetric.Registry
	mp         metric.MeterProvider
	tp         trace.TracerProvider
	servers    []*http.Server

	useOtel                bool
	disableHighCardinality bool
	extendedConfig         bool
}

func newColTelemetry(useOtel bool, disableHighCardinality bool, extendedConfig bool) *telemetryInitializer {
	return &telemetryInitializer{
		mp:                     noopmetric.NewMeterProvider(),
		tp:                     nooptrace.NewTracerProvider(),
		useOtel:                useOtel,
		disableHighCardinality: disableHighCardinality,
		extendedConfig:         extendedConfig,
	}
}

func (tel *telemetryInitializer) init(res *resource.Resource, settings servicetelemetry.TelemetrySettings, cfg telemetry.Config, asyncErrorChannel chan error) error {
	if cfg.Metrics.Level == configtelemetry.LevelNone || (cfg.Metrics.Address == "" && len(cfg.Metrics.Readers) == 0) {
		settings.Logger.Info(
			"Skipping telemetry setup.",
			zap.String(zapKeyTelemetryAddress, cfg.Metrics.Address),
			zap.String(zapKeyTelemetryLevel, cfg.Metrics.Level.String()),
		)
		return nil
	}

	// TODO: move this to NewSDK
	if tp, err := textMapPropagatorFromConfig(cfg.Traces.Propagators); err == nil {
		otel.SetTextMapPropagator(tp)
	} else {
		return err
	}

	settings.Logger.Info("Setting up own telemetry...")
	if err := tel.translateLegacyConfig(&cfg, settings.Logger); err != nil {
		return err
	}
	otelSDK, err := config.NewSDK(config.WithOpenTelemetryConfiguration(config.OpenTelemetryConfiguration{
		TracerProvider: &config.TracerProvider{
			Processors: cfg.Traces.Processors,
		},
		// TODO: set propagators from config
		Propagator: cfg.Traces.Propagators,
	}))
	if err != nil {
		return err
	}
	tel.tp = otelSDK.TracerProvider()

	// Initialize the ocRegistry, still used by the process metrics.
	tel.ocRegistry = ocmetric.NewRegistry()
	metricproducer.GlobalManager().AddProducer(tel.ocRegistry)
	if !tel.useOtel && !tel.extendedConfig {
		return tel.initOpenCensus(res, settings.Logger, cfg.Metrics.Address, cfg.Metrics.Level, asyncErrorChannel)
	}

	tel.servers = append(tel.servers, proctelemetry.InitPrometheusServer(promRegistry, address, asyncErrorChannel))
	tel.mp = otelSDK.MeterProvider()

	return nil
}

func (tel *telemetryInitializer) translateLegacyConfig(cfg *telemetry.Config, logger *zap.Logger) error {
	if len(cfg.Metrics.Address) == 0 {
		return nil
	}

	if tel.extendedConfig {
		logger.Warn("service::telemetry::metrics::address is being deprecated in favor of service::telemetry::metrics::readers")
	}

	host, port, err := net.SplitHostPort(cfg.Metrics.Address)
	if err != nil {
		return err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if cfg.Metrics.Readers == nil {
		cfg.Metrics.Readers = []config.MetricReader{}
	}
	cfg.Metrics.Readers = append(cfg.Metrics.Readers, config.MetricReader{
		Pull: &config.PullMetricReader{
			Exporter: config.MetricExporter{
				Prometheus: &config.Prometheus{
					Host: &host,
					Port: &portInt,
				},
			},
		},
	})
	return nil
}

func (tel *telemetryInitializer) initOpenCensus(res *resource.Resource, logger *zap.Logger, address string, level configtelemetry.Level, asyncErrorChannel chan error) error {
	promRegistry := prometheus.NewRegistry()

	tel.views = obsreportconfig.AllViews(level)
	if err := view.Register(tel.views...); err != nil {
		return err
	}

	// Until we can use a generic metrics exporter, default to Prometheus.
	opts := ocprom.Options{
		Namespace: "otelcol",
		Registry:  promRegistry,
	}

	opts.ConstLabels = make(map[string]string)
	for _, keyValue := range res.Attributes() {
		opts.ConstLabels[sanitizePrometheusKey(string(keyValue.Key))] = keyValue.Value.AsString()
	}

	pe, err := ocprom.NewExporter(opts)
	if err != nil {
		return err
	}

	view.RegisterExporter(pe)

	logger.Info(
		"Serving Prometheus metrics",
		zap.String(zapKeyTelemetryAddress, address),
		zap.String(zapKeyTelemetryLevel, level.String()),
	)
	tel.servers = append(tel.servers, proctelemetry.InitPrometheusServer(promRegistry, address, asyncErrorChannel))
	return nil
}

func (tel *telemetryInitializer) shutdown() error {
	metricproducer.GlobalManager().DeleteProducer(tel.ocRegistry)
	view.Unregister(tel.views...)

	var errs error
	for _, server := range tel.servers {
		if server != nil {
			errs = multierr.Append(errs, server.Close())
		}
	}
	return errs
}

func sanitizePrometheusKey(str string) string {
	runeFilterMap := func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' {
			return r
		}
		return '_'
	}
	return strings.Map(runeFilterMap, str)
}

func textMapPropagatorFromConfig(props []string) (propagation.TextMapPropagator, error) {
	var textMapPropagators []propagation.TextMapPropagator
	for _, prop := range props {
		switch prop {
		case traceContextPropagator:
			textMapPropagators = append(textMapPropagators, propagation.TraceContext{})
		case b3Propagator:
			textMapPropagators = append(textMapPropagators, b3.New())
		default:
			return nil, errUnsupportedPropagator
		}
	}
	return propagation.NewCompositeTextMapPropagator(textMapPropagators...), nil
}

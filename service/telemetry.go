// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service // import "go.opentelemetry.io/collector/service"

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ocmetric "go.opencensus.io/metric"
	"go.opencensus.io/metric/metricproducer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"go.opentelemetry.io/collector/service/internal/proctelemetry"
	"go.opentelemetry.io/collector/service/telemetry"
)

const (
	zapKeyTelemetryAddress = "address"
	zapKeyTelemetryLevel   = "level"
)

type telemetryInitializer struct {
	ocRegistry *ocmetric.Registry
	mp         metric.MeterProvider
	server     []*http.Server
}

func newColTelemetry() *telemetryInitializer {
	return &telemetryInitializer{
		mp: metric.NewNoopMeterProvider(),
	}
}

func (tel *telemetryInitializer) initPeriodicReader(reader telemetry.Reader) (sdkmetric.Reader, error) {
	exporterType, ok := reader.Args["exporter"]
	if !ok {
		return nil, errors.New("no exporter configured")
	}
	exp, err := proctelemetry.InitExporter(exporterType)
	if err != nil {
		return nil, err
	}
	return sdkmetric.NewPeriodicReader(exp), nil

}

func (tel *telemetryInitializer) toReader(reader telemetry.Reader, logger *zap.Logger, cfg telemetry.Config, asyncErrorChannel chan error) (sdkmetric.Reader, error) {
	switch reader.Type {
	case "prometheus":
		return tel.initPrometheusReader(logger, reader.Args["endpoint"], cfg.Metrics.Level.String(), asyncErrorChannel)
	case "periodic":
		return tel.initPeriodicReader(reader)
	}
	return nil, fmt.Errorf("unsupported metric reader type: %s", reader.Type)
}

func (tel *telemetryInitializer) initPrometheusReader(logger *zap.Logger, address string, level string, asyncErrorChannel chan error) (sdkmetric.Reader, error) {
	promRegistry := prometheus.NewRegistry()
	wrappedRegisterer := prometheus.WrapRegistererWithPrefix("otelcol_", promRegistry)
	logger.Info(
		"Serving Prometheus metrics",
		zap.String(zapKeyTelemetryAddress, address),
		zap.String(zapKeyTelemetryLevel, level),
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	promServer := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	tel.server = append(tel.server, promServer)

	go func() {
		if serveErr := promServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			asyncErrorChannel <- serveErr
		}
	}()

	// We can remove `otelprom.WithoutUnits()` when the otel-go start exposing prometheus metrics using the OpenMetrics format
	// which includes metric units that prometheusreceiver uses to trim unit's suffixes from metric names.
	// https://github.com/open-telemetry/opentelemetry-go/issues/3468
	return otelprom.New(
		otelprom.WithRegisterer(wrappedRegisterer),
		otelprom.WithoutUnits(),
		// Disabled for the moment until this becomes stable, and we are ready to break backwards compatibility.
		otelprom.WithoutScopeInfo())
}

func (tel *telemetryInitializer) init(buildInfo component.BuildInfo, logger *zap.Logger, cfg telemetry.Config, asyncErrorChannel chan error) error {
	if cfg.Metrics.Level == configtelemetry.LevelNone || cfg.Metrics.Address == "" {
		logger.Info(
			"Skipping telemetry setup.",
			zap.String(zapKeyTelemetryAddress, cfg.Metrics.Address),
			zap.String(zapKeyTelemetryLevel, cfg.Metrics.Level.String()),
		)
		return nil
	}

	logger.Info("Setting up own telemetry...")

	// Construct telemetry attributes from build info and config's resource attributes.
	telAttrs := buildTelAttrs(buildInfo, cfg)

	if tp, err := proctelemetry.TextMapPropagatorFromConfig(cfg.Traces.Propagators); err == nil {
		otel.SetTextMapPropagator(tp)
	} else {
		return err
	}

	// Initialize the ocRegistry, still used by the process metrics.
	tel.ocRegistry = ocmetric.NewRegistry()
	metricproducer.GlobalManager().AddProducer(tel.ocRegistry)

	if len(cfg.Metrics.Address) > 0 {
		// TODO: improve message
		logger.Warn("service.telemetry.metrics.address is being deprecated in favour of service.telemetry.metrics.metric_readers")
		// TODO: account for multiple readers trying to use the same port
		cfg.Metrics.Readers = append(cfg.Metrics.Readers, telemetry.Reader{
			Type: "prometheus",
			Args: map[string]string{
				"endpoint": cfg.Metrics.Address,
			},
		})
	}
	readers := []sdkmetric.Option{}
	for _, reader := range cfg.Metrics.Readers {
		r, err := tel.toReader(reader, logger, cfg, asyncErrorChannel)
		if err != nil {
			return err
		}
		r.RegisterProducer(opencensus.NewMetricProducer())
		readers = append(readers, sdkmetric.WithReader(r))
	}
	var err error
	if tel.mp, err = proctelemetry.InitOpenTelemetry(telAttrs, readers); err != nil {
		return err
	}

	return nil
}

func buildTelAttrs(buildInfo component.BuildInfo, cfg telemetry.Config) map[string]string {
	telAttrs := map[string]string{}

	for k, v := range cfg.Resource {
		// nil value indicates that the attribute should not be included in the telemetry.
		if v != nil {
			telAttrs[k] = *v
		}
	}

	if _, ok := cfg.Resource[semconv.AttributeServiceName]; !ok {
		// AttributeServiceName is not specified in the config. Use the default service name.
		telAttrs[semconv.AttributeServiceName] = buildInfo.Command
	}

	if _, ok := cfg.Resource[semconv.AttributeServiceInstanceID]; !ok {
		// AttributeServiceInstanceID is not specified in the config. Auto-generate one.
		instanceUUID, _ := uuid.NewRandom()
		instanceID := instanceUUID.String()
		telAttrs[semconv.AttributeServiceInstanceID] = instanceID
	}

	if _, ok := cfg.Resource[semconv.AttributeServiceVersion]; !ok {
		// AttributeServiceVersion is not specified in the config. Use the actual
		// build version.
		telAttrs[semconv.AttributeServiceVersion] = buildInfo.Version
	}

	return telAttrs
}

func (tel *telemetryInitializer) shutdown() error {
	metricproducer.GlobalManager().DeleteProducer(tel.ocRegistry)
	var errs error
	for _, server := range tel.server {
		if server != nil {
			errs = multierr.Append(errs, server.Close())
		}
	}

	return errs
}

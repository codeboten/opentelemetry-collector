// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package proctelemetry // import "go.opentelemetry.io/collector/service/internal/proctelemetry"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/collector/obsreport"
	semconv "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.opentelemetry.io/collector/service/telemetry"
)

const (

	// gRPC Instrumentation Name
	GRPCInstrumentation = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	// http Instrumentation Name
	HTTPInstrumentation = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	// supported metric readers
	pullMetricReader     = "pull"
	periodicMetricReader = "periodic"

	// supported exporters
	prometheueExporter = "prometheus"
	otlpExporter       = "otlp"
	consoleExporter    = "console"

	// supported protocols
	protocolProtobufHTTP = "http/protobuf"
	protocolProtobufGRPC = "grpc/protobuf"
)

var (
	// GRPCUnacceptableKeyValues is a list of high cardinality grpc attributes that should be filtered out.
	GRPCUnacceptableKeyValues = []attribute.KeyValue{
		attribute.String(semconv.AttributeNetSockPeerAddr, ""),
		attribute.String(semconv.AttributeNetSockPeerPort, ""),
		attribute.String(semconv.AttributeNetSockPeerName, ""),
	}

	// HTTPUnacceptableKeyValues is a list of high cardinality http attributes that should be filtered out.
	HTTPUnacceptableKeyValues = []attribute.KeyValue{
		attribute.String(semconv.AttributeNetHostName, ""),
		attribute.String(semconv.AttributeNetHostPort, ""),
	}
)

func InitMetricReader(ctx context.Context, name string, reader any, asyncErrorChannel chan error) (sdkmetric.Reader, *http.Server, error) {
	readerType := strings.Split(name, "/")[0]
	switch readerType {
	case pullMetricReader:
		r, ok := reader.(telemetry.PullMetricReader)
		if !ok {
			return nil, nil, fmt.Errorf("invalid metric reader configuration: %v", reader)
		}
		return initPrometheusReader(ctx, r.Exporter, asyncErrorChannel)
	case periodicMetricReader:
		r, ok := reader.(telemetry.PeriodicMetricReader)
		if !ok {
			return nil, nil, fmt.Errorf("invalid metric reader configuration: %v", reader)
		}

		exporter, err := initExporter(ctx, r.Exporter, asyncErrorChannel)
		if err != nil {
			return nil, nil, err
		}
		return sdkmetric.NewPeriodicReader(exporter), nil, nil
	default:
		return nil, nil, fmt.Errorf("unsupported metric reader type: %s", readerType)
	}
}

func InitOpenTelemetry(res *resource.Resource, options []sdkmetric.Option, disableHighCardinality bool) (*sdkmetric.MeterProvider, error) {
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
		sdkmetric.WithView(batchViews(disableHighCardinality)...),
	}

	opts = append(opts, options...)
	return sdkmetric.NewMeterProvider(
		opts...,
	), nil
}

func InitPrometheusServer(registry *prometheus.Registry, address string, asyncErrorChannel chan error) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	go func() {
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			asyncErrorChannel <- serveErr
		}
	}()
	return server
}

func batchViews(disableHighCardinality bool) []sdkmetric.View {
	views := []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: obsreport.BuildProcessorCustomMetricName("batch", "batch_send_size")},
			sdkmetric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{10, 25, 50, 75, 100, 250, 500, 750, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 20000, 30000, 50000, 100000},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: obsreport.BuildProcessorCustomMetricName("batch", "batch_send_size_bytes")},
			sdkmetric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{10, 25, 50, 75, 100, 250, 500, 750, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 20000, 30000, 50000,
					100_000, 200_000, 300_000, 400_000, 500_000, 600_000, 700_000, 800_000, 900_000,
					1000_000, 2000_000, 3000_000, 4000_000, 5000_000, 6000_000, 7000_000, 8000_000, 9000_000},
			}},
		),
	}
	if disableHighCardinality {
		views = append(views, sdkmetric.NewView(sdkmetric.Instrument{
			Scope: instrumentation.Scope{
				Name: GRPCInstrumentation,
			},
		}, sdkmetric.Stream{
			AttributeFilter: cardinalityFilter(GRPCUnacceptableKeyValues...),
		}))
		views = append(views, sdkmetric.NewView(sdkmetric.Instrument{
			Scope: instrumentation.Scope{
				Name: HTTPInstrumentation,
			},
		}, sdkmetric.Stream{
			AttributeFilter: cardinalityFilter(HTTPUnacceptableKeyValues...),
		}))
	}
	return views
}

func cardinalityFilter(kvs ...attribute.KeyValue) attribute.Filter {
	filter := attribute.NewSet(kvs...)
	return func(kv attribute.KeyValue) bool {
		return !filter.HasValue(kv.Key)
	}
}

func initPrometheusExporter(prometheusConfig telemetry.Prometheus, asyncErrorChannel chan error) (sdkmetric.Reader, *http.Server, error) {
	promRegistry := prometheus.NewRegistry()
	if prometheusConfig.Host == nil {
		return nil, nil, fmt.Errorf("host must be specified")
	}
	if prometheusConfig.Port == nil {
		return nil, nil, fmt.Errorf("port must be specified")
	}
	wrappedRegisterer := prometheus.WrapRegistererWithPrefix("otelcol_", promRegistry)
	// We can remove `otelprom.WithoutUnits()` when the otel-go start exposing prometheus metrics using the OpenMetrics format
	// which includes metric units that prometheusreceiver uses to trim unit's suffixes from metric names.
	// https://github.com/open-telemetry/opentelemetry-go/issues/3468
	exporter, err := otelprom.New(
		otelprom.WithRegisterer(wrappedRegisterer),
		otelprom.WithoutUnits(),
		// Disabled for the moment until this becomes stable, and we are ready to break backwards compatibility.
		otelprom.WithoutScopeInfo())
	if err != nil {
		return nil, nil, fmt.Errorf("error creating otel prometheus exporter: %w", err)
	}

	exporter.RegisterProducer(opencensus.NewMetricProducer())
	return exporter, InitPrometheusServer(promRegistry, fmt.Sprintf("%s:%d", *prometheusConfig.Host, *prometheusConfig.Port), asyncErrorChannel), nil
}

func initOTLPgRPCExporter(ctx context.Context, otlpConfig telemetry.Otlp) (sdkmetric.Exporter, error) {
	opts := []otlpmetricgrpc.Option{}

	if len(otlpConfig.Endpoint) > 0 {
		opts = append(opts, otlpmetricgrpc.WithEndpoint(otlpConfig.Endpoint))
	}
	if otlpConfig.Compression != nil {
		opts = append(opts, otlpmetricgrpc.WithCompressor(*otlpConfig.Compression))
	}
	if otlpConfig.Timeout != nil {
		opts = append(opts, otlpmetricgrpc.WithTimeout(time.Millisecond*time.Duration(*otlpConfig.Timeout)))
	}
	if len(otlpConfig.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(toStringMap(otlpConfig.Headers)))
	}

	return otlpmetricgrpc.New(ctx, opts...)
}

func initOTLPHTTPExporter(ctx context.Context, otlpConfig telemetry.Otlp) (sdkmetric.Exporter, error) {
	opts := []otlpmetrichttp.Option{}

	if len(otlpConfig.Endpoint) > 0 {
		opts = append(opts, otlpmetrichttp.WithEndpoint(otlpConfig.Endpoint))
	}
	if otlpConfig.Compression != nil {
		opts = append(opts, otlpmetrichttp.WithCompressor(*otlpConfig.Compression))
	}
	if otlpConfig.Timeout != nil {
		opts = append(opts, otlpmetrichttp.WithTimeout(time.Millisecond*time.Duration(*otlpConfig.Timeout)))
	}
	if len(otlpConfig.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(toStringMap(otlpConfig.Headers)))
	}

	return otlpmetrichttp.New(ctx, opts...)
}

func initPrometheusReader(ctx context.Context, exporters telemetry.MetricExporter, asyncErrorChannel chan error) (sdkmetric.Reader, *http.Server, error) {
	for exporterType, exporter := range exporters {
		switch exporterType {
		case prometheueExporter:
			e, ok := exporter.(telemetry.Prometheus)
			if !ok {
				return nil, nil, fmt.Errorf("prometheus exporter invalid: %v", exporter)
			}
			return initPrometheusExporter(e, asyncErrorChannel)
		default:
			return nil, nil, fmt.Errorf("unsupported metric exporter type: %s", exporterType)
		}
	}
	return nil, nil, fmt.Errorf("no valid exporter: %v", exporters)
}

func initExporter(ctx context.Context, exporters telemetry.MetricExporter, asyncErrorChannel chan error) (sdkmetric.Exporter, error) {
	for exporterType, exporter := range exporters {
		switch exporterType {
		case consoleExporter:
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return stdoutmetric.New(
				stdoutmetric.WithEncoder(enc),
			)
		case otlpExporter:
			e, ok := exporter.(telemetry.Otlp)
			if !ok {
				return nil, fmt.Errorf("otlp exporter invalid: %v", exporter)
			}
			if e.Protocol == protocolProtobufHTTP {
				return initOTLPHTTPExporter(ctx, e)
			}
			return initOTLPgRPCExporter(ctx, e)
		default:
			return nil, fmt.Errorf("unsupported metric exporter type: %s", exporterType)
		}
	}
	return nil, fmt.Errorf("no valid exporter: %v", exporters)
}

func toStringMap(in map[string]interface{}) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = fmt.Sprint(v)
	}
	return out
}

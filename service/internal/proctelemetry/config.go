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

package proctelemetry // import "go.opentelemetry.io/collector/service/internal/proctelemetry"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/collector/obsreport"
)

const (
	// supported trace propagators
	traceContextPropagator = "tracecontext"
	b3Propagator           = "b3"

	// supported exporters
	stdoutmetricExporter   = "console"
	otlpmetricgrpcExporter = "otlp"

	// supported metric readers
	prometheusMetricReader = "prometheus"
	periodMetricReader     = "periodic"
)

var (
	errUnsupportedPropagator = errors.New("unsupported trace propagator")
)

func TextMapPropagatorFromConfig(props []string) (propagation.TextMapPropagator, error) {
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

func InitOpenTelemetry(attrs map[string]string, options []sdkmetric.Option) (*sdkmetric.MeterProvider, error) {
	var resAttrs []attribute.KeyValue
	for k, v := range attrs {
		resAttrs = append(resAttrs, attribute.String(k, v))
	}

	res, err := resource.New(context.Background(), resource.WithAttributes(resAttrs...))
	if err != nil {
		return nil, fmt.Errorf("error creating otel resources: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating otel prometheus exporter: %w", err)
	}
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
		sdkmetric.WithView(batchViews()...),
	}

	opts = append(opts, options...)
	return sdkmetric.NewMeterProvider(
		opts...,
	), nil
}

func InitExporter(exporterType string) (sdkmetric.Exporter, error) {
	switch exporterType {
	case stdoutmetricExporter:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return stdoutmetric.New(
			stdoutmetric.WithEncoder(enc),
		)
	case otlpmetricgrpcExporter:
		// TODO: pass context in
		ctx := context.Background()
		opts := []otlpmetricgrpc.Option{}
		// for k, v := range exporter.Args {
		// 	switch k {
		// 	case "endpoint":
		// 		opts = append(opts, otlpmetricgrpc.WithEndpoint(v))

		// 		// otlpmetricgrpc.WithAggregationSelector()
		// 		// otlpmetricgrpc.WithCompressor()
		// 		// otlpmetricgrpc.WithDialOption()
		// 		// otlpmetricgrpc.WithGRPCConn()
		// 		// otlpmetricgrpc.WithInsecure()
		// 		// otlpmetricgrpc.WithReconnectionPeriod()
		// 		// otlpmetricgrpc.WithRetry()
		// 		// otlpmetricgrpc.WithServiceConfig()
		// 		// otlpmetricgrpc.WithTLSCredentials()
		// 		// otlpmetricgrpc.WithTemporalitySelector()
		// 		// otlpmetricgrpc.WithTimeout()
		// 		// case "headers":
		// 		// 	opts = append(opts, otlpmetricgrpc.WithHeaders(v))
		// 	}
		// }
		return otlpmetricgrpc.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unsupported metric exporter type: %s", exporterType)
	}
}

func batchViews() []sdkmetric.View {
	return []sdkmetric.View{
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
}

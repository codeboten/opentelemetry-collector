// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package proctelemetry // import "go.opentelemetry.io/collector/service/internal/proctelemetry"

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"

	semconv "go.opentelemetry.io/collector/semconv/v1.18.0"
)

const (

	// gRPC Instrumentation Name
	GRPCInstrumentation = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	// http Instrumentation Name
	HTTPInstrumentation = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

// TODO: batch views should return view config
// func batchViews(disableHighCardinality bool) []sdkmetric.View {
// 	views := []sdkmetric.View{
// 		sdkmetric.NewView(
// 			sdkmetric.Instrument{Name: processorhelper.BuildCustomMetricName("batch", "batch_send_size")},
// 			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
// 				Boundaries: []float64{10, 25, 50, 75, 100, 250, 500, 750, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 20000, 30000, 50000, 100000},
// 			}},
// 		),
// 		sdkmetric.NewView(
// 			sdkmetric.Instrument{Name: processorhelper.BuildCustomMetricName("batch", "batch_send_size_bytes")},
// 			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
// 				Boundaries: []float64{10, 25, 50, 75, 100, 250, 500, 750, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 20000, 30000, 50000,
// 					100_000, 200_000, 300_000, 400_000, 500_000, 600_000, 700_000, 800_000, 900_000,
// 					1000_000, 2000_000, 3000_000, 4000_000, 5000_000, 6000_000, 7000_000, 8000_000, 9000_000},
// 			}},
// 		),
// 	}
// 	if disableHighCardinality {
// 		views = append(views, sdkmetric.NewView(sdkmetric.Instrument{
// 			Scope: instrumentation.Scope{
// 				Name: GRPCInstrumentation,
// 			},
// 		}, sdkmetric.Stream{
// 			AttributeFilter: cardinalityFilter(GRPCUnacceptableKeyValues...),
// 		}))
// 		views = append(views, sdkmetric.NewView(sdkmetric.Instrument{
// 			Scope: instrumentation.Scope{
// 				Name: HTTPInstrumentation,
// 			},
// 		}, sdkmetric.Stream{
// 			AttributeFilter: cardinalityFilter(HTTPUnacceptableKeyValues...),
// 		}))
// 	}
// 	return views
// }

// func cardinalityFilter(kvs ...attribute.KeyValue) attribute.Filter {
// 	filter := attribute.NewSet(kvs...)
// 	return func(kv attribute.KeyValue) bool {
// 		return !filter.HasValue(kv.Key)
// 	}
// }

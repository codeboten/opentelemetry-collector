// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry // import "go.opentelemetry.io/collector/service/telemetry"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

func ViewOptionsFromConfig(views []View) []sdkmetric.Option {
	opts := []sdkmetric.Option{}
	for _, view := range views {
		if view.Selector == nil || view.Stream == nil {
			continue
		}
		opts = append(opts, sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{
					Name: view.Selector.InstrumentNameStr(),
					Kind: instrumentTypeToKind(view.Selector.InstrumentTypeStr()),
					// TODO: add unit once https://github.com/open-telemetry/opentelemetry-configuration/pull/38
					//       is merged
					// Unit: *view.Selector.Unit,
					Scope: instrumentation.Scope{
						Name:      view.Selector.MeterNameStr(),
						Version:   view.Selector.MeterVersionStr(),
						SchemaURL: view.Selector.MeterSchemaUrlStr(),
					},
				},
				sdkmetric.Stream{
					Name:            view.Stream.NameStr(),
					Description:     view.Stream.DescriptionStr(),
					Aggregation:     viewStreamAggregationToAggregation(view.Stream.Aggregation),
					AttributeFilter: attributeKeysToAttributeFilter(view.Stream.AttributeKeys),
				},
			),
		))
	}
	return opts
}

var invalidInstrumentKind = sdkmetric.InstrumentKind(0)

func instrumentTypeToKind(instrument string) sdkmetric.InstrumentKind {
	switch instrument {
	case "counter":
		return sdkmetric.InstrumentKindCounter
	case "histogram":
		return sdkmetric.InstrumentKindHistogram
	case "observable_counter":
		return sdkmetric.InstrumentKindObservableCounter
	case "observable_gauge":
		return sdkmetric.InstrumentKindObservableGauge
	case "observable_updown_counter":
		return sdkmetric.InstrumentKindObservableUpDownCounter
	case "updown_counter":
		return sdkmetric.InstrumentKindUpDownCounter
	}
	return invalidInstrumentKind
}

func attributeKeysToAttributeFilter(keys []string) attribute.Filter {
	kvs := make([]attribute.KeyValue, len(keys))
	for i, key := range keys {
		kvs[i] = attribute.Bool(key, true)
	}
	filter := attribute.NewSet(kvs...)
	return func(kv attribute.KeyValue) bool {
		return !filter.HasValue(kv.Key)
	}
}

func viewStreamAggregationToAggregation(agg *ViewStreamAggregation) aggregation.Aggregation {
	if agg == nil {
		return aggregation.Default{}
	}
	if agg.Sum != nil {
		return aggregation.Sum{}
	}
	if agg.Drop != nil {
		return aggregation.Drop{}
	}
	if agg.LastValue != nil {
		return aggregation.LastValue{}
	}
	if agg.ExplicitBucketHistogram != nil {
		return aggregation.ExplicitBucketHistogram{
			Boundaries: agg.ExplicitBucketHistogram.Boundaries,
			NoMinMax:   !agg.ExplicitBucketHistogram.RecordMinMaxBool(),
		}
	}
	return aggregation.Default{}
}

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

package obsreport // import "go.opentelemetry.io/collector/obsreport"

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/internal/obsreportconfig/obsmetrics"
	"go.opentelemetry.io/collector/processor"
)

var (
	processorName  = "processor"
	processorScope = scopeName + nameSep + processorName
)

// BuildProcessorCustomMetricName is used to be build a metric name following
// the standards used in the Collector. The configType should be the same
// value used to identify the type on the config.
func BuildProcessorCustomMetricName(configType, metric string) string {
	componentPrefix := obsmetrics.ProcessorPrefix
	if !strings.HasSuffix(componentPrefix, obsmetrics.NameSep) {
		componentPrefix += obsmetrics.NameSep
	}
	if configType == "" {
		return componentPrefix
	}
	return componentPrefix + configType + obsmetrics.NameSep + metric
}

// Processor is a helper to add observability to a component.Processor.
type Processor struct {
	level configtelemetry.Level

	logger *zap.Logger

	otelAttrs []attribute.KeyValue

	acceptedSpansCounter        instrument.Int64Counter
	refusedSpansCounter         instrument.Int64Counter
	droppedSpansCounter         instrument.Int64Counter
	acceptedMetricPointsCounter instrument.Int64Counter
	refusedMetricPointsCounter  instrument.Int64Counter
	droppedMetricPointsCounter  instrument.Int64Counter
	acceptedLogRecordsCounter   instrument.Int64Counter
	refusedLogRecordsCounter    instrument.Int64Counter
	droppedLogRecordsCounter    instrument.Int64Counter
}

// ProcessorSettings are settings for creating a Processor.
type ProcessorSettings struct {
	ProcessorID             component.ID
	ProcessorCreateSettings processor.CreateSettings
}

// NewProcessor creates a new Processor.
func NewProcessor(cfg ProcessorSettings) (*Processor, error) {
	return newProcessor(cfg)
}

func newProcessor(cfg ProcessorSettings) (*Processor, error) {
	proc := &Processor{
		level:  cfg.ProcessorCreateSettings.MetricsLevel,
		logger: cfg.ProcessorCreateSettings.Logger,
		otelAttrs: []attribute.KeyValue{
			attribute.String(obsmetrics.ProcessorKey, cfg.ProcessorID.String()),
		},
	}

	if err := proc.createOtelMetrics(cfg); err != nil {
		return nil, err
	}

	return proc, nil
}

func (por *Processor) createOtelMetrics(cfg ProcessorSettings) error {
	meter := cfg.ProcessorCreateSettings.MeterProvider.Meter(processorScope)
	var errors, err error

	por.acceptedSpansCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.AcceptedSpansKey,
		instrument.WithDescription("Number of spans successfully pushed into the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.refusedSpansCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.RefusedSpansKey,
		instrument.WithDescription("Number of spans that were rejected by the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.droppedSpansCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.DroppedSpansKey,
		instrument.WithDescription("Number of spans that were dropped."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.acceptedMetricPointsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.AcceptedMetricPointsKey,
		instrument.WithDescription("Number of metric points successfully pushed into the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.refusedMetricPointsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.RefusedMetricPointsKey,
		instrument.WithDescription("Number of metric points that were rejected by the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.droppedMetricPointsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.DroppedMetricPointsKey,
		instrument.WithDescription("Number of metric points that were dropped."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.acceptedLogRecordsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.AcceptedLogRecordsKey,
		instrument.WithDescription("Number of log records successfully pushed into the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.refusedLogRecordsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.RefusedLogRecordsKey,
		instrument.WithDescription("Number of log records that were rejected by the next component in the pipeline."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	por.droppedLogRecordsCounter, err = meter.Int64Counter(
		obsmetrics.ProcessorPrefix+obsmetrics.DroppedLogRecordsKey,
		instrument.WithDescription("Number of log records that were dropped."),
		instrument.WithUnit("1"),
	)
	errors = multierr.Append(errors, err)

	return errors
}

func (por *Processor) recordData(ctx context.Context, dataType component.DataType, accepted, refused, dropped int64) {
	if por.level == configtelemetry.LevelNone {
		return
	}

	var acceptedCount, refusedCount, droppedCount instrument.Int64Counter
	switch dataType {
	case component.DataTypeTraces:
		acceptedCount = por.acceptedSpansCounter
		refusedCount = por.refusedSpansCounter
		droppedCount = por.droppedSpansCounter
	case component.DataTypeMetrics:
		acceptedCount = por.acceptedMetricPointsCounter
		refusedCount = por.refusedMetricPointsCounter
		droppedCount = por.droppedMetricPointsCounter
	case component.DataTypeLogs:
		acceptedCount = por.acceptedLogRecordsCounter
		refusedCount = por.refusedLogRecordsCounter
		droppedCount = por.droppedLogRecordsCounter
	}

	acceptedCount.Add(ctx, accepted, por.otelAttrs...)
	refusedCount.Add(ctx, refused, por.otelAttrs...)
	droppedCount.Add(ctx, dropped, por.otelAttrs...)
}

// TracesAccepted reports that the trace data was accepted.
func (por *Processor) TracesAccepted(ctx context.Context, numSpans int) {
	por.recordData(ctx, component.DataTypeTraces, int64(numSpans), int64(0), int64(0))
}

// TracesRefused reports that the trace data was refused.
func (por *Processor) TracesRefused(ctx context.Context, numSpans int) {
	por.recordData(ctx, component.DataTypeTraces, int64(0), int64(numSpans), int64(0))
}

// TracesDropped reports that the trace data was dropped.
func (por *Processor) TracesDropped(ctx context.Context, numSpans int) {
	por.recordData(ctx, component.DataTypeTraces, int64(0), int64(0), int64(numSpans))
}

// MetricsAccepted reports that the metrics were accepted.
func (por *Processor) MetricsAccepted(ctx context.Context, numPoints int) {
	por.recordData(ctx, component.DataTypeMetrics, int64(numPoints), int64(0), int64(0))
}

// MetricsRefused reports that the metrics were refused.
func (por *Processor) MetricsRefused(ctx context.Context, numPoints int) {
	por.recordData(ctx, component.DataTypeMetrics, int64(0), int64(numPoints), int64(0))
}

// MetricsDropped reports that the metrics were dropped.
func (por *Processor) MetricsDropped(ctx context.Context, numPoints int) {
	por.recordData(ctx, component.DataTypeMetrics, int64(0), int64(0), int64(numPoints))
}

// LogsAccepted reports that the logs were accepted.
func (por *Processor) LogsAccepted(ctx context.Context, numRecords int) {
	por.recordData(ctx, component.DataTypeLogs, int64(numRecords), int64(0), int64(0))
}

// LogsRefused reports that the logs were refused.
func (por *Processor) LogsRefused(ctx context.Context, numRecords int) {
	por.recordData(ctx, component.DataTypeLogs, int64(0), int64(numRecords), int64(0))
}

// LogsDropped reports that the logs were dropped.
func (por *Processor) LogsDropped(ctx context.Context, numRecords int) {
	por.recordData(ctx, component.DataTypeLogs, int64(0), int64(0), int64(numRecords))
}

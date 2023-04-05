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

package batchprocessor // import "go.opentelemetry.io/collector/processor/batchprocessor"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/internal/obsreportconfig/obsmetrics"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/processor"
)

const (
	scopeName = "go.opentelemetry.io/collector/processor/batchprocessor"
)

type trigger int

const (
	triggerTimeout trigger = iota
	triggerBatchSize
)

type batchProcessorTelemetry struct {
	level    configtelemetry.Level
	detailed bool

	exportCtx context.Context

	processorAttr        []attribute.KeyValue
	batchSizeTriggerSend instrument.Int64Counter
	timeoutTriggerSend   instrument.Int64Counter
	batchSendSize        instrument.Int64Histogram
	batchSendSizeBytes   instrument.Int64Histogram
}

func newBatchProcessorTelemetry(set processor.CreateSettings) (*batchProcessorTelemetry, error) {
	bpt := &batchProcessorTelemetry{
		processorAttr: []attribute.KeyValue{attribute.String(obsmetrics.ProcessorKey, set.ID.String())},
		exportCtx:     context.Background(),
		level:         set.MetricsLevel,
		detailed:      set.MetricsLevel == configtelemetry.LevelDetailed,
	}

	err := bpt.createOtelMetrics(set.MeterProvider)
	if err != nil {
		return nil, err
	}

	return bpt, nil
}

func (bpt *batchProcessorTelemetry) createOtelMetrics(mp metric.MeterProvider) error {
	var err error
	meter := mp.Meter(scopeName)

	bpt.batchSizeTriggerSend, err = meter.Int64Counter(
		obsreport.BuildProcessorCustomMetricName(typeStr, "batch_size_trigger_send"),
		instrument.WithDescription("Number of times the batch was sent due to a size trigger"),
		instrument.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	bpt.timeoutTriggerSend, err = meter.Int64Counter(
		obsreport.BuildProcessorCustomMetricName(typeStr, "timeout_trigger_send"),
		instrument.WithDescription("Number of times the batch was sent due to a timeout trigger"),
		instrument.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	bpt.batchSendSize, err = meter.Int64Histogram(
		obsreport.BuildProcessorCustomMetricName(typeStr, "batch_send_size"),
		instrument.WithDescription("Number of units in the batch"),
		instrument.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	bpt.batchSendSizeBytes, err = meter.Int64Histogram(
		obsreport.BuildProcessorCustomMetricName(typeStr, "batch_send_size_bytes"),
		instrument.WithDescription("Number of bytes in batch that was sent"),
		instrument.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (bpt *batchProcessorTelemetry) record(trigger trigger, sent, bytes int64) {
	switch trigger {
	case triggerBatchSize:
		bpt.batchSizeTriggerSend.Add(bpt.exportCtx, 1, bpt.processorAttr...)
	case triggerTimeout:
		bpt.timeoutTriggerSend.Add(bpt.exportCtx, 1, bpt.processorAttr...)
	}

	bpt.batchSendSize.Record(bpt.exportCtx, sent, bpt.processorAttr...)
	if bpt.detailed {
		bpt.batchSendSizeBytes.Record(bpt.exportCtx, bytes, bpt.processorAttr...)
	}
}

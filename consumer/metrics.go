// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package consumer // import "go.opentelemetry.io/collector/consumer"

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Metrics is the new metrics consumer interface that receives pmetric.Metrics, processes it
// as needed, and sends it to the next processing node if any or to the destination.
type Metrics interface {
	baseConsumer
	// ConsumeMetrics receives pmetric.Metrics for consumption.
	ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error
}

// ConsumeMetricsFunc is a helper function that is similar to ConsumeMetrics.
type ConsumeMetricsFunc func(ctx context.Context, md pmetric.Metrics) error

// ConsumeMetrics calls f(ctx, md).
func (f ConsumeMetricsFunc) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	return f(ctx, md)
}

type baseMetrics struct {
	*baseImpl
	ConsumeMetricsFunc
}

// NewMetrics returns a Metrics configured with the provided options.
func NewMetrics(consume ConsumeMetricsFunc, options ...Option) (Metrics, error) {
	if consume == nil {
		return nil, errNilFunc
	}

	baseImpl := newBaseImpl(options...)
	fn := func(ctx context.Context, ld pmetric.Metrics) error {
		ctx = baseImpl.obsreport.StartTracesOp(ctx)
		err := consume(ctx, ld)
		baseImpl.obsreport.EndTracesOp(ctx, ld.DataPointCount(), err)
		return err
	}

	return &baseMetrics{
		baseImpl:           baseImpl,
		ConsumeMetricsFunc: fn,
	}, nil
}

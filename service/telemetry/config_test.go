// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/confmap"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		success bool
	}{
		{
			name: "basic metric telemetry",
			cfg: &Config{
				Metrics: MetricsConfig{
					Level:   configtelemetry.LevelBasic,
					Address: "127.0.0.1:3333",
				},
			},
			success: true,
		},
		{
			name: "invalid metric telemetry",
			cfg: &Config{
				Metrics: MetricsConfig{
					Level:   configtelemetry.LevelBasic,
					Address: "",
				},
			},
			success: false,
		},
		{
			name: "valid metric telemetry with metric readers",
			cfg: &Config{
				Metrics: MetricsConfig{
					Level: configtelemetry.LevelBasic,
					Readers: MeterProviderJsonMetricReaders{
						"pull/prometheus": PullMetricReader{},
					},
				},
			},
			success: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.success {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUnmarshalMetricReaders(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *confmap.Conf
		success bool
		err     string
	}{
		{
			name: "invalid config",
			cfg:  confmap.NewFromStringMap(map[string]any{"invalid": "invalid"}),
			err:  "unsupported metric reader type: invalid",
		},
	}
	for _, tt := range tests {
		reader := make(MeterProviderJsonMetricReaders)
		err := reader.Unmarshal(tt.cfg)
		assert.ErrorContains(t, err, tt.err)
	}
}

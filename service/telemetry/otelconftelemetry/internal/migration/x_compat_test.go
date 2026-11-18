// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xotelconf "go.opentelemetry.io/contrib/otelconf/x"

	"go.opentelemetry.io/collector/confmap"
)

// These tests characterize how the otelconf/x types handle the configuration
// shapes accepted today (otelconf/v0.3.0 and v0.2.0). They do not change any
// behavior. They document why the collector needs to rewrite legacy keys
// before unmarshaling into the x types, and they pin the x keys that the
// rewrite must target.

// unmarshalFails reports whether unmarshaling into out returns an error or panics.
// The x union types carry a `mapstructure:",remain"` field of type interface{},
// which makes mapstructure panic ("reflect: Key of non-map type interface {}")
// when it sees a key the type does not know.
func unmarshalFails(m map[string]any, out any) (failed bool) {
	defer func() {
		if r := recover(); r != nil {
			failed = true
		}
	}()
	return confmap.NewFromStringMap(m).Unmarshal(out) != nil
}

func pullPrometheus(key string) map[string]any {
	return map[string]any{"readers": []any{map[string]any{"pull": map[string]any{
		"exporter": map[string]any{key: map[string]any{"host": "127.0.0.1", "port": 8888}},
	}}}}
}

func batchOTLP(key string) map[string]any {
	return map[string]any{"processors": []any{map[string]any{"batch": map[string]any{
		"exporter": map[string]any{key: map[string]any{"endpoint": "http://127.0.0.1:4318"}},
	}}}}
}

// TestXRejectsLegacyKeys: the v0.3.0 keys must be rewritten before reaching the x types.
func TestXRejectsLegacyKeys(t *testing.T) {
	t.Run("metrics prometheus", func(t *testing.T) {
		var mp xotelconf.MeterProvider
		assert.True(t, unmarshalFails(pullPrometheus("prometheus"), &mp))
	})
	t.Run("traces otlp", func(t *testing.T) {
		var tp xotelconf.TracerProvider
		assert.True(t, unmarshalFails(batchOTLP("otlp"), &tp))
	})
	t.Run("logs otlp", func(t *testing.T) {
		var lp xotelconf.LoggerProvider
		assert.True(t, unmarshalFails(batchOTLP("otlp"), &lp))
	})
}

// TestXAcceptsNewKeys: the keys the legacy ones must be rewritten to.
func TestXAcceptsNewKeys(t *testing.T) {
	t.Run("metrics prometheus/development", func(t *testing.T) {
		var mp xotelconf.MeterProvider
		require.False(t, unmarshalFails(pullPrometheus("prometheus/development"), &mp))
		require.Len(t, mp.Readers, 1)
		prom := mp.Readers[0].Pull.Exporter.PrometheusDevelopment
		require.NotNil(t, prom)
		assert.Equal(t, "127.0.0.1", *prom.Host)
		assert.Equal(t, 8888, *prom.Port)
	})
	t.Run("traces otlp_http", func(t *testing.T) {
		var tp xotelconf.TracerProvider
		require.False(t, unmarshalFails(batchOTLP("otlp_http"), &tp))
		require.Len(t, tp.Processors, 1)
		require.NotNil(t, tp.Processors[0].Batch.Exporter.OTLPHttp)
	})
	t.Run("logs otlp_grpc", func(t *testing.T) {
		var lp xotelconf.LoggerProvider
		require.False(t, unmarshalFails(batchOTLP("otlp_grpc"), &lp))
		require.Len(t, lp.Processors, 1)
		require.NotNil(t, lp.Processors[0].Batch.Exporter.OTLPGrpc)
	})
}

// Field-level differences between the v0.3.0 `prometheus` exporter and x's
// `prometheus/development`, which the key rewrite has to translate:
//
//	v0.3.0                                 x (prometheus/development)
//	host, port                             host, port                     (same)
//	with_resource_constant_labels          with_resource_constant_labels  (same)
//	without_scope_info                     without_scope_info             (same)
//	without_units + without_type_suffix    translation_strategy           (merged)
//	-                                      without_target_info            (new, not applied by otelconf/x v0.27.0)
//
// otelconf/v0.3.0 maps both without_* flags true to no_translation, and
// anything else to no_utf8_escaping_with_suffixes.

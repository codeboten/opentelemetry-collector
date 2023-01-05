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

package telemetry // import "go.opentelemetry.io/collector/service/telemetry"

import (
	"fmt"

	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/config/configtelemetry"
)

// Config defines the configurable settings for service telemetry.
type Config struct {
	Resource map[string]*string `mapstructure:"resource"`
	Metrics  MetricsConfig      `mapstructure:"metrics"`
	Traces   TracesConfig       `mapstructure:"traces"`
	Logs     LogsConfig         `mapstructure:"logs"`
}

// LogsConfig defines the configurable settings for service telemetry logs.
// This MUST be compatible with zap.Config. Cannot use directly zap.Config because
// the collector uses mapstructure and not yaml tags.
type LogsConfig struct {
	Sampling          *LogsSamplingConfig    `mapstructure:"sampling"`
	InitialFields     map[string]interface{} `mapstructure:"initial_fields"`
	Encoding          string                 `mapstructure:"encoding"`
	OutputPaths       []string               `mapstructure:"output_paths"`
	ErrorOutputPaths  []string               `mapstructure:"error_output_paths"`
	Level             zapcore.Level          `mapstructure:"level"`
	Development       bool                   `mapstructure:"development"`
	DisableCaller     bool                   `mapstructure:"disable_caller"`
	DisableStacktrace bool                   `mapstructure:"disable_stacktrace"`
}

// LogsSamplingConfig sets a sampling strategy for the logger. Sampling caps the
// global CPU and I/O load that logging puts on your process while attempting
// to preserve a representative subset of your logs.
type LogsSamplingConfig struct {
	Initial    int `mapstructure:"initial"`
	Thereafter int `mapstructure:"thereafter"`
}

// MetricsConfig exposes the common Telemetry configuration for one component.
// Experimental: *NOTE* this structure is subject to change or removal in the future.
type MetricsConfig struct {
	Address string                `mapstructure:"address"`
	Level   configtelemetry.Level `mapstructure:"level"`
}

// TracesConfig exposes the common Telemetry configuration for collector's internal spans.
// Experimental: *NOTE* this structure is subject to change or removal in the future.
type TracesConfig struct {
	// Propagators is a list of TextMapPropagators from the supported propagators list. Currently,
	// tracecontext and  b3 are supported. By default, the value is set to empty list and
	// context propagation is disabled.
	Propagators []string `mapstructure:"propagators"`
}

// Validate checks whether the current configuration is valid
func (c *Config) Validate() error {

	// Check when service telemetry metric level is not none, the metrics address should not be empty
	if c.Metrics.Level != configtelemetry.LevelNone && c.Metrics.Address == "" {
		return fmt.Errorf("collector telemetry metric address should exist when metric level is not none")
	}

	return nil
}

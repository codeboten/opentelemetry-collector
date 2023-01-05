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

package otlphttpexporter // import "go.opentelemetry.io/collector/exporter/otlphttpexporter"

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for OTLP/HTTP exporter.
type Config struct {
	TracesEndpoint                string `mapstructure:"traces_endpoint"`
	MetricsEndpoint               string `mapstructure:"metrics_endpoint"`
	LogsEndpoint                  string `mapstructure:"logs_endpoint"`
	exporterhelper.QueueSettings  `mapstructure:"sending_queue"`
	confighttp.HTTPClientSettings `mapstructure:",squash"`
	exporterhelper.RetrySettings  `mapstructure:"retry_on_failure"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" && cfg.TracesEndpoint == "" && cfg.MetricsEndpoint == "" && cfg.LogsEndpoint == "" {
		return errors.New("at least one endpoint must be specified")
	}
	return nil
}

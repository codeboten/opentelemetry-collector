// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package telemetry

type Attributes struct {
	// ServiceName corresponds to the JSON schema field "service.name".
	ServiceName *string `json:"service.name,omitempty" yaml:"service.name,omitempty" mapstructure:"service.name,omitempty"`
}

type LoggerProviderJson struct {
	// LogrecordLimits corresponds to the JSON schema field "logrecord_limits".
	LogrecordLimits *LoggerProviderJsonLogrecordLimits `json:"logrecord_limits,omitempty" yaml:"logrecord_limits,omitempty" mapstructure:"logrecord_limits,omitempty"`
}

type LoggerProviderJsonLogrecordLimits struct {
	// AttributeCountLimit corresponds to the JSON schema field
	// "attribute_count_limit".
	AttributeCountLimit *int `json:"attribute_count_limit,omitempty" yaml:"attribute_count_limit,omitempty" mapstructure:"attribute_count_limit,omitempty"`

	// AttributeValueLengthLimit corresponds to the JSON schema field
	// "attribute_value_length_limit".
	AttributeValueLengthLimit *int `json:"attribute_value_length_limit,omitempty" yaml:"attribute_value_length_limit,omitempty" mapstructure:"attribute_value_length_limit,omitempty"`
}

type MeterProviderJson map[string]interface{}

type ResourceJson struct {
	// Attributes corresponds to the JSON schema field "attributes".
	Attributes *Attributes `json:"attributes,omitempty" yaml:"attributes,omitempty" mapstructure:"attributes,omitempty"`
}

type TracerProviderJson struct {
	// SpanLimits corresponds to the JSON schema field "span_limits".
	SpanLimits *TracerProviderJsonSpanLimits `json:"span_limits,omitempty" yaml:"span_limits,omitempty" mapstructure:"span_limits,omitempty"`
}

type TracerProviderJsonSpanLimits struct {
	// AttributeCountLimit corresponds to the JSON schema field
	// "attribute_count_limit".
	AttributeCountLimit *int `json:"attribute_count_limit,omitempty" yaml:"attribute_count_limit,omitempty" mapstructure:"attribute_count_limit,omitempty"`

	// AttributeValueLengthLimit corresponds to the JSON schema field
	// "attribute_value_length_limit".
	AttributeValueLengthLimit *int `json:"attribute_value_length_limit,omitempty" yaml:"attribute_value_length_limit,omitempty" mapstructure:"attribute_value_length_limit,omitempty"`

	// EventAttributeCountLimit corresponds to the JSON schema field
	// "event_attribute_count_limit".
	EventAttributeCountLimit *int `json:"event_attribute_count_limit,omitempty" yaml:"event_attribute_count_limit,omitempty" mapstructure:"event_attribute_count_limit,omitempty"`

	// EventCountLimit corresponds to the JSON schema field "event_count_limit".
	EventCountLimit *int `json:"event_count_limit,omitempty" yaml:"event_count_limit,omitempty" mapstructure:"event_count_limit,omitempty"`

	// LinkAttributeCountLimit corresponds to the JSON schema field
	// "link_attribute_count_limit".
	LinkAttributeCountLimit *int `json:"link_attribute_count_limit,omitempty" yaml:"link_attribute_count_limit,omitempty" mapstructure:"link_attribute_count_limit,omitempty"`

	// LinkCountLimit corresponds to the JSON schema field "link_count_limit".
	LinkCountLimit *int `json:"link_count_limit,omitempty" yaml:"link_count_limit,omitempty" mapstructure:"link_count_limit,omitempty"`
}

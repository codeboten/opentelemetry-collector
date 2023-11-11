<!-- This file is autogenerated. See CONTRIBUTING.md for instructions to add an entry. -->

# Go API Changelog

This changelog includes only developer-facing changes.
If you are looking for user-facing changes, check out [CHANGELOG.md](./CHANGELOG.md).

<!-- next version -->

## v1.0.0-rcv0017/v0.88.0

### 💡 Enhancements 💡

- `pdata`: Add IsReadOnly() method to p[metrics|logs|traces].[Metrics|Logs|Spans] pdata structs allowing to check if the struct is read-only. (#6794)

## v1.0.0-rcv0016/v0.87.0

### 💡 Enhancements 💡

- `pdata`: Introduce API to control pdata mutability (#6794)
  This change introduces new API pdata methods to control the mutability:
  - p[metric|trace|log].[Metrics|Traces|Logs].MarkReadOnly() - marks the pdata as read-only. Any subsequent
    mutations will result in a panic.
  - p[metric|trace|log].[Metrics|Traces|Logs].IsReadOnly() - returns true if the pdata is marked as read-only.
  Currently, all the data is kept mutable. This API will be used by fanout consumer in the following releases. 

### 🛑 Breaking changes 🛑

- `obsreport`: remove methods/structs deprecated in previous release. (#8492)
- `extension`: remove deprecated Configs and Factories (#8631)

## v1.0.0-rcv0015/v0.86.0

### 🛑 Breaking changes 🛑

- `service`: remove deprecated service.PipelineConfig (#8485)

### 🚩 Deprecations 🚩

- `obsreporttest`: deprecate To*CreateSettings funcs in obsreporttest (#8492)
  The following TestTelemetry methods have been deprecated. Use structs instead:
  -  ToExporterCreateSettings -> exporter.CreateSettings
  -  ToProcessorCreateSettings -> processor.CreateSettings
  -  ToReceiverCreateSettings -> receiver.CreateSettings
  
- `obsreport`: Deprecating `obsreport.Exporter`, `obsreport.ExporterSettings`, `obsreport.NewExporter` (#8492)
  These deprecated methods/structs have been moved to exporterhelper:
  - `obsreport.Exporter` -> `exporterhelper.ObsReport`
  - `obsreport.ExporterSettings` -> `exporterhelper.ObsReportSettings`
  - `obsreport.NewExporter` -> `exporterhelper.NewObsReport`
  
- `obsreport`: Deprecating `obsreport.BuildProcessorCustomMetricName`, `obsreport.Processor`, `obsreport.ProcessorSettings`, `obsreport.NewProcessor` (#8492)
  These deprecated methods/structs have been moved to processorhelper:
  - `obsreport.BuildProcessorCustomMetricName` -> `processorhelper.BuildCustomMetricName`
  - `obsreport.Processor` -> `processorhelper.ObsReport`
  - `obsreport.ProcessorSettings` -> `processorhelper.ObsReportSettings`
  - `obsreport.NewProcessor` -> `processorhelper.NewObsReport`
  
- `obsreport`: Deprecating obsreport scraper and receiver API (#8492)
  These deprecated methods/structs have been moved to receiverhelper and scraperhelper:
  - `obsreport.Receiver` -> `receiverhelper.ObsReport`
  - `obsreport.ReceiverSettings` -> `receiverhelper.ObsReportSettings`
  - `obsreport.NewReceiver` -> `receiverhelper.NewObsReport`
  - `obsreport.Scraper` -> `scraperhelper.ObsReport`
  - `obsreport.ScraperSettings` -> `scraperhelper.ObsReportSettings`
  - `obsreport.NewScraper` -> `scraperhelper.NewObsReport`
  

### 💡 Enhancements 💡

- `otelcol`: Splitting otelcol into its own module. (#7924)
- `service`: Split service into its own module (#7923)

## v0.85.0

## v0.84.0

### 💡 Enhancements 💡

- `exporter/exporterhelper`: Introduce a new exporter helper that operates over client-provided requests instead of pdata (#7874)
  The following experimental API is introduced in exporter/exporterhelper package:
    - `NewLogsRequestExporter`: a new exporter helper for logs.
    - `NewMetricsRequestExporter`: a new exporter helper for metrics.
    - `NewTracesRequestExporter`: a new exporter helper for traces.
    - `Request`: an interface for client-defined requests.
    - `RequestItemsCounter`: an optional interface for counting the number of items in a Request.
    - `LogsConverter`: an interface for converting plog.Logs to Request.
    - `MetricsConverter`: an interface for converting pmetric.Metrics to Request.
    - `TracesConverter`: an interface for converting ptrace.Traces to Request.
    All the new APIs are intended to be used by exporters that need to operate over client-provided requests instead of pdata.
  
- `otlpreceiver`: Export HTTPConfig as part of the API for creating the otlpreceiver configuration. (#8175)
  Changes signature of receiver/otlpreceiver/config.go type httpServerSettings to HTTPConfig.

## v0.83.0

### 🛑 Breaking changes 🛑

- `all`: Remove go 1.19 support, bump minimum to go 1.20 and add testing for 1.21 (#8207)

### 💡 Enhancements 💡

- `changelog`: Generate separate changelogs for end users and package consumers (#8153)
# Developer Guide

Please review the [KServe Developer Guide](https://github.com/kserve/website/blob/main/docs/developer/developer.md) docs.

## Observability

KServe’s controller emits OpenTelemetry spans when the OTLP exporter is enabled. Set `ENABLE_OTEL_EXPORTER=true` in the controller deployment and point `OTEL_COLLECTOR_ENDPOINT` (or the standard `OTEL_EXPORTER_OTLP_ENDPOINT`) at your collector service. By default the exporter connects without TLS; override `OTEL_EXPORTER_OTLP_INSECURE=false` when you need a secure channel.


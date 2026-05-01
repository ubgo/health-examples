# 06-with-otel

Emit an OpenTelemetry span for every health check. Uses an in-memory exporter so the example runs standalone.

## Run

```sh
go run .
```

In production, replace `tracetest.NewInMemoryExporter` with an OTLP / stdout / Jaeger exporter.

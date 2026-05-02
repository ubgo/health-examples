# 06-with-otel

Emit an OpenTelemetry span for every health check. Uses an in-memory exporter so the example runs standalone.

## Run

```sh
go run .
```

The example logs captured spans (one per checker, per check cycle) so you can see what would have been sent to a real backend.

## What it demonstrates

- `healthotel.Register(reg, opts...)` subscribes the OTEL adapter to the registry's observer pattern.
- One span per check with attributes for `checker.name`, `checker.severity`, `result.status`, `result.latency_ms`.
- Errors recorded via `span.RecordError` + `SetStatus(codes.Error, …)`.
- No HTTP exposure of its own — works alongside any framework adapter.

## Real-world: production OTLP exporter

Replace the in-memory exporter with OTLP for production:

```go
import (
    "context"

    "github.com/ubgo/health"
    healthotel "github.com/ubgo/health/contrib/health-otel"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

ctx := context.Background()

exp, err := otlptracegrpc.New(ctx,
    otlptracegrpc.WithEndpoint("otel-collector.observability:4317"),
    otlptracegrpc.WithInsecure(),
)
if err != nil { panic(err) }

res, _ := resource.Merge(resource.Default(), resource.NewSchemaless(
    semconv.ServiceName("my-service"),
    semconv.ServiceVersion("1.4.2"),
))
tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exp),
    sdktrace.WithResource(res),
)
otel.SetTracerProvider(tp)

reg := health.NewRegistry()
healthotel.Register(reg, healthotel.WithTracer(tp.Tracer("health")))
```

## Why correlate health checks with traces?

When a request fails because a downstream is down, you want to jump straight from the failed-request span to the corresponding `health-check` span — same trace context, same time window. That's the value the adapter unlocks: every dependency outage is visible in your tracing backend, not just in `/readyz` HTTP logs.

## Composing with other observers

```go
healthotel.Register(reg, healthotel.WithTracer(otel.Tracer("health")))
healthprom.Register(reg, healthprom.WithRegisterer(promRegistry))
// + your own anonymous reg.Subscribe(...) for ad-hoc Slack alerts
```

The Registry's observer pattern is multi-subscriber — each adapter is independent of the others.

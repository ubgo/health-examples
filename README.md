# health-examples

Runnable example applications for [`github.com/ubgo/health`](https://github.com/ubgo/health).

Each example is a self-contained Go module under its own subdirectory. Clone, `cd`, and run `go run .`.

## Quick start

```sh
git clone https://github.com/ubgo/health-examples.git
cd health-examples/01-basic
go run .
```

## Examples

### Basics

- [`01-basic`](./01-basic) — register checkers, run them, print snapshots for liveness / readiness / startup. No HTTP.

### HTTP framework adapters

The same `Liveness` / `Readiness` / `Startup` / `Mount` shape across all five frameworks.

- [`02-http-nethttp`](./02-http-nethttp) — stdlib `net/http`.
- [`02b-http-gin`](./02b-http-gin) — Gin + auth-middleware example exposing `/internal/*`.
- [`02c-http-chi`](./02c-http-chi) — Chi + Recoverer.
- [`02d-http-echo`](./02d-http-echo) — Echo.
- [`02e-http-fiber`](./02e-http-fiber) — Fiber.

### Observability adapters

- [`06-with-otel`](./06-with-otel) — emit OTEL spans for every check (in-memory exporter).
- [`07-with-prometheus`](./07-with-prometheus) — expose health checks as Prometheus metrics at `/metrics`.

### Patterns

- [`08-custom-checker`](./08-custom-checker) — implement `health.Checker` for a custom dependency (SDK / client wrapping pattern).

## Compatibility

Requires Go 1.24 or later.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).

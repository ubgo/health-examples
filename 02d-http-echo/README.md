# 02d-http-echo

Health endpoints via Echo.

## Run

```sh
go run .
```

In another terminal:

```sh
curl http://localhost:8080/healthz   # liveness — always 200
curl http://localhost:8080/readyz    # readiness — 200 / 503
curl http://localhost:8080/startupz  # startup — 200 once warm
```

## What it demonstrates

- `healthecho.Mount(e, reg)` registers all three probes on an `*echo.Echo`.
- Echo's `MiddlewareFunc` plugs into `healthecho.WithMiddleware` directly — no shape conversion.
- Mounting under a route group works the same: `g := e.Group("/internal", auth)`, then `healthecho.Mount(g, reg)`.

## Same shape across frameworks

See [`02-http-nethttp`](../02-http-nethttp) for the canonical example. The `Mount` / `WithLivenessPath` / `WithReadinessPath` / `WithMiddleware` API is identical across `gin`, `chi`, `echo`, and `fiber`.

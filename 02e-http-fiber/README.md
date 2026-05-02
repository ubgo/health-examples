# 02e-http-fiber

Health endpoints via Fiber.

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

- `healthfiber.Mount(app, reg)` registers all three probes on a `*fiber.App`.
- Fiber's `fiber.Handler` plugs into `healthfiber.WithMiddleware` directly.
- Mounting under a route group works the same: `g := app.Group("/internal", auth)`, then `healthfiber.Mount(g, reg)`.

## Same shape across frameworks

See [`02-http-nethttp`](../02-http-nethttp) for the canonical example. The `Mount` / `WithLivenessPath` / `WithReadinessPath` / `WithMiddleware` API is identical across `gin`, `chi`, `echo`, and `fiber`.

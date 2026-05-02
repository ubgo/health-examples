# 02c-http-chi

Health endpoints via Chi.

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

- `healthchi.Mount(r, reg)` registers all three probes on any `chi.Router`.
- Composing with `chi/middleware.Recoverer` from chi's stock middleware.
- The `func(http.Handler) http.Handler` middleware shape works directly with `healthchi.WithMiddleware`.

## With protected internal endpoints

Chi's sub-router pattern wraps the probes in auth middleware:

```go
r := chi.NewRouter()
r.Route("/internal", func(internal chi.Router) {
    internal.Use(authMiddleware)
    healthchi.Mount(internal, reg)   // → /internal/healthz, /internal/readyz, /internal/startupz
})
```

## Same shape across frameworks

See [`02-http-nethttp`](../02-http-nethttp) for the canonical example. The other framework variants (`gin`, `echo`, `fiber`) share the same `Mount` / `WithLivenessPath` / `WithReadinessPath` / `WithMiddleware` API.

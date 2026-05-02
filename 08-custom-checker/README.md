# 08-custom-checker

Implement `health.Checker` for a custom dependency. Demonstrates the canonical adapter pattern: take an existing SDK / client, wrap its native error semantics in `health.Result`.

## Run

```sh
go run .
```

## What it demonstrates

- The `health.Checker` interface — only `Name() string` and `Check(ctx) Result`.
- How to map a third-party SDK's native error / status to `health.StatusUp` / `StatusDegraded` / `StatusDown`.
- How to surface useful detail (latency, downstream-error message, version) on the `Result` for renderers to expose.

## The canonical pattern

```go
package mydepcheck

import (
    "context"
    "time"

    "github.com/ubgo/health"
)

// thirdPartySDK is whatever client/library you wrap (Stripe, S3,
// Twilio, an internal RPC client, …).
type thirdPartySDK interface {
    Ping(ctx context.Context) error
}

type Checker struct {
    name string
    sdk  thirdPartySDK
}

// Conventional constructor signature for ubgo/health checkers.
func New(name string, sdk thirdPartySDK) *Checker {
    return &Checker{name: name, sdk: sdk}
}

func (c *Checker) Name() string { return c.name }

func (c *Checker) Check(ctx context.Context) health.Result {
    start := time.Now()
    err := c.sdk.Ping(ctx)
    latency := time.Since(start)

    if err != nil {
        return health.Result{
            Status:  health.StatusDown,
            Latency: latency,
            Error:   err.Error(),
        }
    }
    return health.Result{
        Status:  health.StatusUp,
        Latency: latency,
    }
}
```

Use it just like a built-in checker:

```go
reg := health.NewRegistry()
reg.Register(mydepcheck.New("stripe", stripeClient))
```

## Mapping nuanced statuses

Some dependencies have three-way health (Up / Degraded / Down) — replicas in sync vs lagging vs gone. Map them explicitly:

```go
func (c *Checker) Check(ctx context.Context) health.Result {
    info, err := c.sdk.GetReplicaStatus(ctx)
    if err != nil {
        return health.Result{Status: health.StatusDown, Error: err.Error()}
    }
    switch {
    case info.Lag > 30 * time.Second:
        return health.Result{
            Status: health.StatusDegraded,
            Error:  fmt.Sprintf("replica lag %v", info.Lag),
        }
    case info.Lag > 5 * time.Second:
        return health.Result{Status: health.StatusUp} // tolerable
    default:
        return health.Result{Status: health.StatusUp}
    }
}
```

Pair the `Degraded` status with `health.WithSeverity(health.SeverityDegraded)` at registration time — that way the readiness probe stays at 200 OK while still signalling the issue to observers (Prometheus alerts, OTEL spans).

## Don't forget the timeout

Every `Check` call gets a context with the per-checker timeout (default 5s). Always honour it:

```go
err := c.sdk.Ping(ctx) // SDK is responsible for ctx cancellation
```

If the SDK doesn't accept a context, wrap it:

```go
done := make(chan error, 1)
go func() { done <- c.sdk.Ping() }()
select {
case err := <-done:
    if err != nil { return health.Result{Status: health.StatusDown, Error: err.Error()} }
    return health.Result{Status: health.StatusUp}
case <-ctx.Done():
    return health.Result{Status: health.StatusDown, Error: "timeout"}
}
```

A check that ignores `ctx` blocks the registry's `RunAll` for the full SDK timeout — bad for `/readyz` latency.

## When to ship as a contrib

If you wrote a checker for a public dependency (open-source SDK, common SaaS API), consider contributing it as a `contrib/health-<name>/` module. The five built-in concrete checkers (`postgres`, `redis`, `nats`, `httpprobe`, `dns`) follow exactly the pattern shown here.

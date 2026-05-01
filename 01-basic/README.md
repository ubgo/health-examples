# 01-basic

Register a couple of checkers, run them, and print a snapshot for each probe type. No HTTP, no framework — just the registry.

## Run

```sh
go run .
```

## What it demonstrates

- `health.NewRegistry()` constructs the registry.
- Two checkers: a critical `db` (Up) and a degraded-severity `email` (Down).
- `reg.RunAll(ctx)` invokes each checker concurrently and stores the result.
- `reg.SnapshotForProbe(probe)` returns the aggregate state. Liveness is always Up. Readiness is Up because the only Down check is degraded-severity (so it doesn't fail readiness). Startup is Up once readiness is observed Up.

## Expected output

Three JSON blocks printed to stdout: one per probe type. Each block contains the aggregate `status`, the `time`, and the per-check breakdown under `checks`.

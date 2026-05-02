# 07-with-prometheus

Expose health checks as Prometheus metrics at `/metrics`. Three metrics per check:

- `health_check_status{name,severity}` — gauge, 1=up, 0=down, -1=degraded, -2=unknown
- `health_check_latency_seconds{name,severity}` — gauge
- `health_check_runs_total{name,severity,status}` — counter

## Run

```sh
go run .
curl -s http://localhost:8080/metrics | grep health_check_
```

## Sample queries

```promql
# Currently-down critical checks (alert on this)
health_check_status{severity="critical"} == 0

# 99th-percentile check latency over 5m, by checker
quantile_over_time(0.99, health_check_latency_seconds[5m])

# Flapping detector — checks that alternated up/down in last hour
changes(health_check_status[1h]) > 5

# % of recent runs that succeeded, by checker
sum by (name) (rate(health_check_runs_total{status="up"}[5m]))
  /
sum by (name) (rate(health_check_runs_total[5m]))
```

## Sample alert recipe

```yaml
# alertmanager rule
- alert: HealthCheckCritical
  expr: health_check_status{severity="critical"} == 0
  for: 1m
  annotations:
    summary: "Critical dependency {{ $labels.name }} is down"
    description: "{{ $labels.name }} (severity={{ $labels.severity }}) has been down for over 1 minute."

- alert: HealthCheckDegraded
  expr: health_check_status{severity="degraded"} == -1
  for: 5m
  annotations:
    summary: "Degraded dependency {{ $labels.name }}"
    description: "{{ $labels.name }} is in degraded state — service still serving but a non-critical dependency is unreachable."
```

The 5-minute `for:` on degraded prevents pager noise from blip Slack outages, while critical alerts in 1 minute.

## Composing with the OTEL observer

```go
healthprom.Register(reg, healthprom.WithRegisterer(promRegistry))
healthotel.Register(reg, healthotel.WithTracer(otel.Tracer("health")))
```

You get queryable historical data (Prom) AND per-check trace spans (OTEL) from the same Registry. Wire `/metrics` into your existing scrape config — no separate exporter to deploy.

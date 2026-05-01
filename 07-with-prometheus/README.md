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

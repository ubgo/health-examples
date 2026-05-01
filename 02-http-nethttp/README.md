# 02-http-nethttp

Expose `/healthz`, `/readyz`, `/startupz` via stdlib `net/http`. Background re-checks every 5 seconds. Graceful shutdown on SIGINT / SIGTERM.

## Run

```sh
go run .
```

In another terminal:

```sh
curl http://localhost:8080/healthz   # liveness — always 200
curl http://localhost:8080/readyz    # readiness — 200 / 503
curl http://localhost:8080/startupz  # startup — 200 once ready, 503 before
```

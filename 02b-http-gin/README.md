# 02b-http-gin

Health endpoints via Gin. Public endpoints at `/healthz`, `/readyz`, `/startupz` for k8s probes; the same endpoints duplicated under `/internal/` and protected by `X-Internal-Key` middleware for richer detail.

## Run

```sh
go run .
```

```sh
curl http://localhost:8080/readyz
curl -H 'X-Internal-Key: secret-rotate-me' http://localhost:8080/internal/readyz
```

// Example 07-with-prometheus — expose health checks as Prometheus metrics
// at /metrics, with /readyz on the same server for liveness/readiness probes.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ubgo/health"
	healthprom "github.com/ubgo/health/contrib/health-prom"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	promReg := prometheus.NewRegistry()
	reg := health.NewRegistry()
	healthprom.Register(reg, promReg)

	_ = reg.Register(up{name: "db"})
	_ = reg.Register(up{name: "cache"})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(promReg, promhttp.HandlerOpts{}))

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Println("listening on :8080 — try: curl http://localhost:8080/metrics | grep health_check_")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

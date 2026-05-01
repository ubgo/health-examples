// Example 02-http-nethttp — expose /healthz, /readyz, /startupz via stdlib net/http.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ubgo/health"
	healthnethttp "github.com/ubgo/health/contrib/health-nethttp"
)

type fakeChecker struct {
	name string
	up   bool
}

func (f fakeChecker) Name() string { return f.name }
func (f fakeChecker) Check(_ context.Context) health.Result {
	if f.up {
		return health.Result{Status: health.StatusUp}
	}
	return health.Result{Status: health.StatusDown, Error: "simulated"}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	reg := health.NewRegistry()
	_ = reg.Register(fakeChecker{name: "db", up: true})
	_ = reg.Register(fakeChecker{name: "cache", up: true})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	mux := http.NewServeMux()
	healthnethttp.Mount(mux, reg)

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Println("listening on :8080 — try: curl http://localhost:8080/readyz")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

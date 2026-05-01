// Example 02c-http-chi — health endpoints via Chi.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ubgo/health"
	healthchi "github.com/ubgo/health/contrib/health-chi"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	reg := health.NewRegistry()
	_ = reg.Register(up{name: "db"})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	healthchi.Mount(r, reg)

	srv := &http.Server{Addr: ":8080", Handler: r, ReadHeaderTimeout: 5 * time.Second}
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

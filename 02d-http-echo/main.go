// Example 02d-http-echo — health endpoints via Echo.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ubgo/health"
	healthecho "github.com/ubgo/health/contrib/health-echo"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	e := echo.New()
	e.HideBanner = true

	reg := health.NewRegistry()
	_ = reg.Register(up{name: "db"})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	healthecho.Mount(e, reg)

	go func() {
		log.Println("listening on :8080 — try: curl http://localhost:8080/readyz")
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = e.Shutdown(shutdownCtx)
}

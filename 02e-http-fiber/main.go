// Example 02e-http-fiber — health endpoints via Fiber.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ubgo/health"
	healthfiber "github.com/ubgo/health/contrib/health-fiber"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	reg := health.NewRegistry()
	_ = reg.Register(up{name: "db"})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	healthfiber.Mount(app, reg)

	go func() {
		log.Println("listening on :8080 — try: curl http://localhost:8080/readyz")
		if err := app.Listen(":8080"); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(shutdownCtx)
}

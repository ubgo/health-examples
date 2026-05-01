// Example 02b-http-gin — health endpoints via Gin, with optional auth middleware.
package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ubgo/health"
	healthgin "github.com/ubgo/health/contrib/health-gin"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func internalKeyAuth(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if subtle.ConstantTimeCompare([]byte(c.GetHeader("X-Internal-Key")), []byte(expected)) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	reg := health.NewRegistry()
	_ = reg.Register(up{name: "db"})
	reg.Start(ctx, 5*time.Second)
	defer reg.Stop()

	// Public liveness/readiness/startup at default paths.
	healthgin.Mount(r, reg)

	// Detailed internal endpoint guarded by X-Internal-Key.
	internal := r.Group("/internal", internalKeyAuth("secret-rotate-me"))
	healthgin.Mount(internal, reg,
		healthgin.WithLivenessPath("/healthz"),
		healthgin.WithReadinessPath("/readyz"),
		healthgin.WithStartupPath("/startupz"),
	)

	srv := &http.Server{Addr: ":8080", Handler: r, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Println("listening on :8080")
		log.Println("  public:   curl http://localhost:8080/readyz")
		log.Println("  internal: curl -H 'X-Internal-Key: secret-rotate-me' http://localhost:8080/internal/readyz")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

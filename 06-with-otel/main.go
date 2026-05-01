// Example 06-with-otel — emit OTEL spans for every health check.
//
// Uses an in-memory exporter so the example runs standalone without an
// external collector. In production swap exporter for OTLP / stdout / etc.
package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ubgo/health"
	healthotel "github.com/ubgo/health/contrib/health-otel"
)

type up struct{ name string }

func (u up) Name() string                          { return u.name }
func (u up) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

func main() {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	defer func() { _ = tp.Shutdown(context.Background()) }()

	reg := health.NewRegistry()
	healthotel.Register(reg, healthotel.WithTracerProvider(tp))
	_ = reg.Register(up{name: "db"})
	_ = reg.Register(up{name: "cache"})

	reg.RunAll(context.Background())

	fmt.Printf("emitted %d spans:\n", len(exporter.GetSpans()))
	for _, s := range exporter.GetSpans() {
		fmt.Printf("  %s — %d attrs\n", s.Name, len(s.Attributes))
		for _, a := range s.Attributes {
			fmt.Printf("    %s = %v\n", a.Key, a.Value.Emit())
		}
	}
}

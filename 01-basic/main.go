// Example 01-basic — register checkers, run them, print snapshots.
//
// Demonstrates the core API: NewRegistry, Register, RunAll, SnapshotForProbe.
// No HTTP server, no framework — just the registry itself.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ubgo/health"
)

// alwaysUp is a trivial Checker that reports Up.
type alwaysUp struct{ name string }

func (a alwaysUp) Name() string                          { return a.name }
func (a alwaysUp) Check(_ context.Context) health.Result { return health.Result{Status: health.StatusUp} }

// alwaysDown reports Down with an error string.
type alwaysDown struct {
	name, msg string
}

func (a alwaysDown) Name() string { return a.name }
func (a alwaysDown) Check(_ context.Context) health.Result {
	return health.Result{Status: health.StatusDown, Error: a.msg}
}

func main() {
	reg := health.NewRegistry()

	// Critical (default): a Down result fails readiness.
	if err := reg.Register(alwaysUp{name: "db"}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Degraded severity: Down is reported but readiness stays up.
	_ = reg.Register(alwaysDown{name: "email", msg: "smtp unreachable"}, health.WithSeverity(health.SeverityDegraded))

	reg.RunAll(context.Background())

	for _, probe := range []health.Probe{health.ProbeLiveness, health.ProbeReadiness, health.ProbeStartup} {
		fmt.Printf("--- %s snapshot ---\n", probe)
		snap := reg.SnapshotForProbe(probe)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(snap)
	}
}

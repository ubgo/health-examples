// Example 08-custom-checker — implement health.Checker for a custom dependency.
//
// In production a Checker wraps an existing client (DB pool, HTTP client,
// SDK) and translates its native error semantics into health.Result. This
// example simulates a feature-flag service.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ubgo/health"
)

// FeatureFlagsService is the kind of dependency we'd usually inject from an SDK.
type FeatureFlagsService struct {
	ready bool
}

func (s *FeatureFlagsService) IsReady(ctx context.Context) error {
	if !s.ready {
		return errors.New("flags: control-plane unreachable")
	}
	return nil
}

// FeatureFlagsChecker adapts the SDK-style ready check into health.Checker.
type FeatureFlagsChecker struct {
	name string
	svc  *FeatureFlagsService
}

func NewFeatureFlagsChecker(name string, svc *FeatureFlagsService) *FeatureFlagsChecker {
	return &FeatureFlagsChecker{name: name, svc: svc}
}

func (c *FeatureFlagsChecker) Name() string { return c.name }

func (c *FeatureFlagsChecker) Check(ctx context.Context) health.Result {
	start := time.Now()
	err := c.svc.IsReady(ctx)
	latency := time.Since(start)

	r := health.Result{
		Severity:  health.SeverityCritical,
		Latency:   latency,
		Timestamp: time.Now(),
	}
	if err != nil {
		r.Status = health.StatusDown
		r.Error = err.Error()
		return r
	}
	r.Status = health.StatusUp
	r.Metadata = map[string]any{"sdk_version": "1.2.3"}
	return r
}

func main() {
	svc := &FeatureFlagsService{ready: true}
	reg := health.NewRegistry()
	_ = reg.Register(NewFeatureFlagsChecker("flags", svc))

	reg.RunAll(context.Background())
	fmt.Println("--- Up case ---")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(reg.SnapshotForProbe(health.ProbeReadiness))

	svc.ready = false
	reg.RunAll(context.Background())
	fmt.Println("--- Down case ---")
	_ = enc.Encode(reg.SnapshotForProbe(health.ProbeReadiness))
}

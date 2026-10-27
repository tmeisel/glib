package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRegisterHealthCheckValidation(t *testing.T) {
	m := NewMonitor()
	check := HealthCheckFunc(func(context.Context) error { return nil })

	if err := m.RegisterHealthCheck("", check); err == nil {
		t.Fatal("expected empty name to be rejected")
	}

	if err := m.RegisterHealthCheck("nil", nil); err == nil {
		t.Fatal("expected nil function to be rejected")
	}

	if err := m.RegisterHealthCheck("db", check); err != nil {
		t.Fatalf("register check: %v", err)
	}

	if err := m.RegisterHealthCheck("db", check); err == nil {
		t.Fatal("expected duplicate name to be rejected")
	}

	if err := m.RegisterHealthCheck("multiple-options", check, Options{}, Options{}); err == nil {
		t.Fatal("expected multiple Options values to be rejected")
	}
}

func TestStatusTracksFailureAndRecovery(t *testing.T) {
	m := NewMonitor(10 * time.Millisecond)
	failure := errors.New("database unavailable")

	var runs atomic.Int32
	if err := m.RegisterHealthCheck("database", func(context.Context) error {
		if runs.Add(1) == 1 {
			return failure
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer m.Stop()

	waitFor(t, func() bool {
		state := m.Status().HealthChecks["database"]
		return state.Error != nil && *state.Error != nil && (*state.Error).Error() == failure.Error()
	})

	if status := m.Status(); status.Healthy {
		t.Fatal("overall status should be unhealthy after a failed check")
	}

	waitFor(t, func() bool {
		state := m.Status().HealthChecks["database"]
		return state.Healthy && state.Error == nil && state.LastHealthy != nil
	})

	status := m.Status()

	if !status.Healthy {
		t.Fatal("overall status should be healthy after recovery")
	}

	if status.LastHealthy == nil {
		t.Fatal("overall LastHealthy should be set after success")
	}
}

func TestHealthCheckIntervalOverride(t *testing.T) {
	m := NewMonitor(200 * time.Millisecond)
	var defaultRuns, fastRuns atomic.Int32
	if err := m.RegisterHealthCheck("default", func(context.Context) error {
		defaultRuns.Add(1)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := m.RegisterHealthCheck("fast", func(context.Context) error {
		fastRuns.Add(1)
		return nil
	}, Options{HealthCheckInterval: 10 * time.Millisecond}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := m.Start(ctx); err != nil {
		cancel()
		t.Fatal(err)
	}
	defer func() {
		cancel()
		m.Stop()
	}()

	waitFor(t, func() bool { return fastRuns.Load() >= 4 })
	if got := defaultRuns.Load(); got != 1 {
		t.Fatalf("default interval check ran %d times before override check reached four runs; want 1", got)
	}
}

func TestHandlerReturnsJSONStatus(t *testing.T) {
	m := NewMonitor(time.Hour)
	failure := errors.New("upstream timeout")
	if err := m.RegisterHealthCheck("upstream", func(context.Context) error { return failure }); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := m.Start(ctx); err != nil {
		cancel()
		t.Fatal(err)
	}
	defer func() {
		cancel()
		m.Stop()
	}()
	waitFor(t, func() bool { return m.Status().HealthChecks["upstream"].Error != nil })

	recorder := httptest.NewRecorder()
	m.Handler().ServeHTTP(recorder, httptest.NewRequest("GET", "/health", nil))
	if recorder.Code != 200 {
		t.Fatalf("status code = %d, want 200", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	var got struct {
		Healthy      bool `json:"healthy"`
		HealthChecks map[string]struct {
			Healthy bool    `json:"healthy"`
			Error   *string `json:"error"`
		} `json:"healthChecks"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	state := got.HealthChecks["upstream"]
	if got.Healthy || state.Healthy || state.Error == nil || *state.Error != failure.Error() {
		t.Fatalf("unexpected JSON status: %+v", got)
	}
}

func TestContextCancellationStopsChecks(t *testing.T) {
	m := NewMonitor()
	entered := make(chan struct{}, 1)
	if err := m.RegisterHealthCheck("blocking", func(ctx context.Context) error {
		entered <- struct{}{}
		<-ctx.Done()
		return ctx.Err()
	}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := m.Start(ctx); err != nil {
		cancel()
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		cancel()
		m.Stop()
		t.Fatal("health check did not start")
	}
	cancel()
	stopped := make(chan struct{})
	go func() {
		m.Stop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("Stop did not return after context cancellation")
	}
}

func TestCanceledContextSkipsFirstRun(t *testing.T) {
	m := NewMonitor()

	var runs atomic.Int32

	if err := m.RegisterHealthCheck("check", func(context.Context) error {
		runs.Add(1)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	m.Stop()

	if got := runs.Load(); got != 0 {
		t.Fatalf("health check ran %d times with an already-canceled context; want 0", got)
	}
}

func TestHandlerRejectsUnsupportedMethod(t *testing.T) {
	recorder := httptest.NewRecorder()

	NewMonitor().Handler().ServeHTTP(recorder, httptest.NewRequest("POST", "/health", nil))

	if recorder.Code != 405 {
		t.Fatalf("status code = %d, want 405", recorder.Code)
	}

	if got := recorder.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("Allow = %q, want GET, HEAD", got)
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatal("condition was not met before timeout")
}

// Package healthcheck runs periodic health checks and exposes their status as JSON.
package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

// DefaultDuration is used when neither NewMonitor nor a health check specifies an interval.
const DefaultDuration = 30 * time.Second

// HealthCheckFunc is a function that reports whether a dependency or service is healthy.
type HealthCheckFunc func(context.Context) error

// Options configures a registered health check.
type Options struct {
	// HealthCheckInterval overrides the module's default interval when greater than zero.
	HealthCheckInterval time.Duration
}

// HealthCheck is the most recently observed state of one registered check.
// Error contains the latest returned error, or nil when the latest run succeeded
// or the check has not run yet. LastHealthy is nil until a run succeeds.
type HealthCheck struct {
	Healthy     bool       `json:"healthy"`
	Error       *error     `json:"error"`
	LastHealthy *time.Time `json:"lastHealthy"`
}

// MarshalJSON encodes errors as useful strings instead of the empty JSON object
// produced by encoding/json for concrete error implementations.
func (h HealthCheck) MarshalJSON() ([]byte, error) {
	var errorText *string
	if h.Error != nil && *h.Error != nil {
		message := (*h.Error).Error()
		errorText = &message
	}
	return json.Marshal(struct {
		Healthy     bool       `json:"healthy"`
		Error       *string    `json:"error"`
		LastHealthy *time.Time `json:"lastHealthy"`
	}{Healthy: h.Healthy, Error: errorText, LastHealthy: h.LastHealthy})
}

// ServerStatus is the combined state of the module and all registered checks.
type ServerStatus struct {
	Healthy      bool                   `json:"healthy"`
	LastHealthy  *time.Time             `json:"lastHealthy"`
	HealthChecks map[string]HealthCheck `json:"healthChecks"`
}

type registeredCheck struct {
	fn       HealthCheckFunc
	interval time.Duration
}

// HealthCheckMonitor manages a set of periodically executed health checks.
type HealthCheckMonitor struct {
	mu              sync.RWMutex
	checks          map[string]registeredCheck
	states          map[string]HealthCheck
	defaultDuration time.Duration
	started         bool
	cancel          context.CancelFunc
	done            chan struct{}
}

// NewMonitor creates a HealthCheckMonitor. If no positive duration is supplied, DefaultDuration is used.
func NewMonitor(defaultDuration ...time.Duration) *HealthCheckMonitor {
	duration := DefaultDuration
	if len(defaultDuration) > 0 && defaultDuration[0] > 0 {
		duration = defaultDuration[0]
	}
	return &HealthCheckMonitor{
		checks:          make(map[string]registeredCheck),
		states:          make(map[string]HealthCheck),
		defaultDuration: duration,
		done:            make(chan struct{}),
	}
}

// RegisterHealthCheck registers a check before Start is called. The optional
// interval overrides the module default when it is greater than zero.
func (m *HealthCheckMonitor) RegisterHealthCheck(name string, fn HealthCheckFunc, options ...Options) error {
	if name == "" {
		return errors.New("healthcheck: name must not be empty")
	}
	if fn == nil {
		return errors.New("healthcheck: function must not be nil")
	}
	if len(options) > 1 {
		return errors.New("healthcheck: at most one Options value may be provided")
	}

	interval := m.defaultDuration
	if len(options) == 1 && options[0].HealthCheckInterval > 0 {
		interval = options[0].HealthCheckInterval
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return errors.New("healthcheck: cannot register checks after Start")
	}

	if _, exists := m.checks[name]; exists {
		return errors.New("healthcheck: duplicate name: " + name)
	}

	m.checks[name] = registeredCheck{fn: fn, interval: interval}
	m.states[name] = HealthCheck{}

	return nil
}

// Start begins running every check immediately, then repeats it at its interval.
// It returns after starting the workers. The supplied context controls their lifetime.
func (m *HealthCheckMonitor) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("healthcheck: context must not be nil")
	}

	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return errors.New("healthcheck: already started")
	}

	m.started = true

	runCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	checks := make(map[string]registeredCheck, len(m.checks))
	for name, check := range m.checks {
		checks[name] = check
	}
	m.mu.Unlock()

	var workers sync.WaitGroup
	workers.Add(len(checks))

	for name, check := range checks {
		go func(name string, check registeredCheck) {
			defer workers.Done()
			m.runCheck(runCtx, name, check)
		}(name, check)
	}

	go func() {
		workers.Wait()
		close(m.done)
	}()

	return nil
}

func (m *HealthCheckMonitor) runCheck(ctx context.Context, name string, check registeredCheck) {
	ticker := time.NewTicker(check.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m.execute(ctx, name, check.fn)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (m *HealthCheckMonitor) execute(ctx context.Context, name string, fn HealthCheckFunc) {
	err := fn(ctx)
	m.mu.Lock()
	state := m.states[name]
	if err != nil {
		state.Healthy = false
		state.Error = &err
	} else {
		now := time.Now().UTC()
		state.Healthy = true
		state.Error = nil
		state.LastHealthy = &now
	}
	m.states[name] = state
	m.mu.Unlock()
}

// Stop cancels all checks and waits for their goroutines to exit. It is safe to
// call repeatedly. Health checks should return when their context is canceled.
func (m *HealthCheckMonitor) Stop() {
	m.mu.RLock()
	started, cancel, done := m.started, m.cancel, m.done
	m.mu.RUnlock()
	if !started {
		return
	}
	cancel()
	<-done
}

// Status returns a consistent snapshot of the combined and per-check status.
// The module is healthy only when at least one check is registered and all checks
// have succeeded at least once and have no current error.
func (m *HealthCheckMonitor) Status() ServerStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	checks := make(map[string]HealthCheck, len(m.states))
	allHealthy := len(m.states) > 0
	var lastHealthy *time.Time
	for name, state := range m.states {
		checks[name] = state
		if !state.Healthy {
			allHealthy = false
		}
		if state.LastHealthy != nil && (lastHealthy == nil || state.LastHealthy.After(*lastHealthy)) {
			value := *state.LastHealthy
			lastHealthy = &value
		}
	}
	return ServerStatus{
		Healthy:      allHealthy,
		LastHealthy:  lastHealthy,
		HealthChecks: checks,
	}
}

// Handler returns an HTTP handler that serves the current status as JSON.
func (m *HealthCheckMonitor) Handler() http.Handler {
	return http.HandlerFunc(m.ServeHTTP)
}

// ServeHTTP writes the current status as JSON. Only GET and HEAD are accepted.
func (m *HealthCheckMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodHead {
		return
	}

	_ = json.NewEncoder(w).Encode(m.Status())
}

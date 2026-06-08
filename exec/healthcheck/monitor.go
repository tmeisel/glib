package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"sort"
	"time"

	ctxPkg "github.com/tmeisel/glib/ctx"
	"github.com/tmeisel/glib/net/http/response"
	"github.com/tmeisel/glib/utils/ptr"
)

type Monitor struct {
	hcInterval time.Duration
	hcs        []healthCheck
}

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

type healthCheck struct {
	Name   string       `json:"name"`
	Status HealthStatus `json:"status"`
	fn     HealthCheckFunc
	ch     chan error

	LastHealthy *timestamp `json:"lastHealthy"`
}

type HealthStatus struct {
	e error
}

func (hs HealthStatus) MarshalJSON() ([]byte, error) {
	if hs.e == nil {
		return json.Marshal(StatusHealthy)
	}

	return json.Marshal(hs.e.Error())
}

type timestamp time.Time

func (t *timestamp) MarshalJSON() ([]byte, error) {
	if t == nil {
		return json.Marshal("never")
	}

	return json.Marshal(time.Time(*t))
}

type HealthCheckFunc func(ctx context.Context, statusChan chan error, interval *time.Duration)

func NewHealthCheckMonitor(hcInterval time.Duration) *Monitor {
	return &Monitor{
		hcInterval: hcInterval,
		hcs:        make([]healthCheck, 0),
	}
}

func (s *Monitor) Start(ctx context.Context) {
	logger := ctxPkg.GetLogger(ctx)

	selectFrom := make([]reflect.SelectCase, len(s.hcs))

	for idx, hc := range s.hcs {
		go hc.fn(ctx, hc.ch, ptr.Duration(s.hcInterval))

		selectFrom[idx] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(hc.ch),
		}
	}

	// append done channel
	selectFrom = append(selectFrom, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())})

	for {
		hcIdx, value, ok := reflect.Select(selectFrom)

		logger.Debugf(ctx, "received value from select case %d: %v", hcIdx, value)

		// is it the done chan, which was closed?
		if ok == false {
			if errors.Is(ctx.Err(), context.Canceled) {
				logger.Info(ctx, "stopping health checker")
				return
			}

			logger.Infof(ctx, "invalid receive")
			continue
		}

		if value.IsNil() {
			ts := timestamp(time.Now())

			s.hcs[hcIdx].Status = HealthStatus{e: nil} // unset previous error
			s.hcs[hcIdx].LastHealthy = &ts             // last healthy is now

			continue
		}

		err, ok := value.Interface().(error)
		if ok {
			s.hcs[hcIdx].Status = HealthStatus{e: err}
		}
	}
}

type Response struct {
	Status       string        `json:"status"`
	LastHealthy  *timestamp    `json:"lastHealthy"`
	HealthChecks []healthCheck `json:"healthChecks"`
}

func (s *Monitor) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	resp := Response{
		Status:       StatusHealthy,
		HealthChecks: s.hcs,
	}

	for _, hc := range s.hcs {
		if hc.LastHealthy == nil {
			continue
		}

		resp.Status = StatusUnhealthy
	}

	statusCode := http.StatusOK
	if resp.LastHealthy != nil {
		statusCode = http.StatusServiceUnavailable
	}

	response.WriteJson(w, statusCode, resp)
}

func (s *Monitor) AddHealthCheck(name string, hc HealthCheckFunc) {
	s.hcs = append(s.hcs, healthCheck{
		Name:        name,
		Status:      HealthStatus{e: nil},
		LastHealthy: nil,
		fn:          hc,
		ch:          make(chan error),
	})

	sort.Slice(s.hcs, func(i, j int) bool {
		return s.hcs[i].Name < s.hcs[j].Name
	})
}

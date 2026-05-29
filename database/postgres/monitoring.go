package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	ctxPkg "github.com/tmeisel/glib/ctx"
	"github.com/tmeisel/glib/utils/ptr"
)

type Monitoring struct {
	pool       *pgxpool.Pool
	statusChan chan error
}

const (
	defaultMonitoringInterval = time.Second * 10
	minMonitoringInterval     = time.Second * 5
)

func NewMonitoring(pool *pgxpool.Pool, statusChan chan error) *Monitoring {
	return &Monitoring{
		pool:       pool,
		statusChan: statusChan,
	}
}

// Monitor starts health checking the connection with IdAM.
// It keeps running until the given [context.Context] expires
func (m Monitoring) Monitor(ctx context.Context, interval *time.Duration) {
	logger := ctxPkg.GetLogger(ctx)

	if interval == nil {
		interval = ptr.Duration(defaultMonitoringInterval)
	}

	if *interval <= minMonitoringInterval {
		interval = ptr.Duration(minMonitoringInterval)
	}

	logger.Infof(ctx, "starting postgres monitoring with an interval of %s", interval.Truncate(time.Second))

	checkChan := time.Tick(*interval)

	for {
		select {
		case <-checkChan:
			go m.checkHealth(ctx)
		case <-ctx.Done():
			logger.Info(ctx, "stopping postgres monitoring")
			return
		}
	}
}

func (m Monitoring) checkHealth(ctx context.Context) {
	logger := ctxPkg.GetLogger(ctx)

	logger.Debug(ctx, "checking postgres health")

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := m.pool.Ping(ctx); err != nil {
		logger.Errorf(ctx, "postgres health check failed: %s", err)
		m.statusChan <- err

		return
	}

	logger.Debug(ctx, "postgres health check succeeded")

	m.statusChan <- nil
}

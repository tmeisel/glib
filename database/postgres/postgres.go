package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	errPkg "github.com/tmeisel/glib/error"
	"github.com/tmeisel/glib/exec/backoff"
)

func Init(ctx context.Context, conf Config, retryConf *RetryConfig) (pool *pgxpool.Pool, err error) {
	if retryConf == nil {
		return pgxpool.New(ctx, conf.DSN())
	}

	b := retryConf.Backoff
	b.With(retryConf.Options...)

	pool, err = pgxpool.New(ctx, conf.DSN())
	if err != nil {
		return nil, errPkg.NewInternal(err)
	}

	err = b.Do(ctx, func(ctx context.Context) error {
		err = pool.Ping(ctx)
		if err != nil {
			return nil
		}

		if errors.Is(err, &pgconn.ConnectError{}) {
			return backoff.RetryableError(err)
		}

		return err
	})

	return
}

package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tmeisel/glib/exec/backoff"
)

func testInitDb() {
	retryConf := &RetryConfig{
		Backoff: backoff.NewFibonacci(time.Millisecond * 50),
		Options: []backoff.OptionFn{
			backoff.WithJitter(20),
			backoff.WithCap(time.Second * 3),
		},
	}

	Init(context.Background(), Config{}, retryConf)
}

func Init(ctx context.Context, conf Config, retryConf *RetryConfig) (pool *pgxpool.Pool, err error) {
	if retryConf == nil {
		return pgxpool.New(ctx, conf.DSN())
	}

	b := retryConf.Backoff
	b.With(retryConf.Options...)

	err = b.Do(ctx, func(ctx context.Context) error {
		pool, err = pgxpool.New(ctx, conf.DSN())

		if err != nil {
			if errors.Is(err, &pgconn.ConnectError{}) {
				return backoff.RetryableError(err)
			}

			return err
		}

		return nil
	})

	return
}

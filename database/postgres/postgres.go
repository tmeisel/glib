package postgres

import (
	"context"
	"errors"
	"fmt"
	stdLibLog "log"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	ctxPkg "github.com/tmeisel/glib/ctx"
	errPkg "github.com/tmeisel/glib/error"
	"github.com/tmeisel/glib/exec/backoff"
	"github.com/tmeisel/glib/log/fields"
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

	var log = defaultLog
	if logger := ctxPkg.GetLogger(ctx); logger != nil {
		log = logger.Error
	}

	err = b.Do(ctx, func(ctx context.Context) error {
		err = pool.Ping(ctx)
		if err == nil {
			return nil
		}

		var pgErr *pgconn.ConnectError
		if errors.As(err, &pgErr) {
			if errPkg.Is(ProcessError(err), errPkg.CodeInvalidCredentials) {
				return err
			}

			log(ctx, "failed to connect to postgres", fields.Bool("retry", true), fields.String("err", err.Error()))
			return backoff.RetryableError(pgErr)
		}

		log(ctx, "failed to connect to postgres", fields.Bool("retry", false), fields.String("err", err.Error()))

		return errPkg.NewInternal(err)
	})

	return
}

// defaultLog writes the given msg using log.Printf (incl. fields.Field f)
func defaultLog(_ context.Context, msg string, f ...fields.Field) {
	var logFields []string
	for _, field := range f {
		logFields = append(logFields, fmt.Sprintf("%s=%v", field.Key, field.Interface))
	}

	// stdLibLog is the standard library log
	stdLibLog.Printf("%s (fields: %s)", msg, strings.Join(logFields, ","))
}

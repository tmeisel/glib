package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"

	"github.com/tmeisel/glib/database"
)

const (
	CodeDuplicateKey = "23505"
)

func ProcessError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return database.ErrNoRows
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		column := pgErr.ColumnName
		if column == "" {
			column = pgErr.ConstraintName
		}

		switch pgErr.Code {
		case CodeDuplicateKey:
			return database.NewDuplicateKeyError(pgErr, &column)
		}
	}

	return database.NewDbError(err)
}

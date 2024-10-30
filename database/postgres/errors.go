package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/tmeisel/glib/database"
)

const (
	CodeDuplicateKey = "23505"
)

func ProcessError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return database.ErrNoRows
	}

	if pgconnErr, ok := err.(*pgconn.PgError); ok {
		column := pgconnErr.ColumnName
		if column == "" {
			column = pgconnErr.ConstraintName
		}

		switch pgconnErr.Code {
		case CodeDuplicateKey:
			return database.NewDuplicateKeyError(pgconnErr, &column)
		}
	}

	return database.NewDbError(err)
}

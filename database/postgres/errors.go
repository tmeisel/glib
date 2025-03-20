package postgres

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/tmeisel/glib/database"
	errPkg "github.com/tmeisel/glib/error"
)

const (
	CodeDuplicateKey = "23505"
	CodeInvalidLogin = "28P01"
)

func ProcessError(err error) *errPkg.Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return database.ErrNoRows
	}

	var pgErr *pgconn.ConnectError
	if errors.As(err, &pgErr) {
		if strings.Contains(pgErr.Error(), CodeInvalidLogin) {
			return database.ErrInvalidLogin
		}
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

	return database.NewError(err)
}

func IsDuplicateKeyError(err error) bool {
	if pkgErr, ok := err.(*errPkg.Error); ok {
		return pkgErr.GetCode() == errPkg.CodeDuplicateKey
	}

	if pgconnErr, ok := err.(*pgconn.PgError); ok {
		return pgconnErr.Code == CodeDuplicateKey
	}

	return false
}

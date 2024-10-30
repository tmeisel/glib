package database

import (
	"fmt"

	errorPkg "github.com/tmeisel/glib/error"
)

var (
	ErrNoRows = errorPkg.New(errorPkg.CodeNotFound, errorPkg.CodeNotFound.String(), nil)
)

type Error errorPkg.Error

func NewDbError(err error) error {
	return NewDbErrorMsg(err, errorPkg.CodeInternal.String())
}

func NewDbErrorMsg(err error, msg string) error {
	return errorPkg.New(
		errorPkg.CodeInternal,
		msg,
		err,
	)
}

func NewDuplicateKeyError(prev error, column *string) error {
	msg := "duplicate key error"
	if column != nil {
		msg = fmt.Sprintf("duplicate key in column %s", *column)
	}

	return errorPkg.New(errorPkg.CodeConflict, msg, prev)
}

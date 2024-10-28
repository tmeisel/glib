package database

import errorPkg "github.com/tmeisel/glib/error"

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

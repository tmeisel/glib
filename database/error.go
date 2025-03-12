package database

import (
	"fmt"

	errorPkg "github.com/tmeisel/glib/error"
)

var (
	ErrInvalidLogin = errorPkg.New(errorPkg.CodeInvalidCredentials, errorPkg.CodeInvalidCredentials.String(), nil)
	ErrNoRows       = errorPkg.New(errorPkg.CodeNotFound, errorPkg.CodeNotFound.String(), nil)
)

func NewError(err error) *errorPkg.Error {
	return NewErrorMsg(err, errorPkg.CodeInternal.String())
}

func NewErrorMsg(err error, msg string) *errorPkg.Error {
	return errorPkg.New(errorPkg.CodeInternal, msg, err)
}

func NewDuplicateKeyError(prev error, column *string) *errorPkg.Error {
	msg := "duplicate key error"
	if column != nil {
		msg = fmt.Sprintf("duplicate key in column %s", *column)
	}

	return errorPkg.New(errorPkg.CodeDuplicateKey, msg, prev)
}

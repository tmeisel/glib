package error

import (
	"net/http"
	"strconv"

	"github.com/tmeisel/glib/utils/strutils"
)

// Code is an int describing the type of error.
// It begins with a suitable http response code
// like 409 (hence, 3 digits) plus 2 more digits
// that specify the error. If the last 2 digits
// are 0s (00), it's a generic error of that type
//
// Example:
//
//	40900 is a generic conflict
//	40901 is a duplicate key error thrown by e.g. the database
type Code int

const (
	CodeUser               Code = 40000
	CodeAuthRequired       Code = 40100
	CodeForbidden          Code = 40300
	CodeNotFound           Code = 40400
	CodeConflict           Code = 40900
	CodeDuplicateKey       Code = 40901
	CodePreconditionFailed Code = 41200
	CodeGone               Code = 41000
	CodeTooManyRequests    Code = 42900
	CodeInternal           Code = 50000
)

func (c Code) String() string {
	return c.HttpStatusText()
}

// HttpStatus returns the http status code for the given Code
// That's the first 3 digit of code.
func (c Code) HttpStatus() int {
	asString := strutils.SubString(strconv.Itoa(int(c)), 0, 3)
	asInt, err := strconv.Atoi(asString)
	if err != nil {
		return http.StatusInternalServerError
	}

	return asInt
}

// HttpStatusText returns the http status text corresponding to
// the Code
func (c Code) HttpStatusText() string {
	return http.StatusText(c.HttpStatus())
}

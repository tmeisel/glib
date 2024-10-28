package error

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/tmeisel/glib/utils/strutils"
)

const (
	MaxStackDepth = 5
)

type Error struct {
	code     Code
	msg      string
	previous error
	stack    []uintptr
}

func New(code Code, msg string, prev error) error {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	return Error{
		code:     code,
		msg:      msg,
		previous: prev,
		stack:    stack[:length],
	}
}

func NewUser(prev error) error {
	return NewUserMsg(prev, statusText(CodeUser))
}

func NewUserMsg(prev error, msg string) error {
	return New(CodeUser, msg, prev)
}

func NewInternal(prev error) error {
	return NewInternalMsg(prev, statusText(CodeInternal))
}

func NewInternalMsg(prev error, msg string) error {
	return New(CodeInternal, msg, prev)
}

func (e Error) GetCode() Code {
	return e.code
}

func (e Error) GetStatus() int {
	return status(e.code)
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Message() string {
	return e.msg
}

func (e Error) Unwrap() error {
	if e.previous != nil {
		return e.previous
	}

	return nil
}

func (e Error) GetStack() []uintptr {
	return e.stack
}

func status(code Code) int {
	asString := strutils.SubString(strconv.Itoa(int(code)), 0, 3)
	asInt, _ := strconv.Atoi(asString)

	return asInt
}

// statusText returns the http status text corresponding to
// the given
func statusText(code Code) string {
	return http.StatusText(status(code))
}

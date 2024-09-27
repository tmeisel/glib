package error

import (
	"errors"
	"runtime"
)

const (
	MaxStackDepth = 5
)

type Error struct {
	code     int
	msg      string
	previous error
	stack    []uintptr
}

func New(code int, msg string, prev error) error {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	return Error{
		code:     code,
		msg:      msg,
		previous: prev,
		stack:    stack[:length],
	}
}

func (e Error) Error() string {
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

func GetStack(e error) []uintptr {
	if errors.Is(e, Error{}) {
		return e.(Error).GetStack()
	}

	return nil
}

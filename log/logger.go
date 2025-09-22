package log

import (
	"context"
	"strings"

	"github.com/tmeisel/glib/log/fields"
)

type Level int8

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
)

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
)

func LevelFromString(s string) Level {
	switch strings.ToLower(s) {
	case Debug:
		return LevelDebug
	case Info:
		return LevelInfo
	case Warn:
		return LevelWarn
	case Error:
		return LevelError
	default:
		return LevelInfo
	}
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...fields.Field)
	Info(ctx context.Context, msg string, fields ...fields.Field)
	Warn(ctx context.Context, msg string, fields ...fields.Field)
	Error(ctx context.Context, msg string, fields ...fields.Field)

	// Debugf writes a formatted log entry. Field elements can be appended
	// to args
	Debugf(ctx context.Context, format string, args ...interface{})

	// Infof writes a formatted log entry. Field elements can be appended
	// to args
	Infof(ctx context.Context, format string, args ...interface{})

	// Warnf writes a formatted log entry. Field elements can be appended
	// to args
	Warnf(ctx context.Context, format string, args ...interface{})

	// Errorf writes a formatted log entry. Field elements can be appended
	// to args
	Errorf(ctx context.Context, format string, args ...interface{})

	// Printf is a generic formatted print compatible with fmt.Printf.
	// It should be avoided but makes it compatible with other libraries
	Printf(format string, args ...interface{})

	// Write is a generic writer. It should be avoided but makes it compatible
	// with other libraries
	Write(p []byte) (bytesWritten int, err error)

	// SetLevel changes the loglevel to the given Level
	SetLevel(level Level) error

	// Shutdown must be called before the application exits. It
	// flushes all remaining messages to the writer
	Shutdown() error
}

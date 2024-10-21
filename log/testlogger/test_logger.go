package testlogger

import (
	"context"

	"github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/common"
	"github.com/tmeisel/glib/log/fields"
)

type TestLogger struct {
	lvl     log.Level
	entries []logEntry
}

type logEntry struct {
	Lvl    log.Level
	Msg    string
	Fields []fields.Field
}

func New(logLevel log.Level) *TestLogger {
	return &TestLogger{
		lvl:     logLevel,
		entries: make([]logEntry, 0),
	}
}

func (t *TestLogger) GetEntries() []logEntry {
	return t.entries
}

func (t *TestLogger) Debug(_ context.Context, msg string, fields ...fields.Field) {
	if t.lvl > log.LevelDebug {
		return
	}

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelDebug,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Info(_ context.Context, msg string, fields ...fields.Field) {
	if t.lvl > log.LevelInfo {
		return
	}

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelInfo,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Warn(_ context.Context, msg string, fields ...fields.Field) {
	if t.lvl > log.LevelWarn {
		return
	}

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelWarn,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Error(_ context.Context, msg string, fields ...fields.Field) {
	if t.lvl > log.LevelError {
		return
	}

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelError,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Debugf(_ context.Context, format string, args ...interface{}) {
	if t.lvl > log.LevelDebug {
		return
	}

	msg, fields := common.ProcessFormatted(format, args...)

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelDebug,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Infof(_ context.Context, format string, args ...interface{}) {
	if t.lvl > log.LevelInfo {
		return
	}

	msg, fields := common.ProcessFormatted(format, args...)

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelInfo,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Warnf(_ context.Context, format string, args ...interface{}) {
	if t.lvl > log.LevelWarn {
		return
	}

	msg, fields := common.ProcessFormatted(format, args...)

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelWarn,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) Errorf(_ context.Context, format string, args ...interface{}) {
	if t.lvl > log.LevelError {
		return
	}

	msg, fields := common.ProcessFormatted(format, args...)

	t.entries = append(t.entries, logEntry{
		Lvl:    log.LevelError,
		Msg:    msg,
		Fields: fields,
	})
}

func (t *TestLogger) SetLevel(l log.Level) error {
	t.lvl = l
	return nil
}

func (t *TestLogger) Shutdown() error {
	t.entries = make([]logEntry, 0)
	return nil
}

package writer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	ctxPkg "github.com/tmeisel/glib/ctx"
	logPkg "github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/fields"
)

// Writer is a log.Logger implementation, that writes all
// output to an io.Writer. It hence can be used to write
// log messages to os.Stdout
type Writer struct {
	writer     io.Writer
	level      logPkg.Level
	production bool
}

var _ logPkg.Logger = &Writer{}

func New(writer io.Writer, production bool, level logPkg.Level) *Writer {
	if writer == nil {
		panic("nil writer")
	}

	return &Writer{
		writer:     writer,
		production: production,
		level:      level,
	}
}

// NewStdWriter returns a Writer that uses io.Stdout as destination
func NewStdWriter(production bool, level logPkg.Level) *Writer {
	return New(os.Stdout, production, level)
}

func (w *Writer) Debug(ctx context.Context, msg string, fields ...fields.Field) {
	w.write(ctx, logPkg.LevelDebug, msg, fields...)
}

func (w *Writer) Debugf(ctx context.Context, format string, args ...interface{}) {
	w.Debug(ctx, fmt.Sprintf(format, args...))
}

func (w *Writer) Info(ctx context.Context, msg string, fields ...fields.Field) {
	w.write(ctx, logPkg.LevelInfo, msg, fields...)
}

func (w *Writer) Infof(ctx context.Context, format string, args ...interface{}) {
	w.Info(ctx, fmt.Sprintf(format, args...))
}

func (w *Writer) Warn(ctx context.Context, msg string, fields ...fields.Field) {
	w.write(ctx, logPkg.LevelWarn, msg, fields...)
}

func (w *Writer) Warnf(ctx context.Context, format string, args ...interface{}) {
	w.Warn(ctx, fmt.Sprintf(format, args...))
}

func (w *Writer) Error(ctx context.Context, msg string, fields ...fields.Field) {
	w.write(ctx, logPkg.LevelError, msg, fields...)
}

func (w *Writer) Errorf(ctx context.Context, format string, args ...interface{}) {
	w.Error(ctx, fmt.Sprintf(format, args...))
}

func (w *Writer) Printf(format string, args ...interface{}) {
	w.writef(context.Background(), logPkg.LevelInfo, format, args...)
}

func (w *Writer) writef(ctx context.Context, level logPkg.Level, format string, args ...interface{}) {
	fieldArgs := make([]fields.Field, 0)
	fmtArgs := make([]interface{}, 0)

	for _, arg := range args {
		if f, ok := arg.(fields.Field); ok {
			fieldArgs = append(fieldArgs, f)
		} else {
			fmtArgs = append(fmtArgs, arg)
		}
	}

	formattedMessage := fmt.Sprintf(format, fmtArgs...)

	w.write(ctx, level, formattedMessage, fieldArgs...)
}

func (w *Writer) write(ctx context.Context, level logPkg.Level, msg string, fields ...fields.Field) {
	if level < w.level {
		return
	}

	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	if w.production {
		timestamp = strconv.FormatInt(now.UnixMilli(), 10)
	}

	msg = fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), msg)

	fields = append(fields, ctxPkg.GetLogFields(ctx)...)
	for _, f := range fields {
		msg += fmt.Sprintf(" {%s: '%v'}", f.Key, fmt.Sprint(f))
	}

	w.writer.Write([]byte(msg))
}

func (w *Writer) Write(p []byte) (bytesWritten int, err error) {
	return w.writer.Write(p)
}

func (w *Writer) SetLevel(level logPkg.Level) error {
	w.level = level

	return nil
}

func (w *Writer) Shutdown() error {
	return nil
}

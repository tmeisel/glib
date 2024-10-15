package zap

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	logPkg "github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/fields"
)

type Zap struct {
	atom   zap.AtomicLevel
	logger *zap.Logger
}

var _ logPkg.Logger = &Zap{}

func New(production bool, level logPkg.Level, options ...zap.Option) *Zap {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(level))

	var encoderCfg zapcore.EncoderConfig
	if production {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)

	// do not change order. it allows overwriting
	// with passed options
	options = append([]zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	}, options...)

	return &Zap{atom: atom, logger: zap.New(core, options...)}
}

func (z *Zap) SetLevel(level logPkg.Level) error {
	z.atom.SetLevel(zapcore.Level(level))

	return nil
}

func (z *Zap) Debug(ctx context.Context, msg string, fields ...fields.Field) {
	z.log(ctx, zap.DebugLevel, msg, fields...)
}

func (z *Zap) Info(ctx context.Context, msg string, fields ...fields.Field) {
	z.log(ctx, zap.InfoLevel, msg, fields...)
}

func (z *Zap) Warn(ctx context.Context, msg string, fields ...fields.Field) {
	z.log(ctx, zap.WarnLevel, msg, fields...)
}

func (z *Zap) Error(ctx context.Context, msg string, fields ...fields.Field) {
	z.log(ctx, zap.ErrorLevel, msg, fields...)
}

func (z *Zap) Debugf(ctx context.Context, format string, args ...interface{}) {
	msg, fields := z.msg(format, args...)
	z.log(ctx, zap.DebugLevel, msg, fields...)
}

func (z *Zap) Infof(ctx context.Context, format string, args ...interface{}) {
	msg, fields := z.msg(format, args...)
	z.log(ctx, zap.InfoLevel, msg, fields...)
}

func (z *Zap) Warnf(ctx context.Context, format string, args ...interface{}) {
	msg, fields := z.msg(format, args...)
	z.log(ctx, zap.WarnLevel, msg, fields...)
}

func (z *Zap) Errorf(ctx context.Context, format string, args ...interface{}) {
	msg, fields := z.msg(format, args...)
	z.log(ctx, zap.ErrorLevel, msg, fields...)
}

func (z *Zap) Shutdown() error {
	return z.logger.Sync()
}

func (z *Zap) log(ctx context.Context, level zapcore.Level, msg string, fields ...fields.Field) {
	z.logger.Log(level, msg, z.fields(ctx, fields...)...)
}
func (z *Zap) fields(_ context.Context, fields ...fields.Field) []zap.Field {
	// @todo: join with fields from context
	return fields
}

func (z *Zap) msg(format string, args ...interface{}) (string, []fields.Field) {
	var fieldsOut []fields.Field
	for _, arg := range args {
		if field, ok := arg.(fields.Field); ok {
			fieldsOut = append(fieldsOut, field)
		}
	}

	fieldCount := len(fieldsOut)
	if fieldCount == 0 {
		return fmt.Sprintf(format, args...), nil
	}

	firstField := len(args) - fieldCount

	return fmt.Sprintf(format, args[:firstField]...), fieldsOut
}

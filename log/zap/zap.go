package zap

import (
	"context"
	"os"

	"github.com/tmeisel/glib/log/common"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	ctxPkg "github.com/tmeisel/glib/ctx"
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
	msg, fields := common.ProcessFormatted(format, args...)
	z.log(ctx, zap.DebugLevel, msg, fields...)
}

func (z *Zap) Infof(ctx context.Context, format string, args ...interface{}) {
	msg, fields := common.ProcessFormatted(format, args...)
	z.log(ctx, zap.InfoLevel, msg, fields...)
}

func (z *Zap) Warnf(ctx context.Context, format string, args ...interface{}) {
	msg, fields := common.ProcessFormatted(format, args...)
	z.log(ctx, zap.WarnLevel, msg, fields...)
}

func (z *Zap) Errorf(ctx context.Context, format string, args ...interface{}) {
	msg, fields := common.ProcessFormatted(format, args...)
	z.log(ctx, zap.ErrorLevel, msg, fields...)
}

func (z *Zap) Shutdown() error {
	return z.logger.Sync()
}

func (z *Zap) log(ctx context.Context, level zapcore.Level, msg string, fields ...fields.Field) {
	f := common.JoinUnique(ctxPkg.GetUniqueLogFields(ctx), fields...)

	z.logger.Log(level, msg, f...)
}

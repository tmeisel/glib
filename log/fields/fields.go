package fields

import (
	"time"

	"go.uber.org/zap"
)

type Field = zap.Field

func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func Error(e error) Field {
	return zap.Error(e)
}

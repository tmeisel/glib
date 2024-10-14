package fields

import "go.uber.org/zap"

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

func Error(e error) Field {
	return zap.Error(e)
}

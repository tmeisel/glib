package ctx

import "context"

type configContextKey struct {
	key string
}

func GetConfigValue(ctx context.Context, key string) string {
	val := ctx.Value(configContextKey{key: key})
	if val == nil {
		return ""
	}

	return val.(string)
}

func WithConfigValue(parent context.Context, key, value string) context.Context {
	return context.WithValue(parent, configContextKey{key: key}, value)
}

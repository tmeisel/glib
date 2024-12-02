package ctx

import (
	"context"
)

type authContextKeyType struct{}

func GetIdentity(ctx context.Context) any {
	val := ctx.Value(authContextKeyType{})
	if val == nil {
		return nil
	}

	return val
}

func WithIdentity(parent context.Context, identity any) context.Context {
	return context.WithValue(parent, authContextKeyType{}, identity)
}

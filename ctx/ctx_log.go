package ctx

import (
	"context"

	"github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/fields"
)

const (
	ctxKeyLogger = contextKey("logger")
	ctxKeyFields = contextKey("fields")
)

type contextKey string

func WithLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func GetLogger(ctx context.Context) log.Logger {
	value := ctx.Value(ctxKeyLogger)
	if value == nil {
		return nil
	}

	return value.(log.Logger)
}

// WithLogFields adds the given fields.Field f to the context
func WithLogFields(ctx context.Context, f ...fields.Field) context.Context {
	current := ctx.Value(ctxKeyFields)
	if current == nil {
		return context.WithValue(ctx, ctxKeyFields, f)
	}

	return context.WithValue(ctx, ctxKeyFields, append(current.([]fields.Field), f...))
}

// GetLogFields returns all logged fields
func GetLogFields(ctx context.Context) []fields.Field {
	current := ctx.Value(ctxKeyFields)
	if current == nil {
		return make([]fields.Field, 0)
	}

	return current.([]fields.Field)
}

// GetUniqueLogFields returns all logged fields with a unique key. If two or more
// fields have the same key, the last one added will be returned
func GetUniqueLogFields(ctx context.Context) []fields.Field {
	current := ctx.Value(ctxKeyFields)
	if current == nil {
		return make([]fields.Field, 0)
	}

	var output []fields.Field
	keys := make(map[string]int)
	for _, field := range current.([]fields.Field) {
		idx, exists := keys[field.Key]
		if exists {
			output[idx] = field
			continue
		}

		keys[field.Key] = len(output)
		output = append(output, field)
	}

	if len(output) == 0 {
		return make([]fields.Field, 0)
	}

	return output
}

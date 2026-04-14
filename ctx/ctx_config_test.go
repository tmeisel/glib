package ctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tmeisel/glib/utils/strutils"
)

func TestWithConfigValue(t *testing.T) {
	ctx := WithConfigValue(context.Background(), "key", "value1")
	value := ctx.Value(configContextKey{key: "key"})
	require.NotNil(t, value)
	require.Equal(t, "value1", value.(string))
}

func TestGetConfigValue(t *testing.T) {
	for name, key := range map[string]string{
		"simple key": "key1",
		"dotted key": "key1.subkey1",
	} {
		t.Run(name, func(t *testing.T) {
			value := strutils.MustRandom(5, strutils.AlphabetReadable)
			ctx := WithConfigValue(context.Background(), key, value)
			assert.Equal(t, value, GetConfigValue(ctx, key))
		})
	}

}

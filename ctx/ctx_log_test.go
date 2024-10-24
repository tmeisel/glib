package ctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logPkg "github.com/tmeisel/glib/log"
	"github.com/tmeisel/glib/log/fields"
	"github.com/tmeisel/glib/log/testlogger"
)

func TestWithLogger(t *testing.T) {
	input := testlogger.New(logPkg.LevelDebug)

	ctx := context.Background()
	require.Nil(t, GetLogger(ctx))

	ctx = WithLogger(ctx, input)
	output := GetLogger(ctx)

	require.NotNil(t, output)
	assert.Implements(t, (*logPkg.Logger)(nil), output)

}

func TestWithLogField(t *testing.T) {
	type testCase struct {
		Input    []fields.Field
		Expected []fields.Field
	}

	for name, tc := range map[string]testCase{
		"none": {
			Input:    nil,
			Expected: []fields.Field{},
		},
		"one": {
			Input:    []fields.Field{fields.String("userID", "ab1839")},
			Expected: []fields.Field{fields.String("userID", "ab1839")},
		},
		"two": {
			Input: []fields.Field{
				fields.String("userID", "ab1839"),
				fields.Int("ID", 1839),
			},
			Expected: []fields.Field{
				fields.String("userID", "ab1839"),
				fields.Int("ID", 1839),
			},
		},
		"duplicate key": {
			Input: []fields.Field{
				fields.String("userID", "ab1839"),
				fields.Int("ID", 1839),
				fields.String("userID", "xy9876"),
			},
			Expected: []fields.Field{
				fields.String("userID", "ab1839"),
				fields.Int("ID", 1839),
				fields.String("userID", "xy9876"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			ctx = WithLogFields(ctx, tc.Input...)

			assert.Equal(t, tc.Expected, GetLogFields(ctx))
		})
	}
}

func TestGetUniqueLogFields(t *testing.T) {
	type testCase struct {
		Input    []fields.Field
		Expected []fields.Field
	}

	for name, tc := range map[string]testCase{
		"none": {
			Input:    nil,
			Expected: []fields.Field{},
		},
		"one": {
			Input:    []fields.Field{fields.String("userID", "ab1839")},
			Expected: []fields.Field{fields.String("userID", "ab1839")},
		},
		"two": {
			Input: []fields.Field{
				fields.String("userID", "ab1839"),
				fields.String("userID", "ab1840"),
			},
			Expected: []fields.Field{
				fields.String("userID", "ab1840"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			ctx = WithLogFields(ctx, tc.Input...)

			assert.Equal(t, tc.Expected, GetUniqueLogFields(ctx))
		})
	}
}

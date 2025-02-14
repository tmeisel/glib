package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromCtx(t *testing.T) {
	ctx := Disconnect(context.Background())

	require.NotNil(t, ctx)
	assert.Implements(t, (*context.Context)(nil), ctx)
}

func TestDisconnected_Value(t *testing.T) {
	const key = "key"
	var value any = "test"

	parent := context.Background()
	parent = context.WithValue(parent, key, value)

	ctx := Disconnect(parent)

	returnVal := ctx.Value(key)

	require.NotNil(t, returnVal)
	assert.Equal(t, value, returnVal)
}

func TestDisconnected_Deadline(t *testing.T) {
	parent, cancelParent := context.WithDeadline(context.Background(), time.Now().Add(time.Microsecond))
	ctx := Disconnect(parent)

	cancelParent()
	assert.Error(t, parent.Err())

	deadline, ok := ctx.Deadline()
	assert.Equal(t, time.Time{}, deadline)
	assert.False(t, ok)
}

func TestDisconnected_Err(t *testing.T) {
	parent, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	ctx := Disconnect(parent)

	assert.Error(t, parent.Err())
	assert.NoError(t, ctx.Err())
}

func TestDisconnected_Done(t *testing.T) {
	parent, cancelFn := context.WithCancel(context.Background())

	ctx := Disconnect(parent)
	cancelFn()

	select {
	case <-ctx.Done():
		t.Error(ctx.Err())
	default:
	}

	assert.Error(t, parent.Err())
}

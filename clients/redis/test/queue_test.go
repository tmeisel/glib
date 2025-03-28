package test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	queuePkg "github.com/tmeisel/glib/queue"
)

func TestQueue(t *testing.T) {
	queue := client.Queue(uuid.NewString(), time.Second)
	require.Implements(t, (*queuePkg.Queue)(nil), queue)
}

func TestLPop(t *testing.T) {
	queue := client.Queue(uuid.NewString(), time.Second)

	// empty queue
	_, err := queue.RPop()
	require.ErrorIs(t, err, queuePkg.ErrEmpty)

	// push
	require.NoError(t, queue.LPush("hello"))

	// pop
	val, err := queue.LPop()
	require.NoError(t, err)
	assert.Equal(t, "hello", val)
}

func TestRPop(t *testing.T) {
	queue := client.Queue(uuid.NewString(), time.Second)

	// empty queue
	_, err := queue.RPop()
	require.ErrorIs(t, err, queuePkg.ErrEmpty)

	// push
	require.NoError(t, queue.RPush("hello"))

	// pop
	val, err := queue.RPop()
	require.NoError(t, err)
	assert.Equal(t, "hello", val)
}

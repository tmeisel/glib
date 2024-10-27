package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	queuePkg "github.com/tmeisel/glib/queue"
)

func TestQueue(t *testing.T) {
	queue := client.Queue()
	require.Implements(t, (*queuePkg.Queue)(nil), queue)
}

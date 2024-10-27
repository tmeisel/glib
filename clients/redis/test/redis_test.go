package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedis_Ping(t *testing.T) {
	require.NoError(t, client.Ping())
}

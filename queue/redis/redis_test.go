package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/tmeisel/glib/clients/redis/docker"
)

var (
	container *docker.Container
)

func TestMain(m *testing.M) {
	var err error
	container, err = docker.NewContainer(context.Background(), docker.DefaultVersion, time.Second*30, time.Second*10)
	if err != nil {
		os.Exit(1)
	}
	defer container.Close()

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	const name = "testQueue"

	client := New(container.GetConfig(), name)
	assert.Equal(t, name, client.Name())

	const key = "testQueue"
	const value = "testValue"

	require.NoError(t, client.LPush(key, value))
}

func TestNewFromClient(t *testing.T) {
	const name = "testQueue"

	client := NewFromClient(container.GetClient(), name)
	assert.Equal(t, name, client.Name())

	const key = "testQueue"
	const value = "testValue"

	require.NoError(t, client.LPush(key, value))
}

package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tmeisel/glib/clients/redis"
	"github.com/tmeisel/glib/clients/redis/docker"
)

var (
	container *docker.Container
	client    redis.Redis
)

func TestMain(m *testing.M) {
	var err error
	container, err = docker.NewContainer(context.Background(), docker.DefaultVersion, time.Second*30, time.Second*10)
	if err != nil {
		os.Exit(1)
	}

	client = redis.New(container.GetConfig())

	m.Run()
}

func TestRedis_Ping(t *testing.T) {
	require.NoError(t, client.Ping())
}

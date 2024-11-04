package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tmeisel/glib/clients/redis/docker"
	"github.com/tmeisel/glib/queue"
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
	client := New(container.GetConfig(), uuid.NewString(), time.Second)

	assert.Implements(t, (*queue.Queue)(nil), client)
}

func TestNewFromClient(t *testing.T) {
	client := NewFromClient(container.GetClient(), uuid.NewString(), time.Second)

	assert.Implements(t, (*queue.Queue)(nil), client)
}

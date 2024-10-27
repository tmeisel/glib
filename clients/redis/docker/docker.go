package docker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	redisPkg "github.com/tmeisel/glib/clients/redis"

	"github.com/tmeisel/glib/exec/backoff"
	dockerPkg "github.com/tmeisel/glib/testing/docker"
)

const (
	repo = "redis"
	port = 6379
)

const (
	DefaultVersion = V7
	V7             = "7-alpine"
	V6             = "6-alpine"
)

type Container struct {
	*dockerPkg.Container

	exposedPort uint
}

// NewContainer runs a redis server in a docker container. The maxContainerLifetime must
// be set according to the tests it is used for. The waitTime specifies how long you want to
// wait for the client to successfully connect to the redis server AFTER the container
// started. Hence, downloading the docker image and starting the container may exceed the
// given waitTime.
// If the waitTime is eventually exceeded, an error will be returned.
// version needs to be a valid tag, e.g. (redis:)7.4-alpine.
//
// The caller should eventually call Cleanup when tests are finished (usually in TestMain)
func NewContainer(ctx context.Context, version string, maxContainerLifetime, waitTime time.Duration) (*Container, error) {
	t := &Container{}

	if err := t.initContainer(version, maxContainerLifetime); err != nil {
		return nil, err
	}

	if err := t.waitForContainer(ctx, waitTime); err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Container) Close() error {
	return c.Container.Cleanup()
}

func (c *Container) GetClient() redisPkg.Redis {
	return redisPkg.New(c.GetConfig())
}

func (c *Container) GetAddresses() []string {
	return []string{fmt.Sprintf("localhost:%d", c.exposedPort)}
}

func (c *Container) GetConfig() redisPkg.Config {
	return redisPkg.Config{Addresses: c.GetAddresses()}
}

func (c *Container) initContainer(version string, maxLifetime time.Duration) error {
	var err error

	// create the container pool
	c.Container, err = dockerPkg.New(repo, version, maxLifetime)
	if err != nil {
		return err
	}

	portString := c.Container.GetPortString(fmt.Sprintf("%d/tcp", port))

	port, err := strconv.Atoi(portString)
	if err != nil {
		return err
	}

	c.exposedPort = uint(port)

	return err
}

func (c *Container) waitForContainer(ctx context.Context, maxWaitTime time.Duration) error {
	client := redisPkg.New(c.GetConfig())

	connectFn := func(_ context.Context) error {
		if err := client.Ping(); err != nil {
			return backoff.RetryableError(err)
		}

		return nil
	}

	retry, err := backoff.New(
		backoff.Fibonacci,
		time.Millisecond*50,
		backoff.WithMaxDuration(maxWaitTime),
		backoff.WithCap(time.Millisecond*500),
	)

	if err != nil {
		// backoff init failed, try direct connect
		return connectFn(ctx)
	}

	return retry.Do(ctx, connectFn)
}

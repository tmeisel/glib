package docker

import (
	"fmt"
	"time"

	"github.com/ory/dockertest"
)

// Container is the base implementation of test containers
// using github.com/ory/dockertest
type Container struct {
	ports map[string]string

	pool     *dockertest.Pool
	resource *dockertest.Resource
}

// New creates a new pool and runs the given repo/tag. The pool will expire
// automatically after the given expiration time.Duration
func New(repo, tag string, expiration time.Duration, env ...string) (*Container, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %w", err)
	}

	resource, err := pool.Run(repo, tag, env)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	if err = resource.Expire(uint(expiration / time.Second)); err != nil {
		return nil, err
	}

	return &Container{
		ports:    make(map[string]string),
		pool:     pool,
		resource: resource,
	}, nil
}

func (l *Container) GetPool() *dockertest.Pool {
	return l.pool
}

func (l *Container) GetResource() *dockertest.Resource {
	return l.resource
}

// GetPortString returns the container port the service was
// published at. servicePort must be passed in the form
// port/protocol, e.g. "6379/tcp"
func (l *Container) GetPortString(servicePort string) string {
	if _, known := l.ports[servicePort]; !known {
		l.ports[servicePort] = l.resource.GetPort(servicePort)
	}

	return l.ports[servicePort]
}

func (l *Container) Cleanup() error {
	return l.resource.Close()
}

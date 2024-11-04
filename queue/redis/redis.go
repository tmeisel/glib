package redis

import (
	"time"

	redisPkg "github.com/tmeisel/glib/clients/redis"
)

// Queue is a wrapper for redisPkg.Queue
type Queue = redisPkg.Queue

func New(conf redisPkg.Config, name string, readTimeout time.Duration) Queue {
	r := redisPkg.New(conf)
	return r.Queue(name, readTimeout)
}

func NewFromClient(r redisPkg.Redis, name string, readTimeout time.Duration) Queue {
	return r.Queue(name, readTimeout)
}

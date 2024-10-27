package redis

import redisPkg "github.com/tmeisel/glib/clients/redis"

// Queue is a wrapper for redisPkg.Queue
type Queue = redisPkg.Queue

func New(conf redisPkg.Config, name string) Queue {
	r := redisPkg.New(conf)
	return r.Queue(name)
}

func NewFromClient(r redisPkg.Redis, name string) Queue {
	return r.Queue(name)
}

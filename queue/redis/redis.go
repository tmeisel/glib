package redis

import redisPkg "github.com/tmeisel/glib/clients/redis"

// Queue is a wrapper for redisPkg.Queue
type Queue = redisPkg.Queue

func New(conf redisPkg.Config) Queue {
	r := redisPkg.New(conf)
	return r.Queue()
}

func NewFromClient(r redisPkg.Redis) Queue {
	return r.Queue()
}

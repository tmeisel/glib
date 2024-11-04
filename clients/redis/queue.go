package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	errPkg "github.com/tmeisel/glib/error"
	queuePkg "github.com/tmeisel/glib/queue"
)

type Queue struct {
	r           *Redis
	name        string
	readTimeout time.Duration
}

// Queue initializes the client to connect to the queue with the given name.
// readTimeout must be at least 1 second. If the given value is smaller,
// it will be overwritten with time.Second
func (r Redis) Queue(name string, readTimeout time.Duration) Queue {
	if readTimeout < time.Second {
		readTimeout = time.Second
	}

	return Queue{
		r:           &r,
		name:        fmt.Sprintf("queue-%v", name),
		readTimeout: readTimeout,
	}
}

func (q Queue) Name() string {
	return q.name
}

func (q Queue) Empty() error {
	if err := q.r.client.Del(q.name).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to empty list")
	}

	return nil
}

func (q Queue) Push(value string) error {
	return q.LPush(value)
}

func (q Queue) Pop() (string, error) {
	return q.RPop()
}

func (q Queue) LPush(value string) error {
	if err := q.r.client.LPush(q.name, value).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to push value")
	}

	return nil
}

func (q Queue) RPush(value string) error {
	if err := q.r.client.RPush(q.name, value).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to push value")
	}

	return nil
}

func (q Queue) LPop() (string, error) {
	res := q.r.client.BLPop(q.readTimeout, q.name)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", queuePkg.ErrEmpty
		}

		return "", err
	}

	vals := res.Val()
	if len(vals) != 2 {
		return "", errPkg.NewInternalMsg(nil, "unexpected response from redis")
	}

	// vals[0] is the name of the queue (aka the key)
	// vals[1] is the actual value
	return vals[1], nil
}

func (q Queue) RPop() (string, error) {
	res := q.r.client.BRPop(q.readTimeout, q.name)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", queuePkg.ErrEmpty
		}

		return "", err
	}

	vals := res.Val()
	if len(vals) != 2 {
		return "", errPkg.NewInternalMsg(nil, "unexpected response from redis")
	}

	// vals[0] is the name of the queue (aka the key)
	// vals[1] is the actual value
	return vals[1], nil
}

package redis

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	errPkg "github.com/tmeisel/glib/error"
	queuePkg "github.com/tmeisel/glib/queue"
)

type Queue struct {
	r    *Redis
	name string
}

func (r Redis) Queue() Queue {
	return Queue{
		r:    &r,
		name: fmt.Sprintf("queue-%v", uuid.NewString()),
	}
}

func (q Queue) Empty() error {
	if err := q.r.client.Del(q.name).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to empty list")
	}

	return nil
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
	res := q.r.client.LPop(q.name)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", queuePkg.ErrEmpty
		}

		return "", err
	}

	return res.Val(), nil
}

func (q Queue) RPop() (string, error) {
	res := q.r.client.RPop(q.name)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", queuePkg.ErrEmpty
		}

		return "", err
	}

	return res.Val(), nil
}

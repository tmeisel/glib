package redis

import errPkg "github.com/tmeisel/glib/error"

type Queue struct {
	r    *Redis
	name string
}

func (r Redis) Queue(name string) Queue {
	return Queue{
		r:    &r,
		name: name,
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

func (q Queue) LPush(key, value string) error {
	if err := q.r.client.LPush(q.name, key, value).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to push value")
	}

	return nil
}

func (q Queue) RPush(key, value string) error {
	if err := q.r.client.RPush(q.name, key, value).Err(); err != nil {
		return errPkg.NewInternalMsg(err, "failed to push value")
	}

	return nil
}

func (q Queue) LPop(key string) (string, error) {
	res := q.r.client.LPop(key)

	if err := res.Err(); err != nil {
		return "", errPkg.NewInternalMsg(err, "failed to pop value")
	}

	return res.String(), nil
}

func (q Queue) RPop(key string) (string, error) {
	res := q.r.client.RPop(key)

	if err := res.Err(); err != nil {
		return "", errPkg.NewInternalMsg(err, "failed to pop value")
	}

	return res.String(), nil
}

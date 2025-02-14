package ctx

import (
	"context"
	"time"
)

type disconnected struct {
	parent context.Context
}

var _ context.Context = disconnected{}

func Disconnect(ctx context.Context) context.Context {
	return &disconnected{parent: ctx}
}

func (c disconnected) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c disconnected) Done() <-chan struct{} {
	return nil
}

func (c disconnected) Err() error {
	return nil
}

func (c disconnected) Value(key any) any {
	return c.parent.Value(key)
}

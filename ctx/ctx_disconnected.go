package ctx

import (
	"context"
	"time"
)

type Disconnected struct {
	parent context.Context
}

var _ context.Context = Disconnected{}

func FromCtx(ctx context.Context) *Disconnected {
	return &Disconnected{parent: ctx}
}

func (c Disconnected) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c Disconnected) Done() <-chan struct{} {
	return nil
}

func (c Disconnected) Err() error {
	return nil
}

func (c Disconnected) Value(key any) any {
	return c.parent.Value(key)
}

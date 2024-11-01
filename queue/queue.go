package queue

import "errors"

var (
	ErrEmpty = errors.New("queue is empty")
)

type Queue interface {
	// LPush adds an element at the head of the queue
	LPush(value string) error
	// RPush adds an element at the tail of the queue
	RPush(value string) error
	// LPop removes and returns an element from the head of the queue
	LPop() (string, error)
	// RPop removes and returns an element from the tail of the queue
	RPop() (string, error)
	// Empty removes all elements from the queue
	Empty() error
}

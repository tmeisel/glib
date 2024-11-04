package queue

import "errors"

var (
	ErrEmpty = errors.New("queue is empty")
)

type Queue interface {
	// Push add an element to the queue
	Push(value string) error
	// Pop returns a single element from the queue
	Pop() (string, error)
	// Empty removes all elements from the queue
	Empty() error
}

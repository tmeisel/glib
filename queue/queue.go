package queue

type Queue interface {
	// LPush adds an element at the head of the queue
	LPush(key, value string) error
	// RPush adds an element at the tail of the queue
	RPush(key, value string) error
	// LPop removes and returns an element from the head of the queue
	LPop(key string) (string, error)
	// RPop removes and returns an element from the tail of the queue
	RPop(key string) (string, error)
	// Empty removes all elements from the queue
	Empty() error
}

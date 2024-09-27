package exec

import (
	"log"
)

func Deferred(fn func() error) {
	if err := fn(); err != nil {
		log.Printf("deferred fn: %v", err)
	}
}

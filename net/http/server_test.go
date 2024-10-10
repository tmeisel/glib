package http

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	const (
		addr = "localhost"
		port = 80
	)

	s := NewServer(addr, port)
	assert.IsType(t, &Server{}, s)

	parts := strings.Split(s.srv.Addr, ":")

	assert.Equal(t, addr, parts[0], "invalid server address")
	assert.Equal(t, fmt.Sprintf("%d", port), parts[1], "invalid server port")
}

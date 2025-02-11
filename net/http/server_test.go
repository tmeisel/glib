package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	addr = "localhost"
	port = 8000
)

func TestNewServer(t *testing.T) {
	s := NewServer(context.Background(), addr, port)
	assert.IsType(t, &Server{}, s)

	parts := strings.Split(s.srv.Addr, ":")

	assert.Equal(t, addr, parts[0], "invalid server address")
	assert.Equal(t, fmt.Sprintf("%d", port), parts[1], "invalid server port")
}

func TestServer_ListenAndServe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewServer(ctx, addr, port)

	handleFn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	s.AddRoute(http.MethodGet, "/", handleFn)

	var serverError error
	go func() {
		serverError = s.ListenAndServe()
	}()

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/", addr, port))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	cancel()

	require.NoError(t, serverError)
}

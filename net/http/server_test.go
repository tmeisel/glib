package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rs/cors"
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
	require.NoError(t, s.Shutdown(ctx))
}

func TestServer_WithCORS(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	s := NewServerFromConf(ctx, ServerConfig{
		ListenAddr: addr,
		ListenPort: port + 1,
		WithCORS:   true,
		CORSOptions: &cors.Options{
			OptionsSuccessStatus: http.StatusOK,
			AllowedMethods:       []string{http.MethodDelete},
			AllowedOrigins:       []string{"http://localhost", "http://127.0.0.1"},
		},
	})

	handleFn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	s.AddRoute(http.MethodGet, "/", handleFn)

	var serverError error
	go func() {
		require.NoError(t, s.ListenAndServe())
	}()

	// wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	requestMethod := http.MethodDelete

	req, err := http.NewRequest(http.MethodOptions, fmt.Sprintf("http://%s:%d/", addr, port+1), nil)
	req.Header.Set("Access-Control-Request-Method", requestMethod)
	req.Header.Set("Access-Control-Request-Headers", "x-requested-with")
	req.Header.Set("Origin", "http://localhost")
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Origin"), "http://localhost")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), requestMethod)

	cancel()

	require.NoError(t, serverError)
}

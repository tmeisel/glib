package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	certFile string
	keyFile  string

	mwf    []mux.MiddlewareFunc
	router *mux.Router
	srv    http.Server
}

type ServerConfig struct {
	ListenAddr   string         `envconfig:"LISTEN_ADDR" default:"127.0.0.1"`
	ListenPort   uint           `envconfig:"LISTEN_PORT" default:"8080"`
	ReadTimeout  *time.Duration `envconfig:"READ_TIMEOUT" default:"20s"`
	WriteTimeout *time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	IdleTimeout  *time.Duration `envconfig:"IDLE_TIMEOUT" default:"10s"`
	CertFile     string         `envconfig:"CERT_FILE"`
	KeyFile      string         `envconfig:"KEY_FILE"`
}

func NewServer(ctx context.Context, addr string, port uint) *Server {
	return &Server{
		router: mux.NewRouter(),
		srv: http.Server{
			Addr: fmt.Sprintf("%s:%d", addr, port),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
	}
}

func NewServerFromConf(ctx context.Context, conf ServerConfig) *Server {
	s := NewServer(ctx, conf.ListenAddr, conf.ListenPort)

	if conf.ReadTimeout != nil {
		s.SetReadTimeout(*conf.ReadTimeout)
	}

	if conf.WriteTimeout != nil {
		s.SetWriteTimeout(*conf.WriteTimeout)
	}

	if conf.IdleTimeout != nil {
		s.SetIdleTimeout(*conf.IdleTimeout)
	}

	if conf.CertFile != "" && conf.KeyFile != "" {
		s.WithTLS(conf.CertFile, conf.KeyFile)
	}

	return s
}

func (s *Server) WithTLS(certFile, keyFile string) {
	s.certFile = certFile
	s.keyFile = keyFile
}

func (s *Server) SetReadTimeout(t time.Duration) {
	s.srv.ReadTimeout = t
}

func (s *Server) SetWriteTimeout(t time.Duration) {
	s.srv.WriteTimeout = t
}

func (s *Server) SetIdleTimeout(t time.Duration) {
	s.srv.IdleTimeout = t
}

func (s *Server) StrictSlash(value bool) {
	s.router.StrictSlash(value)
}

func (s *Server) SetErrorLog(l *log.Logger) {
	s.srv.ErrorLog = l
}

func (s *Server) AddRoute(method, path string, handler http.HandlerFunc) {
	s.router.HandleFunc(path, handler).Methods(method)
}

func (s *Server) PathHandler(pfx string, handler http.HandlerFunc) {
	s.router.PathPrefix(pfx).HandlerFunc(handler)
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

type MiddlewareFunc = func(http.Handler) http.Handler

func (s *Server) Use(fn ...MiddlewareFunc) {
	for _, mwf := range fn {
		s.mwf = append(s.mwf, mwf)
	}
}

func (s *Server) ListenAndServe() error {
	router := s.router
	router.Use(s.mwf...)

	s.srv.Handler = handlers.LoggingHandler(os.Stdout, router)

	if s.certFile != "" && s.keyFile != "" {
		return s.srv.ListenAndServeTLS(s.certFile, s.keyFile)
	}

	return s.srv.ListenAndServe()
}

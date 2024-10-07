package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	srv    http.Server
}

type ServerConfig struct {
	ListenAddr   string         `envvar:"LISTEN_ADDR" default:"127.0.0.1"`
	ListenPort   uint           `envvar:"LISTEN_PORT" default:"8080"`
	ReadTimeout  *time.Duration `envvar:"READ_TIMEOUT" default:"20s"`
	WriteTimeout *time.Duration `envvar:"WRITE_TIMEOUT" default:"5s"`
	IdleTimeout  *time.Duration `envvar:"IDLE_TIMEOUT" default:"10s"`
}

func NewServer(addr string, port uint) *Server {
	return &Server{
		router: mux.NewRouter(),
		srv: http.Server{
			Addr: fmt.Sprintf("%s:%d", addr, port),
		},
	}
}

func NewServerFromConf(conf ServerConfig) *Server {
	s := NewServer(conf.ListenAddr, conf.ListenPort)

	if conf.ReadTimeout != nil {
		s.SetReadTimeout(*conf.ReadTimeout)
	}

	if conf.WriteTimeout != nil {
		s.SetWriteTimeout(*conf.WriteTimeout)
	}

	if conf.IdleTimeout != nil {
		s.SetIdleTimeout(*conf.IdleTimeout)
	}

	return s
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

func (s *Server) ListenAndServe() error {
	s.srv.Handler = s.router

	return s.srv.ListenAndServe()
}

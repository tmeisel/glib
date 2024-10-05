package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	router *mux.Router
	srv    http.Server
}

type ServerConfig struct {
	ListAddr     string         `envvar:"LISTEN_ADDR"`
	ListPort     uint           `envvar:"LISTEN_PORT"`
	ReadTimeout  *time.Duration `envvar:"READ_TIMEOUT"`
	WriteTimeout *time.Duration `envvar:"WRITE_TIMEOUT"`
	IdleTimeout  *time.Duration `envvar:"IDLE_TIMEOUT"`
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
	s := NewServer(conf.ListAddr, conf.ListPort)

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

func (s *Server) SetErrorLog(l *log.Logger) {
	s.srv.ErrorLog = l
}

func (s *Server) AddRoute(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(path, handler).Methods(method)
}

func (s *Server) ListenAndServe() error {
	s.srv.Handler = s.router

	return s.srv.ListenAndServe()
}

package application

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Router http.Handler

type Server struct {
	l   *slog.Logger
	cfg ServerConfiguration
	r   Router
}

func NewServer(logger *slog.Logger, cfg ServerConfiguration, r Router) *Server {
	if logger == nil {
		logger = NewDefaultJsonLogger()
	}

	return &Server{
		l:   logger,
		cfg: cfg,
		r:   r,
	}
}

func (s *Server) Run() error {
	s.l.Info("Server is running", slog.String("port", s.cfg.Port))

	http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), s.r)

	return nil
}

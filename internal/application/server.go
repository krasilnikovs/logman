package application

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Router is alias to http.Handler
type Router http.Handler

// ApiServer contains main things related to incoming input requests
type ApiServer struct {
	l   *slog.Logger
	cfg ApiServerConfiguration
	r   Router
}

// NewApiServer constructs ApiServer
func NewApiServer(logger *slog.Logger, cfg ApiServerConfiguration, r Router) *ApiServer {
	if logger == nil {
		logger = NewDefaultJsonLogger()
	}

	return &ApiServer{
		l:   logger,
		cfg: cfg,
		r:   r,
	}
}

// Run method launch ApiServer, in case of some error it will return error
func (s *ApiServer) Run() error {
	s.l.Info("Server is running", slog.String("port", s.cfg.Port))

	return http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), s.r)
}

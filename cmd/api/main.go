// Package main in cmd/api directory contains codes which initilize api server, read configurations, create services.
package main

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/krasilnikovm/logman/internal/application"
	"github.com/krasilnikovm/logman/internal/handler"
)

// main is an entrypoint of laucnhing api server
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	registerRoutes(r)

	configuration := application.ServerConfiguration{}
	logger := application.NewDefaultJsonLogger()

	if err := cleanenv.ReadEnv(&configuration); err != nil {
		logger.Error("can not read envs", slog.String("error", err.Error()))
	}

	s := application.NewApiServer(logger, configuration, r)

	if err := s.Run(); err != nil {
		logger.Error("Server is not laucnhed")
	}
}

// registerRoutes method initilized routes
func registerRoutes(r *chi.Mux) {
	r.Get("/", handler.Index)
}

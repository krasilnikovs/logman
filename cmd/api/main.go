// Package main in cmd/api directory contains codes which initilize api server, read configurations, create services.
package main

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/krasilnikovm/logman/internal/application"
	"github.com/krasilnikovm/logman/internal/handler"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

const (
	driverName = "sqlite3"
)

// main is an entrypoint of launching api server
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	registerRoutes(r)

	configuration := application.ApiServerConfiguration{}
	logger := application.NewDefaultJsonLogger()

	if err := cleanenv.ReadEnv(&configuration); err != nil {
		logger.Error("can not read envs", slog.String("error", err.Error()))
	}

	if err := runMigrations(configuration); err != nil {
		logger.Error("migrations is not executed", slog.String("error", err.Error()))
	}

	s := application.NewApiServer(logger, configuration, r)

	if err := s.Run(); err != nil {
		logger.Error("Server is not launched")
	}
}

// registerRoutes method initialized routes
func registerRoutes(r *chi.Mux) {
	r.Get("/", handler.Index)
}

func runMigrations(cfg application.ApiServerConfiguration) error {

	db, err := sql.Open(driverName, cfg.DataStoragePath)

	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		return err
	}

	mgr, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		driverName,
		driver,
	)

	if err != nil {
		return err
	}

	if err = mgr.Up(); err != nil {
		return err
	}

	return nil
}

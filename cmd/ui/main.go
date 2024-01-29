// Package main in cmd/api directory contains codes which initialize api server, read configurations, create services.
package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/krasilnikovm/logman/internal/application"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// main is an entrypoint of launching api server
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	configuration := application.UiServerConfiguration{}
	logger := application.NewDefaultJsonLogger()

	if err := cleanenv.ReadEnv(&configuration); err != nil {
		logger.Error("can not read envs", slog.String("error", err.Error()))
	}

	registerRoutes(r, configuration, logger)

	s := application.NewUiServer(logger, configuration, r)

	if err := s.Run(); err != nil {
		logger.Error("UI Server is not launched")
	}
}

// registerRoutes method initialized routes
func registerRoutes(r *chi.Mux, cfg application.UiServerConfiguration, logger *slog.Logger) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		tmp, err := template.ParseFiles("web/template/index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmp.Execute(w, nil)
	})

	r.Get("/controllers/{filename}", func(w http.ResponseWriter, r *http.Request) {
		fn := chi.URLParam(r, "filename")

		content, err := os.ReadFile(fmt.Sprintf("web/controllers/%s", fn))

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	r.Get("/servers", func(w http.ResponseWriter, r *http.Request) {
		tmp, err := template.ParseFiles("web/template/servers/index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmp.Execute(w, nil)
	})
}

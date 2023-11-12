package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/krasilnikovm/logman/internal/application"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"application": "logman",
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(data)
	})

	configuration := application.ServerConfiguration{}
	logger := application.NewDefaultJsonLogger()

	if err := cleanenv.ReadEnv(&configuration); err != nil {
		logger.Error("can not read envs", slog.String("error", err.Error()))
	}

	s := application.NewServer(logger, configuration, r)

	if err := s.Run(); err != nil {
		logger.Error("Server is not laucnhed")
	}
}

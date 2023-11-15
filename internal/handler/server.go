package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/krasilnikovm/logman/internal/service"
)

type CreateServerRequest struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	LogLocationPath   string `json:"logLocationPath"`
	LogLocationFormat string `json:"logLocationFormat"`
}

type ServerHandlers struct {
	s service.ServerServiceContract
}

func NewServerHandlers(s service.ServerServiceContract) *ServerHandlers {
	return &ServerHandlers{
		s: s,
	}
}

func (l *ServerHandlers) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := l.s.FetchById(r.Context(), id)

	if err != nil {
		slog.Error("Unexpected error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (l *ServerHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody service.ServerData

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := l.s.Create(r.Context(), requestBody)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

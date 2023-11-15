package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/krasilnikovm/logman/internal/service"
)

type LogInfoHandlers struct {
	s service.LogInfoServiceContract
}

func NewLogInfoHandlers(s service.LogInfoServiceContract) *LogInfoHandlers {
	return &LogInfoHandlers{
		s: s,
	}
}

func (l *LogInfoHandlers) FetchLogInfoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := l.s.GetById(r.Context(), id)

	if err != nil {
		slog.Error("Unexpected error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

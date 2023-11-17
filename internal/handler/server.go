package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/krasilnikovm/logman/internal/service"
)

type ServerServiceContract interface {
	Create(ctx context.Context, data service.ServerData) (*service.ServerResponse, error)
	FetchById(ctx context.Context, id int) (*service.ServerResponse, error)
	DeleteById(ctx context.Context, id int) error
	GetList(ctx context.Context, limit, page int) ([]service.ServerResponse, error)
	Update(ctx context.Context, id int, data service.ServerData) (*service.ServerResponse, error)
}

type CreateServerRequest struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	LogLocationPath   string `json:"logLocationPath"`
	LogLocationFormat string `json:"logLocationFormat"`
}

type ServerHandlers struct {
	serverService ServerServiceContract
}

func NewServerHandlers(s ServerServiceContract) *ServerHandlers {
	return &ServerHandlers{
		serverService: s,
	}
}

func (s *ServerHandlers) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.serverService.FetchById(r.Context(), id)

	if err != nil {
		slog.Error("Unexpected error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeOkJson(w, response)
}

func (s *ServerHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody service.ServerData

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.serverService.Create(r.Context(), requestBody)

	if errors.As(err, &service.ErrValidation{}) {
		writeValidationJson(w, err.(service.ErrValidation))
		return
	}

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeOkJson(w, response)
}

func (s *ServerHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.serverService.DeleteById(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeWithEmptyBody(w)
}

func (s *ServerHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	response, err := s.serverService.GetList(r.Context(), limit, page)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeOkJson(w, response)
}

func (s *ServerHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody service.ServerData

	err = json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.serverService.Update(r.Context(), id, requestBody)

	if errors.Is(err, service.ErrValidation{}) {
		writeValidationJson(w, err.(service.ErrValidation))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeOkJson(w, response)
}

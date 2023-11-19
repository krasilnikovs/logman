package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/krasilnikovm/logman/internal/service"
)

type CredentialServiceContract interface {
	Create(ctx context.Context, data service.CredentialData) (service.CredentialResponse, error)
	Update(ctx context.Context, id int, data service.CredentialData) (*service.CredentialResponse, error)
	DeleteById(ctx context.Context, id int) error
	GetList(ctx context.Context, page, limit int) ([]service.CredentialResponse, error)
	GetById(ctx context.Context, id int) (*service.CredentialResponse, error)
}

type CredentialHandlers struct {
	credentialService CredentialServiceContract
}

func NewCredentialHandlers(s CredentialServiceContract) *CredentialHandlers {
	return &CredentialHandlers{
		credentialService: s,
	}
}

func (s *CredentialHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var request service.CredentialData

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.credentialService.Create(r.Context(), service.CredentialData{
		Path: request.Path,
		Name: request.Name,
	})

	if errors.As(err, &service.ErrValidation{}) {
		writeValidationJson(w, err.(service.ErrValidation))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeOkJson(w, response)
}

func (s *CredentialHandlers) FetchById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.credentialService.GetById(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeOkJson(w, response)
}

func (s *CredentialHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.credentialService.DeleteById(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *CredentialHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	response, err := s.credentialService.GetList(r.Context(), page, limit)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeOkJson(w, response)
}

func (s *CredentialHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request service.CredentialData

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.credentialService.Update(
		r.Context(),
		id,
		service.CredentialData{
			Path: request.Path,
			Name: request.Name,
		},
	)

	if errors.As(err, &service.ErrValidation{}) {
		writeValidationJson(w, err.(service.ErrValidation))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeOkJson(w, response)
}

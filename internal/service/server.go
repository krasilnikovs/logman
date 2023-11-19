package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/krasilnikovm/logman/internal/entity"
)

type ErrValidation struct {
	Errors []string `json:"errors"`
}

func (e ErrValidation) Error() string {
	return strings.Join(e.Errors, ", ")
}

type ServerStorager interface {
	Create(ctx context.Context, server *entity.Server) error
	GetById(ctx context.Context, id int) (*entity.Server, error)
	DeleteById(ctx context.Context, id int) error
	GetList(ctx context.Context, limit, page int) ([]entity.Server, error)
	Update(ctx context.Context, server *entity.Server, id int) error
}

type Validator interface {
	Struct(s interface{}) error
}

type ServerData struct {
	Name         string           `json:"name"`
	Host         string           `json:"host"`
	CredentialId int              `json:"credentialId"`
	LogLocation  LogLocationModel `json:"logLocation"`
}

type ServerResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Host          string `json:"host"`
	LogFolderPath string `json:"log_folder_path"`
	LogFormat     string `json:"log_format"`
	CredentialId  int    `json:"credentialId"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type LogLocationModel struct {
	Path   string `json:"path"`
	Format string `json:"format"`
}

type ServerService struct {
	storage           ServerStorager
	credentialStorage CredentialStorager
	l                 Logger
	v                 Validator
}

func NewServerService(storage ServerStorager, credentialStorage CredentialStorager, l Logger, v Validator) *ServerService {
	return &ServerService{
		storage:           storage,
		credentialStorage: credentialStorage,
		v:                 v,
	}
}

func (l *ServerService) FetchById(ctx context.Context, id int) (*ServerResponse, error) {
	server, err := l.storage.GetById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error during Server search by id: %w", err)
	}

	if server == nil {
		return nil, nil
	}

	return createServerResponseFromServerEntity(*server), nil
}

func (s *ServerService) Create(ctx context.Context, data ServerData) (*ServerResponse, error) {

	credential, err := s.credentialStorage.GetById(ctx, data.CredentialId)

	if err != nil {
		return nil, fmt.Errorf("error during Credential search by id: %w", err)
	}

	if credential == nil {
		return nil, ErrValidation{Errors: []string{fmt.Sprintf("credential with id %d not found", data.CredentialId)}}
	}

	now := time.Now()

	server := &entity.Server{
		Name:       data.Name,
		Host:       data.Host,
		Credential: *credential,
		LogPath:    entity.LogFolderPath(data.LogLocation.Path),
		LogFormat:  entity.LogFormat(data.LogLocation.Format),
		CreatedAt:  now.Format(time.RFC3339),
		UpdatedAt:  now.Format(time.RFC3339),
	}

	if err := s.v.Struct(server); err != nil {
		return nil, buildValidationError(err)
	}

	if err := s.storage.Create(ctx, server); err != nil {
		return nil, fmt.Errorf("error during creating server: %w", err)
	}

	return createServerResponseFromServerEntity(*server), nil
}

func (s *ServerService) DeleteById(ctx context.Context, id int) error {
	if err := s.storage.DeleteById(ctx, id); err != nil {
		return fmt.Errorf("delete by id failed: %w", err)
	}

	return nil
}

func (s *ServerService) GetList(ctx context.Context, limit, page int) ([]ServerResponse, error) {
	servers, err := s.storage.GetList(ctx, limit, page)

	if err != nil {
		slog.Error(err.Error())
		return []ServerResponse{}, fmt.Errorf("error during reading data from storage: %w", err)
	}

	responses := make([]ServerResponse, len(servers))

	for i, server := range servers {
		responses[i] = *createServerResponseFromServerEntity(server)
	}

	return responses, nil
}

func (s *ServerService) Update(ctx context.Context, id int, data ServerData) (*ServerResponse, error) {
	server, err := s.storage.GetById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error during Server search by id: %w", err)
	}

	if server == nil {
		return nil, nil
	}

	now := time.Now()

	server.Name = data.Name
	server.Host = data.Host
	server.LogPath = entity.LogFolderPath(data.LogLocation.Path)
	server.LogFormat = entity.LogFormat(data.LogLocation.Format)
	server.UpdatedAt = now.Format(time.RFC3339)

	if err := s.v.Struct(server); err != nil {
		return nil, buildValidationError(err)
	}

	if err := s.storage.Update(ctx, server, id); err != nil {
		return nil, fmt.Errorf("error during updating server: %w", err)
	}

	return createServerResponseFromServerEntity(*server), nil
}

func createServerResponseFromServerEntity(s entity.Server) *ServerResponse {
	return &ServerResponse{
		Id:            s.Id,
		Name:          s.Name,
		Host:          s.Host,
		CredentialId:  s.Credential.Id,
		LogFolderPath: string(s.LogPath),
		LogFormat:     string(s.LogFormat),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func buildValidationError(err error) ErrValidation {
	var errs []string

	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, fmt.Sprintf(`invalid '%s' field, please check the '%s' is an %s`, e.Field(), e.Field(), e.Tag()))
	}

	return ErrValidation{Errors: errs}
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/krasilnikovm/logman/internal/entity"
)

type CredentialStorager interface {
	Create(ctx context.Context, credential *entity.Credential) error
	GetById(ctx context.Context, id int) (*entity.Credential, error)
	GetList(ctx context.Context, page, limit int) ([]*entity.Credential, error)
	DeleteById(ctx context.Context, id int) error
	Update(ctx context.Context, credential *entity.Credential) error
}

type CredentialData struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CredentialResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CredentialService struct {
	storage CredentialStorager

	validator Validator
}

func NewCredentialService(storage CredentialStorager, validator Validator) *CredentialService {
	return &CredentialService{
		storage:   storage,
		validator: validator,
	}
}

func (c *CredentialService) Create(ctx context.Context, data CredentialData) (CredentialResponse, error) {

	now := time.Now()

	credential := &entity.Credential{
		Path:      entity.KeyPath(data.Path),
		Name:      data.Name,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	if err := c.validator.Struct(credential); err != nil {
		return CredentialResponse{}, buildValidationError(err)
	}

	err := c.storage.Create(ctx, credential)

	if err != nil {
		return CredentialResponse{}, fmt.Errorf("error during Credential creation: %w", err)
	}

	response := CredentialResponse{
		Id:        credential.Id,
		Name:      credential.Name,
		Path:      string(credential.Path),
		CreatedAt: credential.CreatedAt,
		UpdatedAt: credential.UpdatedAt,
	}

	return response, nil
}

func (c *CredentialService) Update(ctx context.Context, id int, data CredentialData) (*CredentialResponse, error) {
	now := time.Now()

	credential := &entity.Credential{
		Id:        id,
		Name:      data.Name,
		Path:      entity.KeyPath(data.Path),
		UpdatedAt: now.Format(time.RFC3339),
		CreatedAt: now.Format(time.RFC3339),
	}

	if err := c.validator.Struct(credential); err != nil {
		return nil, buildValidationError(err)
	}

	err := c.storage.Update(ctx, credential)

	if err != nil {
		return nil, fmt.Errorf("error during Credential update: %w", err)
	}

	return c.GetById(ctx, id)
}

func (c *CredentialService) DeleteById(ctx context.Context, id int) error {
	if err := c.storage.DeleteById(ctx, id); err != nil {
		return fmt.Errorf("error during Credential deletion: %w", err)
	}

	return nil
}

func (c *CredentialService) GetList(ctx context.Context, page, limit int) ([]CredentialResponse, error) {
	credentials, err := c.storage.GetList(ctx, page, limit)

	if err != nil {
		return nil, fmt.Errorf("error during Credential list getting: %w", err)
	}

	responses := make([]CredentialResponse, len(credentials))

	for i, credential := range credentials {
		responses[i] = CredentialResponse{
			Id:        credential.Id,
			Name:      credential.Name,
			Path:      string(credential.Path),
			CreatedAt: credential.CreatedAt,
			UpdatedAt: credential.UpdatedAt,
		}
	}

	return responses, nil
}

func (c *CredentialService) GetById(ctx context.Context, id int) (*CredentialResponse, error) {
	credential, err := c.storage.GetById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error during Credential search by id: %w", err)
	}

	if credential == nil {
		return nil, nil
	}

	response := &CredentialResponse{
		Id:        credential.Id,
		Name:      credential.Name,
		Path:      string(credential.Path),
		CreatedAt: credential.CreatedAt,
		UpdatedAt: credential.UpdatedAt,
	}

	return response, nil
}

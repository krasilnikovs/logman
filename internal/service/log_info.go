package service

import (
	"context"
	"fmt"
	"time"

	"github.com/krasilnikovm/logman/internal/entity"
)

type LogInfoStorager interface {
	Create(ctx context.Context, logInfo entity.LogInfo) error
	GetById(ctx context.Context, id int) (*entity.LogInfo, error)
}

type LogInfoServiceContract interface {
	GetById(ctx context.Context, id int) (*LogInfoResponse, error)
}

type LogInfoResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogInfoService struct {
	storage LogInfoStorager
}

func NewLogInfoService(storage LogInfoStorager) *LogInfoService {
	return &LogInfoService{
		storage: storage,
	}
}

func (l *LogInfoService) GetById(ctx context.Context, id int) (*LogInfoResponse, error) {
	logInfo, err := l.storage.GetById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error during LogInfo search by id: %w", err)
	}

	if logInfo == nil {
		return nil, nil
	}

	response := &LogInfoResponse{
		Id:        logInfo.Id,
		Name:      logInfo.Name,
		Location:  logInfo.Location,
		Format:    logInfo.Format,
		CreatedAt: logInfo.CreatedAt,
		UpdatedAt: logInfo.UpdatedAt,
	}

	return response, nil
}

package services

import (
	"context"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

type LogService struct {
	Repository repositories.Repository[models.LogRequest, models.LogData]
}

func (s *LogService) DeleteLog(ctx context.Context, jobId string, logId string) error {
	return s.Repository.Delete(ctx, jobId, logId)
}

func (s *LogService) GetAllLogs(ctx context.Context) ([]models.LogData, error) {
	return s.Repository.GetAll(ctx)
}

func (s *LogService) GetLog(ctx context.Context, jobId string, logId string) (*models.LogData, error) {
	return s.Repository.Get(ctx, jobId, logId)
}

func (s *LogService) PutLog(ctx context.Context, logRequest *models.LogRequest) error {
	return s.Repository.Put(ctx, logRequest)
}

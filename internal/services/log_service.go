package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

func NewLogService(repository *repositories.Repository, modelIdentifiers *models.ModelIdentifiers, routeKey string) (Service, error) {
	if routeKey == "DELETE /{PartitionType}/{PartitionId}/{SortType}" || routeKey == "PUT /{PartitionType}/{PartitionId}/{SortType}" {
		return nil, errors.New("invalid service action")
	}

	if strings.Contains(routeKey, "/{PartitionType}") && modelIdentifiers.PartitionType != models.ModelTypeJob {
		return nil, errors.New("invalid partition type")
	}

	if strings.Contains(routeKey, "/{PartitionId}") && modelIdentifiers.PartitionId == "" {
		return nil, errors.New("invalid partition ID")
	}

	if strings.Contains(routeKey, "/{SortType}") && modelIdentifiers.SortType != models.ModelTypeLog {
		return nil, errors.New("invalid sort type")
	}

	if strings.Contains(routeKey, "/{SortId}") && modelIdentifiers.SortId == "" {
		return nil, errors.New("invalid sort ID")
	}

	return &LogService{Repository: repository, ModelIdentifiers: modelIdentifiers}, nil
}

type LogService struct {
	Repository       *repositories.Repository
	ModelIdentifiers *models.ModelIdentifiers
}

func (s *LogService) DeleteByPartitionIdAndSortId(ctx context.Context) error {
	return s.Repository.DeleteByPartitionIdAndSortId(ctx, s.ModelIdentifiers)
}

func (s *LogService) GetByPartitionId(ctx context.Context) (interface{}, error) {
	return s.Repository.GetByPartitionId(ctx, s.ModelIdentifiers, new(models.LogItem))
}

func (s *LogService) GetByPartitionIdAndSortId(ctx context.Context) (models.ModelData, error) {
	return s.Repository.GetByPartitionIdAndSortId(ctx, s.ModelIdentifiers, new(models.LogItem))
}

func (s *LogService) GetBySortType(ctx context.Context) ([]models.ModelData, error) {
	return s.Repository.GetBySortType(ctx, s.ModelIdentifiers, new(models.LogItem))
}

func (s *LogService) PutByPartitionIdAndSortId(ctx context.Context, requestBody string) error {
	modelPayload := new(models.LogPayload)
	err := json.Unmarshal([]byte(requestBody), modelPayload)
	if err != nil {
		return err
	}
	return s.Repository.PutByPartitionIdAndSortId(ctx, s.ModelIdentifiers, modelPayload)
}

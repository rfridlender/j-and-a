package services

import (
	"context"
	"errors"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

func New(repository *repositories.Repository, modelIdentifiers *models.ModelIdentifiers) (Service, error) {
	switch modelIdentifiers.SortType {
	case models.ModelTypeLog:
		return NewLogService(repository, modelIdentifiers)
	case models.ModelTypePersonMetadata:
		return NewPersonMetadataService(repository, modelIdentifiers)
	default:
		return nil, errors.New("unsupported service")
	}
}

type Service interface {
	DeleteByPartitionIdAndSortId(ctx context.Context) error
	GetByPartitionId(ctx context.Context) (interface{}, error)
	GetByPartitionIdAndSortId(ctx context.Context) (models.ModelData, error)
	GetBySortType(ctx context.Context) ([]models.ModelData, error)
	PutByPartitionIdAndSortId(ctx context.Context, requestBody string) error
}

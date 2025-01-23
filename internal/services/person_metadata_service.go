package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

func NewPersonMetadataService(repository *repositories.Repository, modelIdentifiers *models.ModelIdentifiers, routeKey string) (Service, error) {
	if strings.Contains(routeKey, "/{PartitionType}") && modelIdentifiers.PartitionType != models.ModelTypePerson {
		return nil, errors.New("invalid partition type")
	}

	if strings.Contains(routeKey, "/{PartitionId}") && modelIdentifiers.PartitionId == "" {
		return nil, errors.New("invalid partition ID")
	}

	if strings.Contains(routeKey, "/{SortType}") && modelIdentifiers.SortType != models.ModelTypePersonMetadata {
		return nil, errors.New("invalid sort type")
	}

	if strings.Contains(routeKey, "/{SortId}") {
		return nil, errors.New("invalid service action")
	}

	return &PersonMetadataService{Repository: repository, ModelIdentifiers: modelIdentifiers}, nil
}

type PersonMetadataService struct {
	Repository       *repositories.Repository
	ModelIdentifiers *models.ModelIdentifiers
}

func (s *PersonMetadataService) DeleteByPartitionIdAndSortId(ctx context.Context) error {
	s.ModelIdentifiers.SortId = s.ModelIdentifiers.PartitionId
	return s.Repository.DeleteByPartitionIdAndSortId(ctx, s.ModelIdentifiers)
}

func (s *PersonMetadataService) GetByPartitionId(ctx context.Context) (interface{}, error) {
	s.ModelIdentifiers.SortId = s.ModelIdentifiers.PartitionId
	return s.Repository.GetByPartitionIdAndSortId(ctx, s.ModelIdentifiers, new(models.PersonMetadataItem))
}

func (s *PersonMetadataService) GetByPartitionIdAndSortId(ctx context.Context) (models.ModelData, error) {
	return nil, errors.New("invalid service action")
}

func (s *PersonMetadataService) GetBySortType(ctx context.Context) ([]models.ModelData, error) {
	return s.Repository.GetBySortType(ctx, s.ModelIdentifiers, new(models.PersonMetadataItem))
}

func (s *PersonMetadataService) PutByPartitionIdAndSortId(ctx context.Context, requestBody string) error {
	modelPayload := new(models.PersonMetadataPayload)
	err := json.Unmarshal([]byte(requestBody), modelPayload)
	if err != nil {
		return err
	}
	s.ModelIdentifiers.SortId = s.ModelIdentifiers.PartitionId
	return s.Repository.PutByPartitionIdAndSortId(ctx, s.ModelIdentifiers, modelPayload)
}

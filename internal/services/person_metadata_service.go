package services

import (
	"context"
	"encoding/json"
	"errors"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

func NewPersonMetadataService(repository *repositories.Repository, modelIdentifiers *models.ModelIdentifiers) (Service, error) {
	if modelIdentifiers.SortId != "" {
		return nil, errors.New("sort ID must not be specified")
	}

	if modelIdentifiers.PartitionType != models.ModelTypePerson {
		return nil, errors.New("invalid partition type")
	}

	if modelIdentifiers.SortType != models.ModelTypePersonMetadata {
		return nil, errors.New("invalid sort type")
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
	personMetadataPayload := new(models.PersonMetadataPayload)
	err := json.Unmarshal([]byte(requestBody), personMetadataPayload)
	if err != nil {
		return err
	}
	s.ModelIdentifiers.SortId = s.ModelIdentifiers.PartitionId
	return s.Repository.PutByPartitionIdAndSortId(ctx, s.ModelIdentifiers, personMetadataPayload)
}

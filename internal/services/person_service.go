package services

import (
	"context"

	"j-and-a/internal/models"
	"j-and-a/internal/repositories"
)

type PersonService struct {
	Repository repositories.Repository[models.PersonRequest, models.PersonData]
}

func (s *PersonService) DeletePerson(ctx context.Context, personId string) error {
	return s.Repository.Delete(ctx, personId, personId)
}

func (s *PersonService) GetAllPersons(ctx context.Context) ([]models.PersonData, error) {
	return s.Repository.GetAll(ctx)
}

func (s *PersonService) GetPerson(ctx context.Context, personId string) (*models.PersonData, error) {
	return s.Repository.Get(ctx, personId, personId)
}

func (s *PersonService) PutPerson(ctx context.Context, personRequest *models.PersonRequest) error {
	return s.Repository.Put(ctx, personRequest)
}

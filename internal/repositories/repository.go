package repositories

import "context"

type Repository[Request any, Data any] interface {
	Delete(ctx context.Context, partitionId string, sortId string) error
	GetAll(ctx context.Context) ([]Data, error)
	Get(ctx context.Context, partitionId string, sortId string) (*Data, error)
	Put(ctx context.Context, request *Request) error
}

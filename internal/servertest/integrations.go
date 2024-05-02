package servertest

import "context"

type Repository struct {
	Storage Storage
}

type StorageEntity struct {
	ID          string
	Name        string
	Description string
	SomeValue   int
}

type Storage interface {
	CreateEntity(ctx context.Context, entity StorageEntity) error
	ListEntities(ctx context.Context) ([]StorageEntity, error)
}

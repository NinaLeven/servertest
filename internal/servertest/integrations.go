package servertest

import "context"

type Repository struct {
	Storage Storage
}

type StorageEntity struct {
	RiderID string
	Lat     float32
	Long    float32
}

type StorageListLocationEntitiesOptions struct {
	RiderID string
	Limit   *int
}

type Storage interface {
	CreateLocationEntity(ctx context.Context, entity StorageEntity) error
	ListLocationEntities(ctx context.Context, opts StorageListLocationEntitiesOptions) ([]StorageEntity, error)
}

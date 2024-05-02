package memstorage

import (
	"context"
	"sync"

	"servertest/internal/servertest"
)

type MemoryStorage struct {
	rw     sync.RWMutex
	memory map[string]servertest.StorageEntity
}

var _ servertest.Storage = (*MemoryStorage)(nil)

func NewMemStorage() *MemoryStorage {
	return &MemoryStorage{
		memory: make(map[string]servertest.StorageEntity),
	}
}

func (c *MemoryStorage) CreateEntity(ctx context.Context, entity servertest.StorageEntity) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	c.memory[entity.ID] = entity

	return nil
}

func (c *MemoryStorage) ListEntities(ctx context.Context) ([]servertest.StorageEntity, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	res := make([]servertest.StorageEntity, 0, len(c.memory))
	for _, entity := range c.memory {
		res = append(res, entity)
	}

	return res, nil
}

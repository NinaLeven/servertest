package memstorage

import (
	"context"
	"log/slog"
	"sync"

	"servertest/internal/servertest"
)

type MemoryStorage struct {
	rw     sync.RWMutex
	memory map[string][]servertest.StorageEntity
}

var _ servertest.Storage = (*MemoryStorage)(nil)

func NewMemStorage() *MemoryStorage {
	return &MemoryStorage{
		memory: make(map[string][]servertest.StorageEntity),
	}
}

func (c *MemoryStorage) CreateLocationEntity(ctx context.Context, entity servertest.StorageEntity) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	c.memory[entity.RiderID] = append(c.memory[entity.RiderID], entity)

	slog.Info("map", "map", c.memory)

	return nil
}

func (c *MemoryStorage) ListLocationEntities(ctx context.Context, opts servertest.StorageListLocationEntitiesOptions) ([]servertest.StorageEntity, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	history, ok := c.memory[opts.RiderID]
	if !ok {
		return nil, &servertest.ErrNotFound{}
	}

	res := history
	if opts.Limit != nil {
		res = res[len(res)-min(*opts.Limit, len(res)):]
	}

	return res, nil
}

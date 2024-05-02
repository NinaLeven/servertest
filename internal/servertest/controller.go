package servertest

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type CreateEntityOptions struct {
	Name        string
	Description string
	SomeValue   int
}

type Controller interface {
	CreateEntity(ctx context.Context, repo Repository, opts CreateEntityOptions) (string, error)
	AggregateSomeValue(ctx context.Context, repo Repository) (int, error)
}

var _ Controller = (*controller)(nil)

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

func (c *controller) CreateEntity(ctx context.Context, repo Repository, opts CreateEntityOptions) (string, error) {
	id := uuid.NewString()

	slog.Info("kek")

	err := repo.Storage.CreateEntity(ctx, StorageEntity{
		ID:          id,
		Name:        opts.Name,
		Description: opts.Description,
		SomeValue:   opts.SomeValue,
	})
	if err != nil {
		return "", fmt.Errorf("unable to store entity: %w", err)
	}

	return id, nil
}

func (c *controller) AggregateSomeValue(ctx context.Context, repo Repository) (int, error) {
	entities, err := repo.Storage.ListEntities(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to list entities: %w", err)
	}

	var value int
	for _, entity := range entities {
		value += entity.SomeValue
	}

	return value, nil
}

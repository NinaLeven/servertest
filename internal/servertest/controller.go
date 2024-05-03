package servertest

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

type LocationEntry struct {
	RiderID string
	Lat     float32
	Long    float32
}

type GetRiderLocationHistoryOptions struct {
	RiderID string
	Limit   *int
}

type Controller interface {
	AddRiderLocation(ctx context.Context, repo Repository, opts LocationEntry) error
	GetRiderLocationHistory(ctx context.Context, repo Repository, opts GetRiderLocationHistoryOptions) ([]LocationEntry, error)
}

type ErrNotFound struct {
	Err error
}

func (e *ErrNotFound) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("not found")
	}
	return fmt.Sprintf("not found: %s", e.Err.Error())
}

func (e *ErrNotFound) Is(target error) bool {
	_, ok := target.(*ErrNotFound)
	if ok {
		return true
	}
	return errors.Is(e.Err, target)
}

var _ Controller = (*controller)(nil)

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

func (c *controller) AddRiderLocation(ctx context.Context, repo Repository, opts LocationEntry) error {
	return repo.Storage.CreateLocationEntity(ctx, StorageEntity{
		RiderID: opts.RiderID,
		Long:    opts.Long,
		Lat:     opts.Lat,
	})
}

const defaultLimit = 1000

func (c *controller) GetRiderLocationHistory(ctx context.Context, repo Repository, opts GetRiderLocationHistoryOptions) ([]LocationEntry, error) {
	limit := defaultLimit
	if opts.Limit != nil {
		limit = min(defaultLimit, *opts.Limit)
	}

	history, err := repo.Storage.ListLocationEntities(ctx, StorageListLocationEntitiesOptions{
		RiderID: opts.RiderID,
		Limit:   &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list history: %w", err)
	}

	res := make([]LocationEntry, 0, len(history))
	for _, e := range history {
		res = append(res, LocationEntry{
			Long:    e.Long,
			Lat:     e.Lat,
			RiderID: e.RiderID,
		})
	}

	slices.Reverse(res)

	return res, nil
}

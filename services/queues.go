package services

import (
	"context"

	"github.com/gork-io/gork/models"
	"github.com/pkg/errors"
)

// Queues service implements operations that are related to the queues management.
type Queues struct {
	queuesRepo models.QueuesRepository // queues repository
}

// List returns a subset of the queries, based on collection params given.
func (svc *Queues) List(
	ctx context.Context,
	params *models.CollectionParams,
) (records []*models.Queue, info *models.CollectionInfo, err error) {

	records, info, err = svc.queuesRepo.Find(ctx, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Find failed")
	}

	return
}

// Create creates a new queue instance and saves it to the repository.
func (svc *Queues) Create(ctx context.Context, name string, settings ...models.QueueSetting) (record *models.Queue, err error) {

	record = models.NewQueue(name, settings...)

	err = svc.queuesRepo.Save(ctx, record)
	if err != nil {
		return nil, errors.Wrap(err, "Save failed")
	}

	return
}

// Read returns query by its ID.
func (svc *Queues) Read(ctx context.Context, id string) (record *models.Queue, err error) {

	record, err = svc.queuesRepo.GetById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "GetById failed")
	}

	return
}

// Create creates a new queue instance and saves it to the repository.
func (svc *Queues) Update(ctx context.Context, id string, data *models.Queue) (record *models.Queue, err error) {

	record = models.NewQueue(name, settings...)

	err = svc.queuesRepo.Insert(ctx, record)
	if err != nil {
		return nil, errors.Wrap(err, "Save failed")
	}

	return
}

// Delete removes queue with given ID from the repository.
func (svc *Queues) Delete(ctx context.Context, id string) (err error) {
	return errors.Wrap(svc.queuesRepo.Delete(ctx, id), "Delete failed")
}

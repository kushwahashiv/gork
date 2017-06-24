package resources

import (
	"context"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gork-io/gork/models"
	"github.com/pkg/errors"
)

// NewQueues creates a new instance of Queues.
func NewQueues(queuesRepo models.QueuesRepository) (res *Queues) {
	return &Queues{
		queuesRepo: queuesRepo,
	}
}

// Queues resource service implements operations that are related to the queues management.
type Queues struct {
	queuesRepo models.QueuesRepository // queues repository
}

// List returns a subset of the queries, based on collection params given.
func (res *Queues) List(
	ctx context.Context,
	params *models.CollectionParams,
) (records []*models.Queue, info *models.CollectionInfo, err error) {

	// Retrieve collection from the repo
	records, info, err = res.queuesRepo.Find(ctx, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "repository Find failed")
	}

	return
}

// Create creates a new queue instance and saves it to the repository.
func (res *Queues) Create(
	ctx context.Context,
	name string,
	settings map[models.QueueSetting]string,
) (record *models.Queue, err error) {

	// Validate input
	vErr := validation.Errors{
		"name": validateQueueName(name),
	}
	for key, value := range settings {
		vErr["settings["+string(key)+"]"] = validateQueueSetting(key, value)
	}
	vErr.Filter()
	if vErr != nil {
		return nil, errors.Wrap(err, "validation error")
	}
	existing, err := res.queuesRepo.GetByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "repository GetByName failed")
	}
	if existing != nil {
		return nil, errors.New("queue with such name already exists")
	}

	// Save record to the repo
	err = res.queuesRepo.Save(ctx, models.NewQueue(name, settings))
	if err != nil {
		return nil, errors.Wrap(err, "repository Save failed")
	}

	return
}

// Read returns query by its ID.
func (res *Queues) Read(ctx context.Context, id string) (record *models.Queue, err error) {

	// Retrieve record from the repo
	record, err = res.queuesRepo.GetById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository GetById failed")
	}

	return
}

// Delete removes queue with given ID from the repository.
func (res *Queues) Delete(ctx context.Context, id string) (err error) {
	return errors.Wrap(res.queuesRepo.Delete(ctx, id), "repository Delete failed")
}

// validateQueueName checks that queue name is valid.
func validateQueueName(name string) (err error) {
	return validation.Validate(name, validation.Length(1, 255))
}

// validateQueueSetting checks that option with given name an value is valid.
func validateQueueSetting(name models.QueueSetting, value string) (err error) {
	switch name {
	case models.QueueSettingRateLimitEnabled:
		err = validation.Validate(value, validation.In("0", "1"))
	case models.QueueSettingRateLimitTokens:
		err = validation.Validate(value, is.Int)
	case models.QueueSettingRateLimitDuration:
		err = validation.Validate(value, is.Int, validation.Max(86400))
	default:
		err = errors.New("Unknown setting")
	}
	return
}

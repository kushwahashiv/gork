package redis

import (
	"context"

	"strings"

	"github.com/go-redis/redis"
	"github.com/gork-io/gork/models"
	"github.com/pkg/errors"
)

const (
	queuesKeyIndexById            string = "queues:index:id"
	queuesKeyData                 string = "queues"
	queuesSuffixSettingsRateLimit string = "settings:rate_limit"
)

// QueuesRepository implements a Redis-based queues repository.
//
// Redis schema:
//   - SET: `queues:index:id`.
//     An index containing IDs of all known queues.
//   - HASH: `queues:<queue ID>`.
//     Generic queue information.
//     Fields:
//       - `id`;
//       - `name`;
//       - `created_at`.
//   - HASH: `queues:<queue ID>:settings:rate_limit`.
//     Rate limit setting data.
//     Fields:
//       - `tokens`;
//       - `duration`.
type QueuesRepository struct {
	Prefix      string
	redisClient *redis.Client
}

func (repo *QueuesRepository) Save(ctx context.Context, record *models.Queue) (err error) {

	repo.redisClient.TxPipelined(func(pipe redis.Pipeliner) (err error) {
		pipe.HMSet()
		pipe.HMSet()
		return
	})

	return
}

// Delete removes from repository all information regarding the query with id given.
func (repo *QueuesRepository) Delete(ctx context.Context, id string) (err error) {

	_, err = repo.redisClient.TxPipelined(func(pipe redis.Pipeliner) (err error) {
		pipe.SRem()
		pipe.Del()
		return
	})

	return errors.Wrap(err, "transaction failed")
}

// GetById retrieves a single queue by it's ID.
func (repo *QueuesRepository) GetById(ctx context.Context, id string) (record *models.Queue, err error) {

	var (
		queueCmd                  *redis.StringStringMapCmd
		queueSettingsRateLimitCmd *redis.StringStringMapCmd
	)
	_, err = repo.redisClient.Pipelined(func(pipe redis.Pipeliner) (err error) {
		queueCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, id))
		queueSettingsRateLimitCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, id, queuesSuffixSettingsRateLimit))
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "pipeline failed")
	}

	return
}

func (repo *QueuesRepository) MGetById(ctx context.Context, ids ...string) (records []*models.Queue, err error) {

	repo.redisClient.Pipelined(func(pipe redis.Pipeliner) (err error) {
		pipe.HGetAll()
		pipe.HGetAll()
		return
	})

	return
}

func (repo *QueuesRepository) Find(
	ctx context.Context,
	params *models.CollectionParams,
) (records []*models.Queue, info *models.CollectionInfo, err error) {
	return
}

// buildKey is a helper function that builds a Redis key from key parts given.
func (repo *QueuesRepository) buildKey(parts ...string) (key string) {
	return strings.Join(append([]string{repo.Prefix}, parts...), ":")
}

func queueMarshal() (data) {

}

func queueUnmarshal() {

}

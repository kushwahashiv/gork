package redis

import (
	"context"

	"strings"

	"time"

	"strconv"

	"github.com/go-redis/redis"
	"github.com/gork-io/gork/models"
	"github.com/pkg/errors"
)

const (
	queuesKeyData        string = "queues"
	queuesSuffixSettings string = "settings"
	queuesKeyIndexById   string = "queues:index:id"
	queuesKeyIndexByName string = "queues:index:name"
)

// NewQueuesRepository creates a new instance of QueuesRepository.
func NewQueuesRepository(redisClient *redis.Client) (repo *QueuesRepository) {
	return &QueuesRepository{
		redisClient: redisClient,
	}
}

// QueuesRepository implements a Redis-based queues repository.
//
// Redis schema:
//   - HASH: `queues:<queue ID>`.
//     Generic queue information.
//     Fields:
//       - `id`;
//       - `name`;
//       - `created_at`.
//   - HASH: `queues:<queue ID>:settings`.
//     Queue settings data.
//     Field names are values of the corresponding QueueSetting constants.
//   - SORTED SET: `queues:index:id`.
//     An index containing IDs of all known queues and creation timestamp as a score.
//   - STRING: `queues:index:name:<queue name>`.
//     An index containing Name => ID pairs of all known queues.
type QueuesRepository struct {
	redisClient *redis.Client // redis client instance
}

// Save persists given queue instance to the repo.
func (repo *QueuesRepository) Save(ctx context.Context, record *models.Queue) (err error) {

	data, settingsData := queueMarshal(record)
	_, err = repo.redisClient.WithContext(ctx).TxPipelined(func(pipe redis.Pipeliner) (err error) {
		pipe.HMSet(repo.buildKey(queuesKeyData, record.Id), data)
		pipe.HMSet(repo.buildKey(queuesKeyData, record.Id, queuesSuffixSettings), settingsData)
		pipe.ZAdd(repo.buildKey(queuesKeyIndexById), redis.Z{
			Member: record.Id,
			Score:  float64(record.CreatedAt.Nanosecond()),
		})
		pipe.Set(repo.buildKey(queuesKeyIndexByName, record.Name), record.Id, 0)
		return
	})

	return errors.Wrap(err, "transaction failed")
}

// Delete removes queue with given ID from the repo.
func (repo *QueuesRepository) Delete(ctx context.Context, id string) (err error) {

	clientCtx := repo.redisClient.WithContext(ctx)

	// Retrieve queue name
	nameCmd := clientCtx.HGet(repo.buildKey(queuesKeyData, id), "name")
	err = errors.Wrap(nameCmd.Err(), "failed to retrieve queue name")
	if err != nil {
		return err
	}

	// Delete all queue data
	_, err = repo.redisClient.TxPipelined(func(pipe redis.Pipeliner) (err error) {
		pipe.ZRem(repo.buildKey(queuesKeyIndexById), id)
		pipe.Del(
			repo.buildKey(queuesKeyIndexByName, nameCmd.Val()),
			repo.buildKey(queuesKeyData, id),
			repo.buildKey(queuesKeyData, id, queuesSuffixSettings),
		)
		return
	})

	return errors.Wrap(err, "transaction failed")
}

// GetById retrieves queue with given ID from the repo.
func (repo *QueuesRepository) GetById(ctx context.Context, id string) (record *models.Queue, err error) {

	var (
		dataCmd         *redis.StringStringMapCmd
		settingsDataCmd *redis.StringStringMapCmd
	)
	_, err = repo.redisClient.WithContext(ctx).Pipelined(func(pipe redis.Pipeliner) (err error) {
		dataCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, id))
		settingsDataCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, id, queuesSuffixSettings))
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "pipeline failed")
	}

	record = queueUnmarshal(dataCmd.Val(), settingsDataCmd.Val())
	return
}

// GetByName retrieves queue with given name from the repo.
func (repo *QueuesRepository) GetByName(ctx context.Context, name string) (record *models.Queue, err error) {

	clientCtx := repo.redisClient.WithContext(ctx)

	// Find ID by name
	idCmd := clientCtx.Get(repo.buildKey(queuesKeyIndexByName, name))
	err = errors.Wrap(idCmd.Err(), "failed to retrieve queue id")
	if err != nil {
		return
	}

	var (
		dataCmd         *redis.StringStringMapCmd
		settingsDataCmd *redis.StringStringMapCmd
	)
	_, err = repo.redisClient.Pipelined(func(pipe redis.Pipeliner) (err error) {
		dataCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, idCmd.Val()))
		settingsDataCmd = pipe.HGetAll(repo.buildKey(queuesKeyData, idCmd.Val(), queuesSuffixSettings))
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "pipeline failed")
	}

	record = queueUnmarshal(dataCmd.Val(), settingsDataCmd.Val())
	return
}

// MGetById retrieves queues with given IDs from the repo.
func (repo *QueuesRepository) MGetById(ctx context.Context, ids []string) (records []*models.Queue, err error) {

	var (
		dataCmds         []*redis.StringStringMapCmd
		settingsDataCmds []*redis.StringStringMapCmd
	)
	_, err = repo.redisClient.WithContext(ctx).Pipelined(func(pipe redis.Pipeliner) (err error) {
		for _, id := range ids {
			dataCmds = append(dataCmds, pipe.HGetAll(repo.buildKey(queuesKeyData, id)))
			settingsDataCmds = append(settingsDataCmds, pipe.HGetAll(repo.buildKey(queuesKeyData, id, queuesSuffixSettings)))
		}
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "pipeline failed")
	}

	for i, dataCmd := range dataCmds {
		records = append(records, queueUnmarshal(dataCmd.Val(), settingsDataCmds[i].Val()))
	}

	return
}

// Find returns a subset of the queries, based on collection params given.
func (repo *QueuesRepository) Find(
	ctx context.Context,
	params *models.CollectionParams,
) (records []*models.Queue, info *models.CollectionInfo, err error) {

	// Parse cursor
	cursor, err := strconv.ParseUint(params.Cursor, 10, 64)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse cursor")
	}

	// Retrieve indexes
	var (
		idxCmd   *redis.ScanCmd
		countCmd *redis.IntCmd
	)
	_, err = repo.redisClient.WithContext(ctx).Pipelined(func(pipe redis.Pipeliner) (err error) {
		key := repo.buildKey(queuesKeyIndexById)
		idxCmd = pipe.ZScan(key, cursor, "*", int64(params.Limit))
		countCmd = pipe.ZCard(key)
		return
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to scan queues index")
	}

	// Retrieve records
	keys, newCursor := idxCmd.Val()
	records, err = repo.MGetById(ctx, keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to retrieve records")
	}
	info = &models.CollectionInfo{
		Total:  uint64(countCmd.Val()),
		Cursor: strconv.FormatUint(newCursor, 10),
	}

	return
}

// buildKey is a helper function that builds a Redis key from key parts given.
func (repo *QueuesRepository) buildKey(parts ...string) (key string) {
	return strings.Join(parts, ":")
}

// queueMarshal is a helper function that marshals record into format Redis understands.
func queueMarshal(record *models.Queue) (data, settingsData map[string]interface{}) {

	data = make(map[string]interface{})
	settingsData = make(map[string]interface{})

	data["id"] = record.Id
	data["name"] = record.Name
	data["created_at"] = record.CreatedAt.Format(time.RFC3339Nano)

	for key, value := range record.Settings {
		settingsData[string(key)] = value
	}

	return
}

// queueUnmarshal is a helper function that unmarshals record from Redis format.
func queueUnmarshal(data, settingsData map[string]string) (record *models.Queue) {

	createdAt, _ := time.Parse(time.RFC3339Nano, data["created_at"])

	record = &models.Queue{
		Id:        data["id"],
		Name:      data["name"],
		CreatedAt: createdAt,
	}
	for key, value := range settingsData {
		record.Settings[models.QueueSetting(key)] = value
	}

	return
}

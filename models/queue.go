package models

import (
	"time"

	"context"

	"github.com/rs/xid"
)

// QueuesRepository is an interface that all queues storage should implement.
type QueuesRepository interface {
	Save(ctx context.Context, record *Queue) (err error)
	Delete(ctx context.Context, id string) (err error)
	GetById(ctx context.Context, id string) (record *Queue, err error)
	MGetById(ctx context.Context, ids ...string) (records []*Queue, err error)
	Find(ctx context.Context, params *CollectionParams) (records []*Queue, info *CollectionInfo, err error)
}

// NewQueue creates a new instance of Queue.
func NewQueue(name string, settings ...QueueSetting) (queue *Queue) {
	queue = &Queue{
		Id:        xid.New().String(),
		Name:      name,
		Settings:  &QueueSettings{},
		CreatedAt: time.Now(),
	}
	for _, setting := range settings {
		setting(queue)
	}
	return
}

// Queue represents a list of tasks of the same type that are pending processing.
type Queue struct {
	Id        string            // unique ID
	Name      string            // unique name
	Settings  map[string]string // settings
	CreatedAt time.Time         // creation time
}

// QueueSettings represents settings of the queue.
type QueueSettings struct {
	RateLimit *QueueSettingRateLimit // rate limiting
}

// QueueSettingRateLimit represents queue settings that are related to rate limiting.
type QueueSettingRateLimit struct {
	Tokens   uint64        // tokens amount
	Duration time.Duration // time period
}

// QueueSetting is used by NewQueue for easier queues creation.
type QueueSetting func(queue *Queue)

// QueueWithRateLimit applies specified rate limit settings to the new queue.
func QueueWithRateLimit(tokens uint64, duration time.Duration) (setting QueueSetting) {
	return func(queue *Queue) {
		queue.Settings.RateLimit = &QueueSettingRateLimit{
			Tokens:   tokens,
			Duration: duration,
		}
	}
}

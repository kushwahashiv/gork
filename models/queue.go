package models

import (
	"time"

	"context"

	"github.com/rs/xid"
)

const (
	QueueSettingRateLimitEnabled  QueueSetting = "rate-limit.enabled"
	QueueSettingRateLimitTokens   QueueSetting = "rate-limit.tokens"
	QueueSettingRateLimitDuration QueueSetting = "rate-limit.duration"
)

var (
	defaultSettings = map[QueueSetting]string{
		QueueSettingRateLimitEnabled:  "0",
		QueueSettingRateLimitTokens:   "0",
		QueueSettingRateLimitDuration: "0",
	}
)

// QueuesRepository is an interface that all queues storage should implement.
type QueuesRepository interface {
	// Save persists given queue instance to the repo.
	Save(ctx context.Context, record *Queue) (err error)
	// Delete removes queue with given ID from the repo.
	Delete(ctx context.Context, id string) (err error)
	// GetById retrieves queue with given ID from the repo.
	GetById(ctx context.Context, id string) (record *Queue, err error)
	// GetByName retrieves queue with given name from the repo.
	GetByName(ctx context.Context, name string) (record *Queue, err error)
	// MGetById retrieves queues with given IDs from the repo.
	MGetById(ctx context.Context, ids []string) (records []*Queue, err error)
	// Find returns a subset of the queries, based on collection params given.
	Find(ctx context.Context, params *CollectionParams) (records []*Queue, info *CollectionInfo, err error)
}

// NewQueue creates a new instance of Queue.
func NewQueue(name string, settings map[QueueSetting]string) (queue *Queue) {
	return &Queue{
		Id:        xid.New().String(),
		Name:      name,
		Settings:  mergeSettings(defaultSettings, settings),
		CreatedAt: time.Now(),
	}
}

// Queue represents a list of tasks of the same type that are pending processing.
type Queue struct {
	Id        string                  // unique ID
	Name      string                  // unique name
	Settings  map[QueueSetting]string // settings
	CreatedAt time.Time               // creation time
}

// QueueSetting represents an identifier of the query setting.
type QueueSetting string

// mergeSettings is a tiny helper that merges default and custom queue settings.
func mergeSettings(defaultSettings, customSettings map[QueueSetting]string) (merged map[QueueSetting]string) {
	for key, value := range customSettings {
		defaultSettings[key] = value
	}
	return defaultSettings
}

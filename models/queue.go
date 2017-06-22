package models

import "time"

// Queue represents a list of tasks of the same type that are pending processing.
type Queue struct {
	Name      string
	RateLimit *QueueRateLimit
	CreatedAt time.Time
}

type QueueRateLimit struct {
	Amount   uint64
	Duration time.Duration
}

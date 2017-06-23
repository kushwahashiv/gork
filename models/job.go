package models

import "time"

// Job represents a single unit of work that is delivered to the worker.
type Job struct {
	Id          string // unique ID
	Input       []byte // task payload
	Progress    uint8
	Logs        []string
	CreatedAt   time.Time
	DeliveredAt time.Time
	FinishedAt  time.Time
}

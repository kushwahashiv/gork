package models

import "time"

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusExpired
)

// TaskStatus represents
type TaskStatus uint8

// Task represents a single unit of work.
type Task struct {
	Id         string
	Status     TaskStatus
	Priority   uint8
	Headers    map[string]string
	Input      []byte
	Output     []byte
	CreatedAt  time.Time
	ExpiresAt  time.Time
	FinishedAt time.Time
}

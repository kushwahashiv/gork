package models

import "time"

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusProcessing
	TaskStatusExpired
	TaskStatusFinished
)

// TaskStatus represents a current state of the task.
type TaskStatus uint8

// Task represents a single unit of work that should be processed by worker(s).
type Task struct {
	Id         string            // unique ID
	Status     TaskStatus        // processing status
	Priority   uint8             // priority level
	Headers    map[string]string // custom key->value pairs
	Input      []byte            // payload data
	CreatedAt  time.Time         // creation time
	ExpiresAt  time.Time         // expiration time
	FinishedAt time.Time         // processing finish time
}

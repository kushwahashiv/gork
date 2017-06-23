package future

import "time"

// Exchange represents an input that accepts tasks from the publishers and routes them to the appropriate queues.
type Exchange struct {
	Id        string    // unique ID
	Name      string    // unique name
	CreatedAt time.Time // creation time
}

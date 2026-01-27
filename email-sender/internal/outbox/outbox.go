package outbox

import (
	"time"
)

type Status string

type EventType string

var (
	StatusNew    Status = "NEW"
	StatusSent   Status = "SENT"
	StatusFailed Status = "FAILED"
)

type OutboxEmail struct {
	ID            int64
	Payload       []byte
	Status        Status
	DedupKey      string
	Attempts      int
	NextAttemptAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Type          string
}

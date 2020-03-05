package persistence

import (
	"time"
)

type EventRepository interface {
	Accumulate(event Event) (int, error)
}

type Event struct {
	Time      time.Time
	Initiator string
}

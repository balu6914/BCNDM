package persistence

import (
	"time"
)

type EventRepository interface {
	Accumulate(event Event, unit string) (int, error)
}

type Event struct {
	Time      time.Time
	Initiator string
}

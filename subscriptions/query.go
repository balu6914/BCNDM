package subscriptions

import (
	"time"
)

// Query struct wraps query parameters and provides
// "query builder" as a convenient way to generate DB query.
type Query struct {
	StreamOwner string
	UserID      string
	Page        uint64
	Limit       uint64
	StreamID    string
	StartTime   *time.Time
}

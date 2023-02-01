package subscriptions

import "time"

// StreamsService contains API for fetching stream data.
type StreamsService interface {
	// One returns specified stream by its id.
	One(string) (Stream, error)
}

// Stream represents state of one stream.
type Stream struct {
	ID         string
	Name       string
	Owner      string
	URL        string
	Price      uint64
	External   bool
	Project    string
	Dataset    string
	Table      string
	Fields     string
	Visibility string
	AccessType string
	MaxCalls   uint64
	MaxUnit    string
	EndDate    *time.Time
}

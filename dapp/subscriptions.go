package dapp

import "time"

type StreamData struct {
	Coordinates [2]float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
	Name        string     `bson:"name,omitempty" json:"name,omitempty"`
	Price       int        `bson:"price,omitempty" json:"price,omitempty"`
}
type Subscription struct {
	UserID     string     `bson:"user_id,omitempty" json:"user_id,omitempty"`
	StreamID   string     `bson:"id,omitempty" json:"id,omitempty"`
	StreamData StreamData `bson:"stream_data,omitempty" json:"stream_data,omitempty"`
	Hours      uint64     `bson:"hours,omitempty" json:"hours,omitempty"`
	StartDate  time.Time  `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate    time.Time  `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StreamURL  string     `bson:"stream_url,omitempty" json:"stream_url,omitempty"`
}

// Validate returns an error if user representation is invalid.
func (s *Subscription) Validate() error {

	if s.Hours == 0 || s.StreamID == "" {
		return ErrMalformedEntity
	}

	return nil
}

// SubscriptionsRepository specifies a subscription persistence API.
type SubscriptionsRepository interface {
	// Create persists the subscription.
	Create(Subscription) error

	// Read retrieves a list of subscription by userID
	Read(string) ([]Subscription, error)

	// Update performs an update of an existing Subscription
	Update(string, Subscription) error
}

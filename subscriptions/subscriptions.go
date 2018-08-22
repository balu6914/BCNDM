package subscriptions

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Subscription represents users purchase of stream.
type Subscription struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string        `bson:"user_id,omitempty" json:"user_id,omitempty"`
	StreamID    string        `bson:"stream_id,omitempty" json:"stream_id,omitempty"`
	StreamOwner string        `bson:"stream_owner,omitempty" json:"stream_owner,omitempty"`
	Hours       uint64        `bson:"hours,omitempty" json:"hours,omitempty"`
	StartDate   time.Time     `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate     time.Time     `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StreamURL   string        `bson:"stream_url,omitempty" json:"stream_url,omitempty"`
}

// Validate returns an error if user representation is invalid.
func (sub *Subscription) Validate() error {
	if sub.Hours <= 0 || !bson.IsObjectIdHex(sub.StreamID) || !bson.IsObjectIdHex(sub.UserID) {
		return ErrMalformedEntity
	}

	return nil
}

// Page represents paged result for list response.
type Page struct {
	Page    uint64         `json:"page"`
	Limit   uint64         `json:"limit"`
	Total   uint64         `json:"total"`
	Content []Subscription `json:"content"`
}

// SubscriptionRepository specifies a subscription persistence API.
type SubscriptionRepository interface {
	// Save persists the subscription.
	Save(Subscription) (string, error)

	// Search retrieves a list of subscription by the Query.
	Search(Query) (Page, error)

	// One retrieves a subscription by its ID.
	One(string) (Subscription, error)
}

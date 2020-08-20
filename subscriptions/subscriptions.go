package subscriptions

import (
	"strconv"
	"time"

	"github.com/datapace/datapace/auth"
	"gopkg.in/mgo.v2/bson"
)

var _ auth.Resource = (*Subscription)(nil)

// Subscription represents users purchase of stream.
type Subscription struct {
	ID          bson.ObjectId `json:"id,omitempty"`
	UserID      string        `json:"user_id,omitempty"`
	StreamID    string        `json:"stream_id,omitempty"`
	StreamOwner string        `json:"stream_owner,omitempty"`
	Hours       uint64        `json:"hours,omitempty"`
	StartDate   time.Time     `json:"start_date,omitempty"`
	EndDate     time.Time     `json:"end_date,omitempty"`
	StreamURL   string        `json:"stream_url,omitempty"`
	StreamPrice uint64        `json:"stream_price,omitempty"`
	StreamName  string        `json:"stream_name,omitempty"`
}

// Attributes returns authz.Resource attributes.
func (sub Subscription) Attributes() map[string]string {
	return map[string]string{
		"id":          sub.ID.String(),
		"ownerID":     sub.UserID,
		"streamID   ": sub.StreamID,
		"streamOwner": sub.StreamOwner,
		"hours":       strconv.FormatUint(sub.Hours, 10),
		"startDate":   sub.StartDate.String(),
		"endDate":     sub.EndDate.String(),
		"streamURL":   sub.StreamURL,
		"streamPrice": strconv.FormatUint(sub.StreamPrice, 10),
		"streamName":  sub.StreamName,
	}
}

// ResourceType returns authz.Resource type.
func (sub Subscription) ResourceType() string {
	return "subscription"
}

const (
	minSubscriptionHours = 0
	maxSubscriptionHours = 365 * 24
)

// Validate returns an error if user representation is invalid.
func (sub *Subscription) Validate() error {
	if sub.Hours <= minSubscriptionHours || sub.Hours > maxSubscriptionHours ||
		!bson.IsObjectIdHex(sub.StreamID) || !bson.IsObjectIdHex(sub.UserID) {
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

	// OneByUserAndStream retrieves a subscription by its owner ID
	// and stream ID.
	OneByUserAndStream(string, string) (Subscription, error)

	// Activate subscription by ID.
	Activate(string) error

	// Removes subscription by ID.
	Remove(string) error
}

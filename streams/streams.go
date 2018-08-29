package streams

import (
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

var _ Service = (*streamService)(nil)

// Location represents Stream location to enable geo
// search streams. Official MongoDB docs could be found here
// http://docs.mongoengine.org/guide/querying.html#geo-queries
type Location struct {
	Type string `json:"type,omitempty"`
	// Coordinates represent latitude and longitude. It's represented this
	// way to match the way MongoDB represents geo data.
	Coordinates [2]float64 `json:"coordinates,omitempty"`
}

// Stream represents data stream to be exchanged through platform.
type Stream struct {
	Owner       string        `bson:"owner,omitempty" json:"owner,omitempty"`
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Type        string        `bson:"type,omitempty" json:"type,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Snippet     string        `bson:"snippet,omitempty" json:"snippet,omitempty"`
	URL         string        `bson:"url,omitempty" json:"url,omitempty"`
	Price       uint64        `bson:"price,omitempty" json:"price,omitempty"`
	Location    Location      `bson:"location,omitempty" json:"location,omitempty"`
}

// Page represents paged result for list response.
type Page struct {
	Page    uint64   `json:"page"`
	Limit   uint64   `json:"limit"`
	Total   uint64   `json:"total"`
	Content []Stream `json:"content"`
}

// Validate returns an error if user representation is invalid.
func (s *Stream) Validate() error {
	if s.Name == "" || s.Type == "" ||
		s.Description == "" || s.URL == "" {
		return ErrMalformedData
	}

	if !govalidator.IsURL(s.URL) {
		return ErrMalformedData
	}

	if s.ID != "" && !bson.IsObjectIdHex(s.Owner) {
		return ErrMalformedData
	}

	if s.Owner != "" && !bson.IsObjectIdHex(s.Owner) {
		return ErrMalformedData
	}

	return nil
}

// StreamRepository specifies a stream persistence API.
type StreamRepository interface {
	// Save persists the stream. A non-nil error is returned to indicate
	// operation failure.
	Save(Stream) (string, error)

	// Save persists an array of streams. A non-nil error is returned to
	// indicate operation failure.
	SaveAll([]Stream) error

	// Search for streams by specified query parameters.
	Search(Query) (Page, error)

	// Update performs an update of an existing stream. A non-nil error is
	// returned to indicate operation failure.
	Update(Stream) error

	// One retrieves a stream by its unique identifier (i.e. id).
	One(string) (Stream, error)

	// Removes the stream having the provided identifier, that is owned
	// by the specified user.
	Remove(string, string) error
}

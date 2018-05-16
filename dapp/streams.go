package dapp

import (
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

type Location struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}
type Stream struct {
	Owner       string        `bson:"owner,omitempty" json:"owner,omitempty"`
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Type        string        `bson:"type,omitempty" json:"type,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	URL         string        `bson:"url,omitempty" json:"url,omitempty"`
	Price       int           `bson:"price,omitempty" json:"price,omitempty"`
	Location    Location      `bson:"location,omitempty" json:"location,omitempty"`
}

// Validate returns an error if user representation is invalid.
func (s *Stream) Validate() error {
	if s.Name == "" || s.Type == "" ||
		s.Description == "" || s.URL == "" {
		return ErrMalformedEntity
	}

	if !govalidator.IsURL(s.URL) {
		return ErrMalformedEntity
	}

	return nil
}

// StreamRepository specifies a stream persistence API.
type StreamRepository interface {
	// Save persists the stream. A non-nil error is returned to indicate
	// operation failure.
	Save(Stream) (string, error)

	// Update performs an update of an existing stream. A non-nil error is
	// returned to indicate operation failure.
	Update(string, Stream) error

	// One retrieves a stream by its unique identifier (i.e. name).
	One(string) (Stream, error)

	// Search for streams by means of geolocation parameters.
	Search([][]float64) ([]Stream, error)

	// Removes the stream having the provided identifier, that is owned
	// by the specified user.
	Remove(string) error
}

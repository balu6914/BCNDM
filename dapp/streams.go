package dapp

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict error = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity error = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess error = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound error = errors.New("non-existent entity")

	ErrUnknownType error = errors.New("unknown type")

	ErrMalformedData error = errors.New("malformed data")

	ErrUnsupportedContentType error = errors.New("unsupported content type")
)

// type Location struct {
// 	Longitude float32
// 	Latitude  float32
// }

type Location struct {
	Type        string
	Coordinates []float64
}
type Stream struct {
	Name        string
	Type        string
	Description string
	URL         string
	Price       int
	// Owner       User
	Location Location
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
	Save(Stream) error

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

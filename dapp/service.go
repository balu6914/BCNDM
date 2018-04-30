package dapp

import (
	"errors"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict error = errors.New("stream id already taken")

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

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Adds new stream to the user identified by the provided key.
	AddStream(string, Stream) (string, error)

	// Updates the stream identified by the provided ID, that
	// belongs to the user identified by the provided key.
	UpdateStream(string, string, Stream) error

	// Retrieves data about the stream identified with the provided
	// ID, that belongs to the user identified by the provided key.
	ViewStream(string, string) (Stream, error)

	// Retrieves data about subset of streams
	// given geolocation coordinates.
	SearchStreams([][]float64) ([]Stream, error)

	// Removes the stream identified with the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveStream(string, string) error
}

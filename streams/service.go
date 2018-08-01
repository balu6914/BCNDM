package streams

import (
	"errors"
)

var (
	// ErrConflict indicates usage of the existing stream id for the new stream.
	ErrConflict = errors.New("stream id already taken")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrWrongType indicates wrong contant type error.
	ErrWrongType = errors.New("wrong type")

	// ErrMalformedData indicates a malformed request.
	ErrMalformedData = errors.New("malformed data")
)

var _ Service = (*streamService)(nil)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Adds new stream to the user identified by the provided user id.
	AddStream(string, Stream) (string, error)

	// Adds new streams via parsed csv file.
	AddBulkStream(string, []Stream) error

	// Updates the stream identified by the provided id, that
	// belongs to the user identified by the provided id.
	UpdateStream(string, Stream) error

	// Retrieves data about the stream identified.
	ViewStream(string) (Stream, error)

	// Retrieves data about subset of streams given geolocation coordinates.
	SearchStreams(Query) (Page, error)

	// Removes the stream identified with the provided id, that
	// belongs to the user identified by the provided id.
	RemoveStream(string, string) error
}

type streamService struct {
	streams StreamRepository
}

// NewService instantiates the domain service implementation.
func NewService(streams StreamRepository) Service {
	return streamService{
		streams: streams,
	}
}

func (ss streamService) AddStream(owner string, stream Stream) (string, error) {
	return ss.streams.Save(stream)
}

func (ss streamService) AddBulkStream(owner string, streams []Stream) error {
	return ss.streams.SaveAll(streams)
}

func (ss streamService) UpdateStream(owner string, stream Stream) error {
	return ss.streams.Update(stream)
}

func (ss streamService) ViewStream(id string) (Stream, error) {
	s, err := ss.streams.One(id)
	if err != nil {
		return s, err
	}

	s.URL = ""
	return s, nil
}

func (ss streamService) SearchStreams(query Query) (Page, error) {
	p, err := ss.streams.Search(query)
	if err != nil {
		return p, err
	}

	// Workaround to prevent sending real URL to the end user.
	for i := range p.Content {
		p.Content[i].URL = ""
	}

	return p, nil
}

func (ss streamService) RemoveStream(owner string, id string) error {
	return ss.streams.Remove(owner, id)
}

package streams

import "errors"

var (
	// ErrConflict indicates usage of the existing stream id
	// for the new stream.
	ErrConflict = errors.New("stream id already taken")

	// ErrUnauthorizedAccess indicates missing or invalid
	// credentials provided when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrWrongType indicates wrong contant type error.
	ErrWrongType = errors.New("wrong type")

	// ErrMalformedData indicates a malformed request.
	ErrMalformedData = errors.New("malformed data")
)

var _ Service = (*streamService)(nil)

// Service specifies an API that must be fullfiled by the
// domain service implementation, and all of its decorators
// (e.g. logging & metrics).
type Service interface {
	// Adds new stream to the user identified by the provided id.
	AddStream(Stream) (string, error)

	// Adds new streams via parsed csv file.
	AddBulkStreams([]Stream) error

	// Retrieves data about subset of streams given geolocation
	// coordinates, name, type, owner or price range. Data is returned
	// in the Page form. Provides check if the user is actual owner of
	// the Stream to prevent access to the real Stream URL.
	SearchStreams(string, Query) (Page, error)

	// Updates the Stream identified by the provided id.
	UpdateStream(Stream) error

	// Retrieves data about the Stream identified by the id.
	// Provides check if the user is actual owner of the
	// Stream to prevent access to the real Stream URL.
	ViewStream(string, string) (Stream, error)

	// Removes the Stream identified with the provided id, that
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

func (ss streamService) AddStream(stream Stream) (string, error) {
	return ss.streams.Save(stream)
}

func (ss streamService) AddBulkStreams(streams []Stream) error {
	return ss.streams.SaveAll(streams)
}

func (ss streamService) SearchStreams(owner string, query Query) (Page, error) {
	p, err := ss.streams.Search(query)
	if err != nil {
		return p, err
	}

	// Prevent sending real URL to the end user.
	for i := range p.Content {
		if p.Content[i].Owner != owner {
			p.Content[i].URL = ""
		}
	}

	return p, nil
}

func (ss streamService) UpdateStream(stream Stream) error {
	return ss.streams.Update(stream)
}

func (ss streamService) ViewStream(id, owner string) (Stream, error) {
	s, err := ss.streams.One(id)
	if err != nil {
		return s, err
	}

	if s.Owner != owner {
		s.URL = ""
	}
	return s, nil
}

func (ss streamService) RemoveStream(owner string, id string) error {
	return ss.streams.Remove(owner, id)
}

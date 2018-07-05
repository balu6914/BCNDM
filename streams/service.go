package streams

import "errors"

const defLocation = "Point"

var (
	// ErrConflict indicates usage of the existing stream id for the new stream.
	ErrConflict = errors.New("stream id already taken")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrUnknownType indicates a unknown type error.
	ErrUnknownType = errors.New("unknown type")

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
	SearchStreams([][]float64) ([]Stream, error)

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
	stream.Owner = owner
	if stream.Location.Type == "" {
		stream.Location.Type = defLocation
	}

	return ss.streams.Save(stream)
}

func (ss streamService) AddBulkStream(owner string, streams []Stream) error {
	if len(streams) < 1 {
		return ErrMalformedData
	}
	for _, stream := range streams {
		stream.Owner = owner
		if stream.Location.Type == "" {
			stream.Location.Type = defLocation
		}
	}

	return ss.streams.SaveAll(streams)
}

func (ss streamService) UpdateStream(owner string, stream Stream) error {
	stream.Owner = owner
	if stream.Location.Type == "" {
		stream.Location.Type = defLocation
	}

	return ss.streams.Update(stream)
}

func (ss streamService) ViewStream(id string) (Stream, error) {
	return ss.streams.One(id)
}

func (ss streamService) SearchStreams(coords [][]float64) ([]Stream, error) {
	return ss.streams.Search(coords)
}

func (ss streamService) RemoveStream(owner string, id string) error {
	return ss.streams.Remove(owner, id)
}

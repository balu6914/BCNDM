package dapp

type Location struct {
	Longitude float32
	Latitude  float32
}

type Stream struct {
	Name        string
	Type        string
	Description string
	URL         string
	Price       int
	// Owner       User
	// Longlat     Location
}

// StreamRepository specifies a stream persistence API.
type StreamRepository interface {
	// Save persists the stream. A non-nil error is returned to indicate
	// operation failure.
	Save(Stream) (Stream, error)

	// Update performs an update of an existing stream. A non-nil error is
	// returned to indicate operation failure.
	Update(string, Stream) error

	// One retrieves a stream by its unique identifier (i.e. name).
	One(string) (Stream, error)

	// Search for streams by means of geolocation parameters.
	// Search([]int) ([]Stream, error)

	// Removes the stream having the provided identifier, that is owned
	// by the specified user.
	Remove(string) error
}

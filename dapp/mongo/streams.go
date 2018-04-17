package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
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
)

// type Location struct {
// 	Longitude float32
// 	Latitude  float32
// }

// type Stream struct {
// 	Name        string
// 	Type        string
// 	Description string
// 	URL         string
// 	Price       int
// 	// Owner       User
// 	// Longlat     Location
// }

// // StreamRepository specifies a stream persistence API.
// type StreamRepository interface {
// 	// Save persists the stream. A non-nil error is returned to indicate
// 	// operation failure.
// 	Save(Stream) (Stream, error)

// 	// Update performs an update of an existing stream. A non-nil error is
// 	// returned to indicate operation failure.
// 	Update(string, Stream) error

// 	// One retrieves a stream by its unique identifier (i.e. name).
// 	One(string) (Stream, error)

// 	// Search for streams by means of geolocation parameters.
// 	// Search([]int) ([]Stream, error)

// 	// Removes the stream having the provided identifier, that is owned
// 	// by the specified user.
// 	Remove(string) error
// }

type streamRepository struct {
	db *mgo.Session
}

// Ensure a type streamRepository implements an interface StreamRepository
var _ StreamRepository = (*streamRepository)(nil)

func NewStreamRepository(db *mgo.Session) *streamRepository {
	return &streamRepository{db}
}

func (sr streamRepository) Save(stream Stream) (Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Insert(stream); err != nil {
		if mgo.IsDup(err) {
			return stream, ErrConflict
		}

		return stream, err
	}

	return stream, nil
}

func (sr streamRepository) Update(name string, stream Stream) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	query := bson.M{"name": name}
	update := bson.M{"$set": stream}
	if err := c.Update(query, update); err != nil {
		return err
	}

	return nil
}

func (sr streamRepository) One(name string) (Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	stream := Stream{}

	if err := c.Find(bson.M{"name": name}).One(&stream); err != nil {
		if err == mgo.ErrNotFound {
			return stream, ErrNotFound
		}

		return stream, err
	}

	return stream, nil
}

// func (sr streamRepository) Search(coords []int) ([]Stream, error) {
// 	fmt.Println("Search")
// 	fmt.Println(req)

// 	count := 10
// 	streams := make([]Stream, count)
// 	for i := 0; i < count; i++ {
// 		streams[i] = Stream{
// 			Name:        fmt.Sprintf("Name: %d", i),
// 			Type:        fmt.Sprintf("Type: %d", i),
// 			Description: fmt.Sprintf("Description: %d", i),
// 			Price:       i,
// 		}
// 	}

// 	return streams, nil
// }

func (sr streamRepository) Remove(name string) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Remove(bson.M{"name": name}); err != nil {
		return err
	}

	return nil
}

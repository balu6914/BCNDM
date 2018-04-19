package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"monetasa/auth"
	"monetasa/dapp"
)

type streamRepository struct {
	db *mgo.Session
}

// Ensure a type streamRepository implements an interface StreamRepository
var _ dapp.StreamRepository = (*streamRepository)(nil)

func NewStreamRepository(db *mgo.Session) *streamRepository {
	return &streamRepository{db}
}

func (sr streamRepository) Save(stream dapp.Stream) (dapp.Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Insert(stream); err != nil {
		if mgo.IsDup(err) {
			return stream, auth.ErrConflict
		}

		return stream, err
	}

	return stream, nil
}

func (sr streamRepository) Update(id string, stream dapp.Stream) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	query := bson.M{"_id": _id}
	update := bson.M{"$set": stream}
	if err := c.Update(query, update); err != nil {
		return err
	}

	return nil
}

func (sr streamRepository) One(id string) (dapp.Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	stream := dapp.Stream{}

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Find(bson.M{"_id": _id}).One(&stream); err != nil {
		if err == mgo.ErrNotFound {
			return stream, auth.ErrNotFound
		}

		return stream, err
	}

	return stream, nil
}

func (sr streamRepository) Remove(id string) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Remove(bson.M{"_id": _id}); err != nil {
		return err
	}

	return nil
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

package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"monetasa/dapp"
)

type streamRepository struct {
	db *mgo.Session
}

// Ensure a type streamRepository implements an interface StreamRepository
var _ dapp.StreamRepository = (*streamRepository)(nil)

func NewStreamRepository(db *mgo.Session) dapp.StreamRepository {
	return &streamRepository{db}
}

func (sr streamRepository) Save(stream dapp.Stream) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Insert(stream); err != nil {
		if mgo.IsDup(err) {
			return dapp.ErrConflict
		}

		return err
	}

	return nil
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
			return stream, dapp.ErrNotFound
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
		if err == mgo.ErrNotFound {
			return dapp.ErrNotFound
		}
		return err
	}

	return nil
}

func (sr streamRepository) Search(coords [][]float64) ([]dapp.Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	var results []dapp.Stream
	err := c.Find(bson.M{
		"location": bson.M{
			"$geoWithin": bson.M{
				"$polygon": coords,
			},
		},
	}).All(&results)

	if err != nil {
		return nil, dapp.ErrNotFound
	}

	return results, nil
}

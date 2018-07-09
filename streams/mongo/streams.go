package mongo

import (
	"monetasa/streams"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName     string = "monetasa"
	collection string = "streams"
)

var _ streams.StreamRepository = (*streamRepository)(nil)

type streamRepository struct {
	db *mgo.Session
}

// New instantiates a Mongo implementation of streams
// repository.
func New(db *mgo.Session) streams.StreamRepository {
	c := db.DB(dbName).C(collection)
	ownersIdx := mgo.Index{
		Name: "owners",
		Key:  []string{"owner"},
	}
	locIdx := mgo.Index{
		Name: "location",
		Key:  []string{"$2d:location.coordinates"},
	}
	c.EnsureIndex(ownersIdx)
	c.EnsureIndex(locIdx)
	return &streamRepository{db}
}

func (sr streamRepository) Save(stream streams.Stream) (string, error) {
	if !bson.IsObjectIdHex(stream.Owner) {
		return "", streams.ErrMalformedData
	}

	s := sr.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(collection)
	stream.ID = bson.NewObjectId()

	if err := c.Insert(stream); err != nil {
		if mgo.IsDup(err) {
			return "", streams.ErrConflict
		}
		return "", err
	}

	return stream.ID.Hex(), nil
}

func (sr streamRepository) SaveAll(blk []streams.Stream) error {
	if len(blk) == 0 {
		return nil
	}

	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)
	bulk := c.Bulk()

	var arr []interface{}
	for _, stream := range blk {
		arr = append(arr, &stream)
	}

	bulk.Insert(arr...)
	if _, err := bulk.Run(); err != nil {
		return err
	}

	return nil
}

func (sr streamRepository) Update(stream streams.Stream) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	query := bson.M{"_id": stream.ID, "owner": stream.Owner}
	update := bson.M{"$set": stream}
	if err := c.Update(query, update); err != nil {
		if err == mgo.ErrNotFound {
			return streams.ErrNotFound
		}
		return err
	}

	return nil
}

func (sr streamRepository) One(id string) (streams.Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	stream := streams.Stream{}

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Find(bson.M{"_id": _id}).One(&stream); err != nil {
		if err == mgo.ErrNotFound {
			return stream, streams.ErrNotFound
		}

		return stream, err
	}

	return stream, nil
}

func (sr streamRepository) Search(coords [][]float64) ([]streams.Stream, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	var results []streams.Stream
	err := c.Find(bson.M{
		"location": bson.M{
			"$within": bson.M{
				"$polygon": coords,
			},
		},
	}).All(&results)

	if results == nil || err != nil {
		return []streams.Stream{}, nil
	}

	return results, nil
}

func (sr streamRepository) Remove(owner, id string) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Remove(bson.M{"_id": _id, "owner": owner}); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}

	return nil
}

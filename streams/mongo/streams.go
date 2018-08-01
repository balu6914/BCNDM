package mongo

import (
	"monetasa/streams"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName     string = "monetasa-streams"
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
		Name: "locations",
		Key:  []string{"$2d:location.coordinates"},
	}
	namesIdx := mgo.Index{
		Name: "names",
		Key:  []string{"name"},
	}
	typeIdx := mgo.Index{
		Name: "types",
		Key:  []string{"type"},
	}
	priceIdx := mgo.Index{
		Name: "prices",
		Key:  []string{"price"},
	}

	c.EnsureIndex(ownersIdx)
	c.EnsureIndex(locIdx)
	c.EnsureIndex(namesIdx)
	c.EnsureIndex(typeIdx)
	c.EnsureIndex(priceIdx)
	return &streamRepository{db}
}

func (sr streamRepository) Save(stream streams.Stream) (string, error) {
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

	arr := make([]interface{}, len(blk))
	for i := range blk {
		arr[i] = blk[i]
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

func (sr streamRepository) Search(query streams.Query) (streams.Page, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)
	limit := int(query.Limit)
	page := int(query.Page)

	ret := streams.Page{
		Page:    query.Page,
		Limit:   query.Limit,
		Content: []streams.Stream{},
	}

	var results []streams.Stream
	q := streams.GenerateQuery(&query)

	total, err := c.Find(q).Count()
	if err != nil {
		return ret, err
	}

	ret.Total = uint64(total)
	start := limit * page
	if total < start {
		return ret, nil
	}

	err = c.Find(q).Skip(start).Limit(limit).All(&results)
	if results == nil || err != nil {
		return ret, nil
	}

	ret.Content = results
	return ret, nil
}

func (sr streamRepository) Remove(owner, id string) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	removeID := bson.ObjectIdHex(id)
	if err := c.Remove(bson.M{"_id": removeID, "owner": owner}); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}

	return nil
}

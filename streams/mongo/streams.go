package mongo

import (
	"strings"

	"github.com/datapace/datapace/errors"
	"github.com/datapace/datapace/streams"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName     = "datapace-streams"
	collection = "streams"
	unknown    = "unknown conflict"
	errMsg     = "Some of the URLs already exist in the database."
	msgMark    = "\""
)

var errFormatUUID = errors.New("invalid uuid format")

var _ streams.StreamRepository = (*streamRepository)(nil)

type streamRepository struct {
	db *mgo.Session
}

// This obscure way of parsing error is forced by the
// way `mgo` handles Bulk error.
func parseError(err error) string {
	e := err.Error()
	start := strings.Index(e, msgMark)
	if start != -1 && start < len(e)-1 {
		e = e[start+1:]
		if end := strings.Index(e, msgMark); end != -1 {
			return e[:end]
		}
	}

	return unknown
}

// New instantiates a Mongo implementation of streams
// repository.
func New(db *mgo.Session) streams.StreamRepository {
	c := db.DB(dbName).C(collection)
	indices := []mgo.Index{
		mgo.Index{
			Name: "owners",
			Key:  []string{"owner"},
		},
		mgo.Index{
			Name: "visibilities",
			Key:  []string{"visibility"},
		},
		mgo.Index{
			Name: "locations",
			Key:  []string{"$2d:location.coordinates"},
		},
		mgo.Index{
			Name:   "urls",
			Unique: true,
			Key:    []string{"url"},
		},
		mgo.Index{
			Name: "names",
			Key:  []string{"name"},
		},
		mgo.Index{
			Name: "types",
			Key:  []string{"type"},
		},
		mgo.Index{
			Name: "prices",
			Key:  []string{"price"},
		},
		mgo.Index{
			Name: "metadatas",
			Key:  []string{"metadata"},
		},
		mgo.Index{
			Name: "terms",
			Key:  []string{"terms"},
		},
	}
	for _, idx := range indices {
		c.EnsureIndex(idx)
	}

	return &streamRepository{db}
}

func (sr streamRepository) Save(stream streams.Stream) (string, error) {
	s := sr.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(collection)
	// ignore error because invalid ID should be ignored in this case
	dbs, _ := toDBStream(stream)
	dbs.ID = bson.NewObjectId()

	if err := c.Insert(dbs); err != nil {
		if mgo.IsDup(err) {
			return "", streams.ErrConflict
		}
		return "", err
	}

	return dbs.ID.Hex(), nil
}

func (sr streamRepository) SaveAll(blk []streams.Stream) error {
	if len(blk) == 0 {
		return nil
	}

	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)
	bulk := c.Bulk()
	bulk.Unordered()

	arr := make([]interface{}, len(blk))
	for i, v := range blk {
		// ignore error because invalid ID should be ignored in this case
		dbs, _ := toDBStream(v)
		dbs.ID = bson.NewObjectId()

		arr[i] = dbs
	}

	bulk.Insert(arr...)
	if _, err := bulk.Run(); err != nil {
		if mgo.IsDup(err) {
			ret := streams.ErrBulkConflict{
				Message:   errMsg,
				Conflicts: []string{},
			}
			if bulkErr, ok := err.(*mgo.BulkError); ok {
				for _, errCase := range bulkErr.Cases() {
					ret.Conflicts = append(ret.Conflicts, parseError(errCase.Err))
				}
			}
			return ret
		}
		return err
	}

	return nil
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

	var results []dbStream
	q := streams.GenQuery(&query)

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

	strms := []streams.Stream{}
	for _, dbs := range results {
		strms = append(strms, fromDBStream(dbs))
	}

	ret.Content = strms
	return ret, nil
}

func (sr streamRepository) Update(stream streams.Stream) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collection)

	dbs, err := toDBStream(stream)
	if err != nil {
		return err
	}

	query := bson.M{"_id": dbs.ID, "owner": dbs.Owner}
	update := bson.M{"$set": dbs}
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

	dbs := dbStream{}

	if !bson.IsObjectIdHex(id) {
		return streams.Stream{}, errFormatUUID
	}

	// It's valid, calling bson.ObjectIdHex() will not panic...
	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Find(bson.M{"_id": _id}).One(&dbs); err != nil {
		if err == mgo.ErrNotFound {
			return streams.Stream{}, streams.ErrNotFound
		}
		return streams.Stream{}, err
	}

	stream := fromDBStream(dbs)
	return stream, nil
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

type dbStream struct {
	Owner       string             `bson:"owner,omitempty"`
	ID          bson.ObjectId      `bson:"_id,omitempty"`
	Visibility  streams.Visibility `bson:"visibility,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Type        string             `bson:"type,omitempty"`
	Description string             `bson:"description,omitempty"`
	Snippet     string             `bson:"snippet,omitempty"`
	URL         string             `bson:"url,omitempty"`
	Price       uint64             `bson:"price,omitempty"`
	Location    dbLocation         `bson:"location,omitempty"`
	Terms       string             `bson:"terms,omitempty"`
	External    bool               `bson:"external"`
	BQ          dbBigQuery         `bson:"big_query,omitempty"`
	Metadata    bson.M             `bson:"metadata,omitempty"`
}

type dbBigQuery struct {
	// Email represents Gmail address of the owner. It can be either Email
	// or ContactEmail of the owner.
	Email   string `bson:"email,omitempty"`
	Project string `bson:"project,omitempty"`
	Dataset string `bson:"dataset,omitempty"`
	Table   string `bson:"table,omitempty"`
	Fields  string `bson:"fields,omitempty"`
}

type dbLocation struct {
	Type string `bson:"type,omitempty"`
	// Coordinates represent longitude and latitude. It's represented this
	// way to match the way MongoDB represents geo data.
	Coordinates [2]float64 `bson:"coordinates,omitempty"`
}

func toDBStream(stream streams.Stream) (dbStream, error) {
	dbl := dbLocation{
		Type:        stream.Location.Type,
		Coordinates: stream.Location.Coordinates,
	}

	dbBQ := dbBigQuery{
		Email:   stream.BQ.Email,
		Project: stream.BQ.Project,
		Dataset: stream.BQ.Dataset,
		Table:   stream.BQ.Table,
		Fields:  stream.BQ.Fields,
	}

	dbs := dbStream{
		Owner:       stream.Owner,
		Visibility:  stream.Visibility,
		Name:        stream.Name,
		Type:        stream.Type,
		Description: stream.Description,
		Snippet:     stream.Snippet,
		URL:         stream.URL,
		Price:       stream.Price,
		Location:    dbl,
		Terms:       stream.Terms,
		External:    stream.External,
		BQ:          dbBQ,
		Metadata:    stream.Metadata,
	}

	if ok := bson.IsObjectIdHex(stream.ID); !ok {
		return dbs, streams.ErrMalformedData
	}
	dbs.ID = bson.ObjectIdHex(stream.ID)

	return dbs, nil
}

func fromDBStream(dbs dbStream) streams.Stream {
	dbl := streams.Location{
		Type:        dbs.Location.Type,
		Coordinates: dbs.Location.Coordinates,
	}

	dbBQ := streams.BigQuery{
		Email:   dbs.BQ.Email,
		Project: dbs.BQ.Project,
		Dataset: dbs.BQ.Dataset,
		Table:   dbs.BQ.Table,
		Fields:  dbs.BQ.Fields,
	}

	return streams.Stream{
		Owner:       dbs.Owner,
		ID:          dbs.ID.Hex(),
		Visibility:  dbs.Visibility,
		Name:        dbs.Name,
		Type:        dbs.Type,
		Description: dbs.Description,
		Snippet:     dbs.Snippet,
		URL:         dbs.URL,
		Price:       dbs.Price,
		Location:    dbl,
		Terms:       dbs.Terms,
		External:    dbs.External,
		BQ:          dbBQ,
		Metadata:    dbs.Metadata,
	}
}

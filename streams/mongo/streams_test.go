package mongo_test

import (
	"fmt"
	log "monetasa/logger"
	"monetasa/streams"
	"monetasa/streams/mongo"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const url = "localhost"

var (
	db         *mgo.Session
	testLog            = log.New(os.Stdout)
	testDB             = "monetasa"
	collection         = "streams"
	long       float64 = 50
	lat        float64 = 50
)

func TestSave(t *testing.T) {
	repo := mongo.New(db)

	id, err := repo.Save(stream())

	assert.Nil(t, err, fmt.Sprintf("expected to save successfully got: %s", err))
	assert.True(t, bson.IsObjectIdHex(id), fmt.Sprintf("create new stream expected to return a valid ID"))
}

func TestSaveAll(t *testing.T) {
	repo := mongo.New(db)
	s := stream()
	s.ID = ""

	bulk := []streams.Stream{}
	for i := 0; i < 100; i++ {
		bulk = append(bulk, s)
	}

	err := repo.SaveAll(bulk)
	assert.Nil(t, err, fmt.Sprintf("expected to save all successfully got: %s", err))
}

func TestOne(t *testing.T) {
	repo := mongo.New(db)
	s := stream()
	id, err := repo.Save(stream())
	s.ID = bson.ObjectIdHex(id)
	assert.Nil(t, err, fmt.Sprintf("create new stream: expected no error, got %s", err))

	nonexisting := s
	nonexisting.ID = bson.NewObjectId()

	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "get an existing stream",
			stream: s,
			err:    nil,
		},
		{
			desc:   "get a non-existing stream",
			stream: nonexisting,
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := repo.One(tc.stream.ID.Hex())
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	repo := mongo.New(db)
	s := stream()
	id, err := repo.Save(s)
	s.ID = bson.ObjectIdHex(id)
	assert.Nil(t, err, fmt.Sprintf("create a new stream: expected no error, got %s", err))

	nonExisting := s
	nonExisting.ID = bson.NewObjectId()

	wrong := s
	wrong.Owner = bson.NewObjectId().Hex()

	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "update an existing stream",
			stream: s,
			err:    nil,
		},
		{
			desc:   "update an existing stream with wrong user ID",
			stream: wrong,
			err:    streams.ErrNotFound,
		},
		{
			desc:   "update a non-existing stream",
			stream: nonExisting,
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := repo.Update(tc.stream)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSearch(t *testing.T) {
	repo := mongo.New(db)
	repo.Save(stream())

	dbSize, _ := db.DB(testDB).C(collection).Count()

	cases := []struct {
		desc   string
		coords [][]float64
		size   int
		err    error
	}{
		{
			desc:   "search get all streams",
			coords: [][]float64{{-180, -90}, {-180, 90}, {180, 90}, {180, -90}},
			size:   dbSize,
			err:    nil,
		},
		{
			desc:   "search empty result",
			coords: [][]float64{{long - 1, lat - 1}, {long, lat - 2}, {long - 3, lat}},
			size:   0,
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		resp, err := repo.Search(tc.coords)
		n := len(resp)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.size, n, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, n))
	}
}

func TestRemove(t *testing.T) {
	repo := mongo.New(db)

	s := stream()
	id, err := repo.Save(s)
	s.ID = bson.ObjectIdHex(id)
	assert.Nil(t, err, fmt.Sprintf("create new stream: expected no error, got %s", err))

	// show that the removal works the same for both existing and non-existing
	// (removed) thing
	for i := 0; i < 2; i++ {
		err := repo.Remove(s.Owner, s.ID.Hex())
		assert.Nil(t, err, "removing a stream should not return an error")

		_, err = repo.One(s.ID.Hex())
		assert.Equal(t, streams.ErrNotFound, err, fmt.Sprintf("#%d: expected %s got %s", i, streams.ErrNotFound, err))
	}
}

func stream() streams.Stream {
	return streams.Stream{
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId(),
		Name:        "stream",
		Type:        "stream type",
		Description: "stream description",
		URL:         "http://www.mystream.com",
		Price:       123,
		Location: streams.Location{
			Type:        "Point",
			Coordinates: [2]float64{long, lat},
		},
	}
}

package mongo_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	log "github.com/datapace/datapace/logger"
	"github.com/datapace/datapace/streams"
	"github.com/datapace/datapace/streams/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	url        = "localhost"
	testDB     = "datapace-streams"
	collection = "streams"
	limit      = uint64(3)
	maxPage    = uint64(100)
	long       = float64(50)
	lat        = float64(50)
	bulkSize   = 10
	searchSize = 50
)

var (
	db          *mgo.Session
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	testLog     = log.New(os.Stdout)
)

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func stream() streams.Stream {
	return streams.Stream{
		Visibility:  streams.Public,
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId().Hex(),
		Name:        "stream",
		Type:        "stream type",
		Description: "stream description",
		URL:         fmt.Sprintf("http://www.%s.com", randomString(10)),
		Price:       123,
		Location: streams.Location{
			Type:        "Point",
			Coordinates: [2]float64{long, lat},
		},
	}
}

func pointer(val uint64) *uint64 {
	return &val
}

func TestSearch(t *testing.T) {
	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	partners := []string{}
	all := []streams.Stream{}
	for i := 0; i < searchSize; i++ {
		s := stream()
		id, err := repo.Save(s)
		require.Nil(t, err, "Repo should save streams.")
		s.ID = id
		all = append(all, s)
		partners = append(partners, s.Owner)
	}
	// Specify two special Streams to match different
	// types of query and different result sets.
	s1Price := uint64(40)
	s2Price := uint64(50)

	s1 := stream()
	s1.Price = s1Price
	s1.Name = "different name"
	s1.Type = "different type"
	id, err := repo.Save(s1)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))
	s1.ID = id

	s2 := stream()
	s2.Price = s2Price
	s2.Owner = bson.NewObjectId().Hex()
	id, err = repo.Save(s2)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))
	s2.ID = id

	all = append(all, s1, s2)
	partners = append(partners, s1.Owner, s2.Owner)

	val, _ := db.DB(testDB).C(collection).Count()
	total := uint64(val)

	cases := []struct {
		desc    string
		query   streams.Query
		page    streams.Page
		content []streams.Stream
	}{
		{
			desc: "search streams with query with only the limit specified",
			query: streams.Query{
				Limit:    limit,
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   total,
				Content: all[:limit],
			},
		},
		{
			desc: "search streams reset too big offset to default value silently",
			query: streams.Query{
				Limit:    limit,
				Page:     uint64(total + maxPage),
				Partners: partners,
			},
			page: streams.Page{
				Page:    maxPage,
				Limit:   limit,
				Total:   total,
				Content: []streams.Stream{},
			},
		},
		{
			desc: "search streams with min price specified",
			// Get all except the one with the price1.
			// Content is caluclated this way because MongoDB
			// pages result from the last insertied entry.
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(s1Price + 1),
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   total - 1,
				Content: append(append([]streams.Stream{}, all[total-1]), all[:limit-1]...),
			},
		},
		{
			desc: "search streams with max price specified",
			query: streams.Query{
				Limit:    limit,
				MaxPrice: pointer(s2Price),
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   1,
				Content: []streams.Stream{s1},
			},
		},
		{
			desc: "search streams with price range specified",
			// GTE price1 and LT price2 + 1 (to include price2)
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(s1Price),
				MaxPrice: pointer(s2Price + 1),
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   2,
				Content: []streams.Stream{s1, s2},
			},
		},
		{
			desc: "search streams by owner",
			query: streams.Query{
				Limit:    limit,
				Owner:    s2.Owner,
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   1,
				Content: []streams.Stream{s2},
			},
		},
		{
			desc: "search streams by name",
			query: streams.Query{
				Limit:    limit,
				Name:     s1.Name,
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   1,
				Content: []streams.Stream{s1},
			},
		},
		{
			desc: "search streams by name partial",
			query: streams.Query{
				Limit:    limit,
				Name:     "str",
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   total - 1,
				Content: all[:limit],
			},
		},
		{
			desc: "search streams by type",
			query: streams.Query{
				Limit:      limit,
				StreamType: s1.Type,
				Partners:   partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   1,
				Content: []streams.Stream{s1},
			},
		},
		{
			desc: "search streams by owner other than provided",
			query: streams.Query{
				Limit:    limit,
				Owner:    fmt.Sprintf("-%s", s2.Owner),
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   total - 1,
				Content: all[:limit],
			},
		},
		{
			desc: "search streams by name other than provided",
			query: streams.Query{
				Limit:    limit,
				Name:     fmt.Sprintf("-%s", s2.Type[:3]),
				Partners: partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   1,
				Content: []streams.Stream{s1},
			},
		},
		{
			desc: "search streams by type other than provided",
			query: streams.Query{
				Limit:      limit,
				StreamType: fmt.Sprintf("-%s", s1.Type),
				Partners:   partners,
			},
			page: streams.Page{
				Limit:   limit,
				Total:   total - 1,
				Content: all[:limit],
			},
		},
	}

	for _, tc := range cases {
		res, err := repo.Search(tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.page.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.page.Limit, res.Limit))
		assert.Equal(t, tc.page.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.page.Total, res.Total))
		assert.ElementsMatch(t, tc.page.Content, res.Content, tc.desc)
	}
}

func TestSearchShared(t *testing.T) {

	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	all := []streams.Stream{}
	s0 := stream()
	s0.Visibility = streams.Protected
	id, err := repo.Save(s0)
	require.Nil(t, err, "Repo should save streams.")
	s0.ID = id
	all = append(all, s0)

	s1 := stream()
	s1.Visibility = streams.Protected
	s1.Name = "different name"
	s1.Type = "different type"
	s1.ID, err = repo.Save(s1)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))
	all = append(all, s1)

	cases := []struct {
		desc    string
		query   streams.Query
		page    streams.Page
		content []streams.Stream
	}{
		{
			desc: "search streams by shared stream id",
			query: streams.Query{
				Limit: limit,
				Partners: []string{
					s0.Owner,
				},
				Shared: map[string]bool{
					s1.ID: true,
				},
			},
			page: streams.Page{
				Limit: limit,
				Total: 2,
				Content: []streams.Stream{
					s0,
					s1,
				},
			},
		},
		{
			desc: "search streams by shared stream id and owner filter",
			query: streams.Query{
				Limit: limit,
				Owner: s0.Owner,
				Partners: []string{
					s0.Owner,
				},
				Shared: map[string]bool{
					s1.ID: true,
				},
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
				Content: []streams.Stream{
					s0,
				},
			},
		},
	}

	for _, tc := range cases {
		res, err := repo.Search(tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.page.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.page.Limit, res.Limit))
		assert.Equal(t, tc.page.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.page.Total, res.Total))
		assert.ElementsMatch(t, tc.page.Content, res.Content, tc.desc)
	}
}

func TestSearchByMetadata(t *testing.T) {

	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	all := []streams.Stream{}
	s0 := stream()
	s0.Visibility = streams.Public
	id, err := repo.Save(s0)
	require.Nil(t, err, "Repo should save streams.")
	s0.ID = id
	all = append(all, s0)

	s1 := stream()
	s1.Owner = s0.Owner
	s1.Visibility = streams.Public
	s1.Name = "different name"
	s1.Metadata = map[string]interface{}{
		"foo": 42,
	}
	s1.ID, err = repo.Save(s1)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))
	all = append(all, s1)

	s2 := stream()
	s2.Owner = s0.Owner
	s2.Visibility = streams.Public
	s2.Name = "different name 2"
	s2.Metadata = bson.M{
		"nested": bson.M{
			"field1": "value1",
			"field2": bson.M{
				"field3": "value2",
			},
		},
	}
	s2.ID, err = repo.Save(s2)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))
	all = append(all, s2)

	cases := []struct {
		desc    string
		query   streams.Query
		page    streams.Page
		content []streams.Stream
	}{
		{
			desc: "search streams by present metadata",
			query: streams.Query{
				Limit: limit,
				Owner: s0.Owner,
				Partners: []string{
					s0.Owner,
				},
				Metadata: map[string]interface{}{
					"foo": 42,
				},
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
				Content: []streams.Stream{
					s1,
				},
			},
		},
		{
			desc: "search streams by empty metadata",
			query: streams.Query{
				Limit: limit,
				Owner: s0.Owner,
				Partners: []string{
					s0.Owner,
				},
			},
			page: streams.Page{
				Limit: limit,
				Total: 3,
				Content: []streams.Stream{
					s0,
					s1,
					s2,
				},
			},
		},
		{
			desc: "search streams by nested metadata",
			query: streams.Query{
				Limit: limit,
				Owner: s0.Owner,
				Partners: []string{
					s0.Owner,
				},
				Metadata: map[string]interface{}{
					"nested": map[string]interface{}{
						"field2": map[string]interface{}{
							"field3": "value2",
						},
					},
				},
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
				Content: []streams.Stream{
					s2,
				},
			},
		},
	}

	for _, tc := range cases {
		res, err := repo.Search(tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.page.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.page.Limit, res.Limit))
		assert.Equal(t, tc.page.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.page.Total, res.Total))
		assert.ElementsMatch(t, tc.page.Content, res.Content, tc.desc)
	}
}

func TestSave(t *testing.T) {
	db.DB(testDB).DropDatabase()
	db.ResetIndexCache()
	db.Refresh()

	repo := mongo.New(db)

	s := stream()
	s1 := stream()
	s1.URL = s.URL
	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "save stream successfully",
			stream: s,
			err:    nil,
		},
		{
			desc:   "save duplicate stream URL",
			stream: s1,
			err:    streams.ErrConflict,
		},
	}

	for _, tc := range cases {
		id, err := repo.Save(tc.stream)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s expected: %s", tc.desc, err))
		if err == nil {
			assert.True(t, bson.IsObjectIdHex(id), "stream id is expected to be valid")
		}
	}
}

func TestSaveAll(t *testing.T) {
	db.DB(testDB).DropDatabase()
	db.ResetIndexCache()
	db.Refresh()

	repo := mongo.New(db)

	validBulk := []streams.Stream{}
	conflictBulk := []streams.Stream{}
	conflicts := []string{}
	for i := 0; i < bulkSize; i++ {
		s := stream()
		validBulk = append(validBulk, s)
		// Add some new Streams and some Streams with
		// an existing URL, but non-existent ID.
		s.ID = bson.NewObjectId().Hex()
		conflicts = append(conflicts, s.URL)
		conflictBulk = append(conflictBulk, stream(), s)
	}

	cases := []struct {
		desc    string
		streams []streams.Stream
		err     error
	}{
		{
			desc:    "save a valid bulk",
			streams: validBulk,
			err:     nil,
		},
		{
			desc:    "save bulk with non-unique URLs",
			streams: conflictBulk,
			err: streams.ErrBulkConflict{
				// Since we don't care about message, this is exact
				// copy of the same message from the repo.
				Message:   "Some of the URLs already exist in the database.",
				Conflicts: conflicts,
			},
		},
	}

	for _, tc := range cases {
		err := repo.SaveAll(tc.streams)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error %s got %s", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	s := stream()
	id, err := repo.Save(s)
	require.Nil(t, err, fmt.Sprintf("got unexpected error when creating stream: %s", err))
	s.ID = id

	nonExisting := s
	nonExisting.ID = bson.NewObjectId().Hex()

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

func TestOne(t *testing.T) {
	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	s := stream()
	id, err := repo.Save(stream())
	s.ID = id
	assert.Nil(t, err, fmt.Sprintf("create new stream: expected no error, got %s", err))

	nonexisting := s
	nonexisting.ID = bson.NewObjectId().Hex()

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
		_, err := repo.One(tc.stream.ID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRemove(t *testing.T) {
	db.DB(testDB).DropDatabase()
	repo := mongo.New(db)

	s := stream()
	id, err := repo.Save(s)
	s.ID = id
	assert.Nil(t, err, fmt.Sprintf("create new stream: expected no error, got %s", err))

	// Show that the removal works the same for both
	// existing and non-existing (removed) stream.
	for i := 0; i < 2; i++ {
		err := repo.Remove(s.Owner, s.ID)
		assert.Nil(t, err, "removing a stream should not return an error")

		_, err = repo.One(s.ID)
		assert.Equal(t, streams.ErrNotFound, err, fmt.Sprintf("#%d: expected %s got %s", i, streams.ErrNotFound, err))
	}
}

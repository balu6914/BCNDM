package streams_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"

	"monetasa/streams"
	"monetasa/streams/mocks"
)

const (
	nameLen = 8
	typeLen = 4
	descLen = 12
	urlLen  = 6
	maxLong = 180
	maxLat  = 90
	maxPage = uint64(100)
	limit   = uint64(3)
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key         = bson.NewObjectId().Hex()
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
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId(),
		Name:        randomString(nameLen),
		Type:        randomString(typeLen),
		Description: randomString(descLen),
		URL:         fmt.Sprintf("http://%s.com", randomString(urlLen)),
		Price:       rand.Uint64(),
		Location: streams.Location{
			Type: "Point",
			Coordinates: [2]float64{
				rand.Float64() * (float64)(rand.Intn(maxLat*2)-maxLat),
				rand.Float64() * (float64)(rand.Intn(maxLong*2)-maxLong),
			},
		},
	}
}

func newService() streams.Service {
	repo := mocks.NewStreamRepository()
	return streams.NewService(repo)
}

func pointer(val uint64) *uint64 {
	return &val
}

func TestAddStream(t *testing.T) {
	svc := newService()
	s := stream()

	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "add a new stream",
			stream: s,
			err:    nil},
		{
			desc:   "add an existing stream",
			stream: s,
			err:    streams.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := svc.AddStream(tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestAddBulkStreams(t *testing.T) {
	svc := newService()
	all := []streams.Stream{}
	for i := 0; i < 100; i++ {
		all = append(all, stream())
	}

	cases := []struct {
		desc    string
		streams []streams.Stream
		key     string
		err     error
	}{
		{
			desc:    "add 100 streams",
			streams: all,
			err:     nil,
		},
		{
			desc:    "add 0 streams",
			streams: []streams.Stream{},
			err:     nil,
		},
	}

	for _, tc := range cases {
		err := svc.AddBulkStreams(tc.streams)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSearchStreams(t *testing.T) {
	svc := newService()
	size := 0
	// all := []streams.Stream{}
	for i := 0; i < 50; i++ {
		size++
		s := stream()
		id, err := svc.AddStream(s)
		s.ID = bson.ObjectIdHex(id)
		// all = append(all, s)
		require.Nil(t, err, "Repo should save streams.")
	}
	// Specify two special Streams to match different
	// types of query and different result sets.
	price1 := uint64(40)
	price2 := uint64(50)
	name := "different name"

	s1 := stream()
	s1.Price = price1
	s1.Name = name
	id, _ := svc.AddStream(s1)
	s1.ID = bson.ObjectIdHex(id)
	// all = append(all, s1)
	size++

	s2 := stream()
	owner := bson.NewObjectId().Hex()
	s2.Price = price2
	s2.Owner = owner
	id, _ = svc.AddStream(s2)
	s2.ID = bson.ObjectIdHex(id)
	// all = append(all, s2)
	size++

	total := uint64(size)
	l := int(limit)

	cases := []struct {
		desc  string
		owner string
		size  int
		query streams.Query
		page  streams.Page
	}{
		{
			desc: "search with query with only the limit specified",
			size: l,
			query: streams.Query{
				Limit: limit,
			},
			page: streams.Page{
				Limit: limit,
				Total: total,
			},
		},
		{
			desc: "search reset too big offest to default value silently",
			size: 0,
			query: streams.Query{
				Limit: limit,
				Page:  uint64(total + maxPage),
			},
			page: streams.Page{
				Page:  maxPage,
				Limit: limit,
				Total: total,
			},
		},
		{
			desc: "search with min price specified",
			size: l,
			// Get all except the one with the price1.
			// Content is caluclated this way because MongoDB
			// pages result from the last insertied entry.
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(price1 + 1),
			},
			page: streams.Page{
				Limit: limit,
				Total: total - 1,
			},
		},
		{
			desc:  "search with max price specified",
			owner: owner,
			size:  1,
			query: streams.Query{
				Limit:    limit,
				MaxPrice: pointer(price2),
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
		{
			desc: "search with price range specified",
			size: 2,
			// GTE price1 and LT price2 + 1 (to include price2)
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(price1),
				MaxPrice: pointer(price2 + 1),
			},
			page: streams.Page{
				Limit: limit,
				Total: 2,
			},
		},
		{
			desc: "search by owner",
			size: 1,
			query: streams.Query{
				Limit: limit,
				Owner: owner,
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},

		{
			desc: "search by name",
			size: 1,
			query: streams.Query{
				Limit: limit,
				Name:  name,
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
	}

	for _, tc := range cases {
		res, err := svc.SearchStreams(tc.owner, tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.page.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.page.Limit, res.Limit))
		assert.Equal(t, tc.page.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.page.Total, res.Total))
		assert.Equal(t, tc.size, len(res.Content), fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.size, len(res.Content)))
		for _, s := range res.Content {
			if tc.owner == s.Owner {
				assert.NotEmpty(t, s.URL, fmt.Sprintf("%s: expected Streams of the owner to have an URL, but got an empty value instead", tc.desc))
				continue
			}
			assert.Empty(t, s.URL, fmt.Sprintf("%s: expected Streams of the other owners not to have an URL, but got a value instead", tc.desc))
		}
	}
}

func TestUpdateStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	wrongOwner := s
	wrongOwner.Owner = bson.NewObjectId().Hex()

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
			stream: stream(),
			err:    streams.ErrNotFound,
		},
		{
			desc:   "update a stream with wrong owner",
			stream: wrongOwner,
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateStream(tc.stream)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "view an existing stream",
			id:   s.ID.Hex(),
			err:  nil,
		},
		{
			desc: "view a non-existing stream",
			id:   bson.NewObjectId().String(),
			err:  streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewStream(tc.id, "")
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRemoveStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{
			desc:     "remove an existing stream",
			streamId: s.ID,
			owner:    s.Owner,
			err:      nil,
		},
		{
			desc:     "remove a non-existing stream",
			streamId: bson.NewObjectId(),
			owner:    s.Owner,
			err:      nil,
		},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.owner, tc.streamId.Hex())
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

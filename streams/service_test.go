package streams_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
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
	key     = "token"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestAddStream(t *testing.T) {
	svc := newService()
	s := generateStream()

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
		_, err := svc.AddStream(tc.stream.Owner, tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStream(t *testing.T) {
	s := generateStream()
	svc := newService()
	svc.AddStream(key, s)

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
		_, err := svc.ViewStream(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateStream(t *testing.T) {
	s := generateStream()
	svc := newService()
	svc.AddStream(key, s)

	cases := []struct {
		desc   string
		stream streams.Stream
		owner  string
		err    error
	}{
		{
			desc:   "update an existing stream",
			stream: s,
			owner:  key,
			err:    nil,
		},
		{
			desc:   "update a non-existing stream",
			stream: generateStream(),
			owner:  key,
			err:    streams.ErrNotFound,
		},
		{
			desc:   "update an existing stream with wrong owner",
			stream: s,
			owner:  bson.NewObjectId().Hex(),
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateStream(tc.owner, tc.stream)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRemoveStream(t *testing.T) {
	s := generateStream()
	svc := newService()
	svc.AddStream(s.Owner, s)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{
			desc:     "remove existing stream with wrong owner",
			streamId: s.ID,
			owner:    "",
			err:      streams.ErrNotFound},
		{
			desc:     "remove existing stream",
			streamId: s.ID,
			owner:    s.Owner,
			err:      nil},
		{
			desc:     "remove non-existing stream",
			streamId: bson.NewObjectId(),
			owner:    s.Owner,
			err:      streams.ErrNotFound},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.owner, tc.streamId.Hex())
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestAddBulkStreams(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc    string
		streams []streams.Stream
		key     string
		err     error
	}{
		{
			desc:    "add 100 streams",
			streams: generateStreams(100),
			key:     key,
			err:     nil,
		},
		{
			desc:    "add 0 streams",
			streams: []streams.Stream{},
			key:     key,
			err:     streams.ErrMalformedData},
	}

	for _, tc := range cases {
		err := svc.AddBulkStream(tc.key, tc.streams)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateStream() streams.Stream {
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

func generateStreams(numStreams int) []streams.Stream {
	var streams []streams.Stream
	for i := 0; i < numStreams; i++ {
		streams = append(streams, generateStream())
	}
	return streams
}

func newService() streams.Service {
	repo := mocks.NewStreamRepository()
	return streams.NewService(repo)
}

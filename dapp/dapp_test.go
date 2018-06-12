package dapp_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	"monetasa/dapp"
	"monetasa/dapp/mocks"
)

const (
	nameLen  int = 8
	typeLen  int = 4
	descLen  int = 12
	urlLen   int = 6
	maxPrice int = 100
	maxLong  int = 180
	maxLat   int = 90
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateStream() dapp.Stream {
	return dapp.Stream{
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId(),
		Name:        randomString(nameLen),
		Type:        randomString(typeLen),
		Description: randomString(descLen),
		URL:         "http://" + randomString(urlLen) + ".com",
		Price:       rand.Intn(maxPrice),
		Location: dapp.Location{
			Type: "Point",
			Coordinates: []float64{
				rand.Float64() * (float64)(rand.Intn(maxLong*2)-maxLong),
				rand.Float64() * (float64)(rand.Intn(maxLat*2)-maxLat)},
		},
	}
}

func generateStreams(numStreams int) []dapp.Stream {
	var streams []dapp.Stream
	for i := 0; i < numStreams; i++ {
		streams = append(streams, generateStream())
	}
	return streams
}

func newService() dapp.Service {
	streams := mocks.NewStreamRepository()
	subs := mocks.NewSubscriptionsRepository()
	return dapp.New(streams, subs)
}

func TestAddStream(t *testing.T) {
	svc := newService()
	s := generateStream()

	cases := []struct {
		desc   string
		stream dapp.Stream
		err    error
	}{
		{"add new stream", s, nil},
		{"add existing stream", s, dapp.ErrConflict},
	}

	for _, tc := range cases {
		_, err := svc.AddStream(tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestAddBulkStreams(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc    string
		streams []dapp.Stream
		err     error
	}{
		{"add 100 streams", generateStreams(100), nil},
		{"add 0 streams", []dapp.Stream{}, dapp.ErrMalformedData},
	}

	for _, tc := range cases {
		err := svc.AddBulkStream(tc.streams)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateStream(t *testing.T) {
	svc := newService()
	s := generateStream()
	svc.AddStream(s)

	cases := []struct {
		desc     string
		stream   dapp.Stream
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{"update existing stream", generateStream(), s.ID, s.Owner, nil},
		{"update non-existing stream", generateStream(), bson.NewObjectId(),
			s.Owner, dapp.ErrNotFound},
		{"update existing stream with wrong owner", generateStream(), s.ID,
			bson.NewObjectId().Hex(), dapp.ErrUnauthorizedAccess},
	}

	for _, tc := range cases {
		err := svc.UpdateStream(tc.owner, tc.streamId.Hex(), tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStream(t *testing.T) {
	svc := newService()
	s := generateStream()
	svc.AddStream(s)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		err      error
	}{
		{"view existing stream", s.ID, nil},
		{"view non-existing stream", bson.NewObjectId(), dapp.ErrNotFound},
	}

	for _, tc := range cases {
		_, err := svc.ViewStream(tc.streamId.Hex())
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRemoveStream(t *testing.T) {
	svc := newService()
	s := generateStream()
	svc.AddStream(s)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{"remove existing stream with wrong owner", s.ID,
			bson.NewObjectId().Hex(), dapp.ErrUnauthorizedAccess},
		{"remove existing stream", s.ID, s.Owner, nil},
		{"remove non-existing stream", s.ID, s.Owner, dapp.ErrNotFound},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.owner, tc.streamId.Hex())
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

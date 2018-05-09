package dapp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	"monetasa/dapp"
	"monetasa/dapp/mocks"
)

var (
	user string = bson.NewObjectId().Hex()

	stream dapp.Stream = dapp.Stream{
		User:        user,
		ID:          bson.NewObjectId(),
		Name:        "stream_name",
		Type:        "stream_type",
		Description: "stream_description",
		URL:         "www.stream_url.com",
		Price:       10,
		Location: dapp.Location{
			Type:        "Point",
			Coordinates: []float64{0, 0},
		},
	}
)

func newService() dapp.Service {
	streams := mocks.NewStreamRepository()
	return dapp.New(streams)
}

func TestAddStream(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc   string
		stream dapp.Stream
		user   string
		err    error
	}{
		{"add new stream", stream, user, nil},
		{"add existing stream", stream, user, dapp.ErrConflict},
	}

	for _, tc := range cases {
		_, err := svc.AddStream(tc.user, tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateStream(t *testing.T) {
	svc := newService()
	svc.AddStream(user, stream)

	cases := []struct {
		desc     string
		stream   dapp.Stream
		streamId bson.ObjectId
		user     string
		err      error
	}{
		{"update existing stream", stream, stream.ID, user, nil},
		{"update non-existing stream", stream, bson.NewObjectId(),
			user, dapp.ErrNotFound},
		{"update existing stream with wrong user", stream, stream.ID,
			bson.NewObjectId().Hex(), dapp.ErrUnauthorizedAccess},
	}

	for _, tc := range cases {
		err := svc.UpdateStream(tc.user, tc.streamId.Hex(), tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStream(t *testing.T) {
	svc := newService()
	svc.AddStream(user, stream)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		err      error
	}{
		{"view existing stream", stream.ID, nil},
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
	svc.AddStream(user, stream)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		user     string
		err      error
	}{
		{"remove existing stream with wrong user", stream.ID,
			bson.NewObjectId().Hex(), dapp.ErrUnauthorizedAccess},
		{"remove existing stream", stream.ID, user, nil},
		{"remove non-existing stream", stream.ID, user, dapp.ErrNotFound},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.user, tc.streamId.Hex())
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

package dapp_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	"monetasa/dapp"
	"monetasa/dapp/mocks"
)

var (
	owner string = bson.NewObjectId().Hex()

	stream dapp.Stream = dapp.Stream{
		Owner:       owner,
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

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func randomURL(len int) string {
	var buffer bytes.Buffer
	buffer.WriteString("https://")
	buffer.WriteString(randomString(len))
	buffer.WriteString(".com")
	return buffer.String()
}

func generateStream() dapp.Stream {
	return dapp.Stream{
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId(),
		Name:        randomString(8),
		Type:        randomString(4),
		Description: randomString(12),
		URL:         randomURL(6),
		Price:       rand.Intn(100),
		Location: dapp.Location{
			Type: "Point",
			Coordinates: []float64{
				rand.Float64() * (float64)(rand.Intn(360)-180),
				rand.Float64() * (float64)(rand.Intn(180)-90)},
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
	return dapp.New(streams)
}

func TestAddStream(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc   string
		stream dapp.Stream
		owner  string
		err    error
	}{
		{"add new stream", stream, owner, nil},
		{"add existing stream", stream, owner, dapp.ErrConflict},
	}

	for _, tc := range cases {
		_, err := svc.AddStream(tc.owner, tc.stream)
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
	svc.AddStream(owner, stream)

	cases := []struct {
		desc     string
		stream   dapp.Stream
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{"update existing stream", stream, stream.ID, owner, nil},
		{"update non-existing stream", stream, bson.NewObjectId(),
			owner, dapp.ErrNotFound},
		{"update existing stream with wrong owner", stream, stream.ID,
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
	svc.AddStream(owner, stream)

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
	svc.AddStream(owner, stream)

	cases := []struct {
		desc     string
		streamId bson.ObjectId
		owner    string
		err      error
	}{
		{"remove existing stream with wrong owner", stream.ID,
			bson.NewObjectId().Hex(), dapp.ErrUnauthorizedAccess},
		{"remove existing stream", stream.ID, owner, nil},
		{"remove non-existing stream", stream.ID, owner, dapp.ErrNotFound},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.owner, tc.streamId.Hex())
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

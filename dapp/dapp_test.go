package dapp_test

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"monetasa/dapp/mocks"
)

var (
	stream dapp.Stream = dapp.Stream{
		ID:          bson.NewObjectId(),
		Name:        "stream_name_01",
		Type:        "type_01",
		Description: "description_01",
		URL:         "www.url.com",
		Price:       10,
		Location: dapp.Location{
			Type:        "Point",
			Coordinates: []float64{0, 0},
		},
	}
)

func newService() dapp.Service {
	streams := mocks.NewStreamRepository()
	return manager.New(streams)
}

// func TestAddStream(t *testing.T) {
// 	svc := newService()

// 	_, err := svc.AddClient("", stream)
// }

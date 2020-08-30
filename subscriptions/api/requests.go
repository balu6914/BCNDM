package api

import (
	"github.com/datapace/datapace/subscriptions"

	"gopkg.in/mgo.v2/bson"
)

const (
	defLimit uint64 = 20
	maxLimit uint64 = 100
)

type apiReq interface {
	validate() error
}

type addSubsReq struct {
	UserID        string
	UserToken     string
	Subscriptions []subscriptions.Subscription
}

type searchSubsReq struct {
	owner string
	subscriptions.Query
}

func (req searchSubsReq) validate() error {
	if !bson.IsObjectIdHex(req.UserID) && !bson.IsObjectIdHex(req.StreamOwner) {
		return subscriptions.ErrMalformedEntity
	}

	return nil
}

type viewSubReq struct {
	userID         string
	subscriptionID string
}

func (req viewSubReq) validate() error {
	if !bson.IsObjectIdHex(req.userID) || !bson.IsObjectIdHex(req.subscriptionID) {
		return subscriptions.ErrMalformedEntity
	}
	return nil
}

type viewSubByUserAndStreamReq struct {
	userID   string
	streamID string
}

func (req viewSubByUserAndStreamReq) validate() error {
	if req.userID == "" || req.streamID == "" {
		return subscriptions.ErrMalformedEntity
	}

	return nil
}

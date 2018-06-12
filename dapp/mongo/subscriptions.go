package mongo

import (
	"monetasa/dapp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type subscriptionsRepository struct {
	db *mgo.Session
}

// Ensure a type subscriptionsRepository implements an interface SubscriptionsRepository
var _ dapp.SubscriptionsRepository = (*subscriptionsRepository)(nil)

func NewSubscriptionsRepository(db *mgo.Session) dapp.SubscriptionsRepository {
	return &subscriptionsRepository{db}
}

func (sr subscriptionsRepository) Create(sub dapp.Subscription) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionSubs)

	if err := c.Insert(sub); err != nil {
		if mgo.IsDup(err) {
			return dapp.ErrConflict
		}

		return err
	}

	return nil
}

func (sr subscriptionsRepository) Read(userID string) ([]dapp.Subscription, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionSubs)

	var subs []dapp.Subscription
	if err := c.Find(bson.M{"user_id": userID}).All(&subs); err != nil {
		if err == mgo.ErrNotFound {
			return subs, dapp.ErrNotFound
		}

		return subs, err
	}

	return subs, nil
}

func (sr subscriptionsRepository) Update(id string, sub dapp.Subscription) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionSubs)

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	query := bson.M{"_id": _id}
	update := bson.M{"$set": sub}
	if err := c.Update(query, update); err != nil {
		return err
	}

	return nil
}

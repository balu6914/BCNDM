package mongo

import (
	"monetasa/subscriptions"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Ensure a type subscriptionsRepository implements an interface SubscriptionsRepository
var _ subscriptions.SubscriptionsRepository = (*mongoRepository)(nil)

type mongoRepository struct {
	db *mgo.Session
}

func NewRepository(db *mgo.Session) subscriptions.SubscriptionsRepository {
	return &mongoRepository{db}
}

func (sr mongoRepository) Create(sub subscriptions.Subscription) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionName)

	if err := c.Insert(sub); err != nil {
		if mgo.IsDup(err) {
			return subscriptions.ErrConflict
		}

		return err
	}

	return nil
}

func (sr mongoRepository) Read(userID string) ([]subscriptions.Subscription, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionName)

	var subs []subscriptions.Subscription
	if err := c.Find(bson.M{"user_id": userID}).All(&subs); err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return []subscriptions.Subscription{}, nil
	}

	return subs, nil
}

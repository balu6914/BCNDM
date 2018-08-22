package mongo

import (
	"monetasa/subscriptions"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Ensure a type subscriptionsRepository implements an interface SubscriptionsRepository
var _ subscriptions.SubscriptionRepository = (*subscriptionRepository)(nil)

type subscriptionRepository struct {
	db *mgo.Session
}

// NewSubscriptionRepository returns new Subscription repository.
func NewSubscriptionRepository(db *mgo.Session) subscriptions.SubscriptionRepository {
	return &subscriptionRepository{db}
}

func (sr subscriptionRepository) Save(sub subscriptions.Subscription) (string, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionName)

	sub.ID = bson.NewObjectId()

	if err := c.Insert(sub); err != nil {
		if mgo.IsDup(err) {
			return "", subscriptions.ErrConflict
		}

		return "", err
	}

	return sub.ID.Hex(), nil
}

func (sr subscriptionRepository) Search(query subscriptions.Query) (subscriptions.Page, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionName)

	limit := int(query.Limit)
	page := int(query.Page)

	ret := subscriptions.Page{
		Page:    query.Page,
		Limit:   query.Limit,
		Content: []subscriptions.Subscription{},
	}

	results := []subscriptions.Subscription{}
	q := subscriptions.GenQuery(&query)

	total, err := c.Find(q).Count()
	if err != nil {
		return ret, err
	}

	ret.Total = uint64(total)
	start := limit * page
	if total < start {
		return ret, nil
	}

	if err = c.Find(q).Skip(start).Limit(limit).All(&results); err != nil {
		return ret, err
	}

	ret.Content = results
	return ret, nil
}

func (sr subscriptionRepository) One(id string) (subscriptions.Subscription, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collentionName)

	sub := subscriptions.Subscription{}

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	_id := bson.ObjectIdHex(id)
	if err := c.Find(bson.M{"_id": _id}).One(&sub); err != nil {
		if err == mgo.ErrNotFound {
			return sub, subscriptions.ErrNotFound
		}

		return sub, err
	}

	return sub, nil
}

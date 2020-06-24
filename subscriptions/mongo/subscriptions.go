package mongo

import (
	"time"

	"github.com/datapace/datapace/subscriptions"

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
	return &subscriptionRepository{db: db}
}

func (sr subscriptionRepository) Save(sub subscriptions.Subscription) (string, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	sub.ID = bson.NewObjectId()
	dbSub := toDBSub(sub)

	if err := c.Insert(dbSub); err != nil {
		if mgo.IsDup(err) {
			return "", subscriptions.ErrConflict
		}

		return "", err
	}

	return dbSub.ID.Hex(), nil
}

func (sr subscriptionRepository) Search(query subscriptions.Query) (subscriptions.Page, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	limit := int(query.Limit)
	page := int(query.Page)

	ret := subscriptions.Page{
		Page:    query.Page,
		Limit:   query.Limit,
		Content: []subscriptions.Subscription{},
	}

	results := []subscription{}
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

	for _, dbSub := range results {
		sub := toSub(dbSub)
		ret.Content = append(ret.Content, sub)
	}

	return ret, nil
}

func (sr subscriptionRepository) One(id string) (subscriptions.Subscription, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	var dbSub subscription
	query := bson.M{
		"_id":    bson.ObjectIdHex(id),
		"active": true,
	}
	if err := c.Find(query).One(&dbSub); err != nil {
		if err == mgo.ErrNotFound {
			return subscriptions.Subscription{}, subscriptions.ErrNotFound
		}

		return subscriptions.Subscription{}, err
	}

	sub := toSub(dbSub)

	return sub, nil
}

func (sr subscriptionRepository) OneByUserAndStream(userID string, streamID string) (subscriptions.Subscription, error) {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	var dbSub subscription
	query := bson.M{
		"user_id":   userID,
		"stream_id": streamID,
		"end_date":  map[string]time.Time{"$gt": time.Now()},
		"active":    true,
	}

	if err := c.Find(query).One(&dbSub); err != nil {
		if err == mgo.ErrNotFound {
			return subscriptions.Subscription{}, subscriptions.ErrNotFound
		}

		return subscriptions.Subscription{}, err
	}

	sub := toSub(dbSub)
	return sub, nil
}

func (sr subscriptionRepository) Activate(id string) error {
	return sr.setActive(id, true)
}

func (sr subscriptionRepository) Remove(id string) error {
	return sr.setActive(id, false)
}

func (sr subscriptionRepository) setActive(id string, active bool) error {
	s := sr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	update := bson.M{
		"$set": subscription{Active: active},
	}
	if err := c.UpdateId(bson.ObjectIdHex(id), update); err != nil {
		if err == mgo.ErrNotFound {
			return subscriptions.ErrNotFound
		}
		if mgo.IsDup(err) {
			return subscriptions.ErrConflict
		}

		return err
	}

	return nil
}

// Subscription is subscription representation in DB.
type subscription struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	UserID      string        `bson:"user_id,omitempty"`
	StreamID    string        `bson:"stream_id,omitempty"`
	StreamOwner string        `bson:"stream_owner,omitempty"`
	Hours       uint64        `bson:"hours,omitempty"`
	StartDate   time.Time     `bson:"start_date,omitempty"`
	EndDate     time.Time     `bson:"end_date,omitempty"`
	StreamURL   string        `bson:"stream_url,omitempty"`
	StreamPrice uint64        `bson:"stream_price,omitempty"`
	StreamName  string        `bson:"stream_name,omitempty"`
	Active      bool          `bson:"active"`
}

func toDBSub(sub subscriptions.Subscription) subscription {
	return subscription{
		ID:          sub.ID,
		UserID:      sub.UserID,
		StreamID:    sub.StreamID,
		StreamOwner: sub.StreamOwner,
		Hours:       sub.Hours,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		StreamURL:   sub.StreamURL,
		StreamPrice: sub.StreamPrice,
		StreamName:  sub.StreamName,
		Active:      false,
	}
}

func toSub(dbSub subscription) subscriptions.Subscription {
	return subscriptions.Subscription{
		ID:          dbSub.ID,
		UserID:      dbSub.UserID,
		StreamID:    dbSub.StreamID,
		StreamOwner: dbSub.StreamOwner,
		Hours:       dbSub.Hours,
		StartDate:   dbSub.StartDate,
		EndDate:     dbSub.EndDate,
		StreamURL:   dbSub.StreamURL,
		StreamPrice: dbSub.StreamPrice,
		StreamName:  dbSub.StreamName,
	}
}

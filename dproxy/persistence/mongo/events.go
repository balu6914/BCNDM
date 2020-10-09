package mongo

import (
	"github.com/datapace/datapace/dproxy/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type eventRepository struct {
	db *mgo.Session
}

var _ persistence.EventRepository = (*eventRepository)(nil)

func NewEventsRepository(db *mgo.Session) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

// Accumulate updates event count for event initiator and returns total count for specific initiator
func (er *eventRepository) Accumulate(event persistence.Event) (int, error) {
	s := er.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)
	query := bson.M{"initiator": event.Initiator}
	d := bson.M{
		"$set": bson.M{
			"last_request": time.Now(),
		},
		"$inc": bson.M{
			"count": 1,
		},
	}
	_, err := c.Upsert(query, d)
	if err != nil {
		return 0, err
	}
	var dbEvent dbEvent
	if err = c.Find(query).One(&dbEvent); err != nil {
		return 0, err
	}
	return dbEvent.Count, nil
}

type dbEvent struct {
	Count int `bson:"count"`
}

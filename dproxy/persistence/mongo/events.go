package mongo

import (
	"time"

	"github.com/datapace/datapace/dproxy/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
func (er *eventRepository) Accumulate(event persistence.Event, unit string) (int, error) {
	s := er.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)
	dbEvent, _ := toDBEvent(event)

	if err := c.Insert(dbEvent); err != nil {
		return 0, err
	}
	var period time.Duration
	switch unit {
	case "second":
		period = time.Second
	case "minute":
		period = time.Minute
	case "hour":
		period = time.Hour
	case "day":
		period = 24 * time.Hour
	case "month":
		period = 30 * 24 * time.Hour
	case "year":
		period = 365 * 24 * time.Hour
	default:
		period = 0 //nije nula vec beskonacno, sredi ovo
	}
	now := time.Now()

	filter := bson.M{"initiator": event.Initiator}

	if period != 0 {
		since := now.Add(-period)
		filter["request_time"] = bson.M{"$gte": since, "$lt": now}
	}

	// Count the number of documents that match the filter
	cnt, err := c.Find(filter).Count()
	if err != nil {
		return 0, err
	}
	return cnt, nil

}

type dbEvent struct {
	Initiator   string     `bson:"initiator,omitempty"`
	RequestTime *time.Time `bson:"request_time,omitempty"`
}

func toDBEvent(event persistence.Event) (dbEvent, error) {
	dbe := dbEvent{
		Initiator:   event.Initiator,
		RequestTime: &event.Time,
	}
	return dbe, nil
}

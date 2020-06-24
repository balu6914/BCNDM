package postgres

import (
	"github.com/datapace/datapace/dproxy/persistence"
	"github.com/jmoiron/sqlx"
)

type eventRepository struct {
	db *sqlx.DB
}

var _ persistence.EventRepository = (*eventRepository)(nil)

func NewEventsRepository(db *sqlx.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

// Accumulate saves new event and returns total number of events for specific initiator
func (er *eventRepository) Accumulate(event persistence.Event) (int, error) {
	tx, err := er.db.Begin()
	if err != nil {
		return 0, err
	}
	if _, err := er.db.Exec("INSERT INTO events(init_time, initiator) VALUES($1,$2)", event.Time, event.Initiator); err != nil {
		return 0, err
	}
	var cnt int
	if err := er.db.Get(&cnt, "SELECT COUNT(*) FROM events WHERE initiator = $1", event.Initiator); err != nil {
		return 0, err
	}
	tx.Commit()
	return cnt, nil
}

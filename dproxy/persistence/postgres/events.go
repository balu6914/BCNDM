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
func (er *eventRepository) Accumulate(event persistence.Event, unit string) (int, error) { //dodaj u upite unit uslov
	tx, err := er.db.Begin()
	if err != nil {
		return 0, err
	}
	if _, err := er.db.Exec("INSERT INTO events(init_time, initiator) VALUES($1,$2)", event.Time, event.Initiator); err != nil {
		return 0, err
	}

	var period string

	switch unit {
	case "second":
		period = "AND init_time >= NOW() - CAST ('1 second' AS INTERVAL)"
	case "minute":
		period = "AND init_time >= NOW() - CAST ('1 minute' AS INTERVAL)"
	case "hour":
		period = "AND init_time >= NOW() - CAST ('1 hour' AS INTERVAL)"
	case "day":
		period = "AND init_time >= NOW() - CAST ('1 day' AS INTERVAL)"
	case "month":
		period = "AND init_time >= NOW() - CAST ('30 days' AS INTERVAL)"
	case "year":
		period = "AND init_time >= NOW() - CAST ('365 days' AS INTERVAL)"
	default:
		period = ""
	}

	var cnt int
	if err := er.db.Get(&cnt, "SELECT COUNT(*) FROM events WHERE initiator = $1 $2", event.Initiator, period); err != nil {
		return 0, err
	}
	tx.Commit()
	return cnt, nil
}

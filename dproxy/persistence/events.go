package persistence

import (
	"time"
)

type EventRepository interface {
	Accumulate(event Event, unit string) (int, error)
	List(Query) ([]Event, error)
}

type Event struct {
	Time      time.Time
	Initiator string
	SubID     string
}

type Query struct {
	Limit uint32

	Cursor Event

	Sort Sort
}

type Sort struct {
	Order SortOrder
	By    SortBy
}

type SortOrder int

const (
	SortOrderAsc = iota
	SortOrderDesc
)

func (so SortOrder) String() string {
	return [...]string{
		"Asc",
		"Desc",
	}[so]
}

type SortBy int

const (
	SortByDate = iota
)

func (sb SortBy) String() string {
	return [...]string{
		"Date",
	}[sb]
}

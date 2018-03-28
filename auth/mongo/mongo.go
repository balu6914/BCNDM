package mongo

import (
	"gopkg.in/mgo.v2"
)

// DBName - DB name
// CollectionName - Collection name
const (
	DBName         string = "monetasa"
	CollectionName string = "auth"
)

// Repository - get Mongo session
type Repository struct {
	Session *mgo.Session
}

// NewRepository - create new Mongo repository
func NewRepository(ms *mgo.Session) *Repository {
	return &Repository{Session: ms}
}

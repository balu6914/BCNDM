package mongo

import (
	"fmt"

	"gitlab.com/blocksense/monetasa/auth"
	"gopkg.in/mgo.v2"
)

// DBName - DB name
// CollectionName - Collection name
const (
	DBName         string = "monetasa"
	CollectionName string = "auth"
)

// Connect creates a connection to the MongoDB instance. A non-nil error
// is returned to indicate failure.
func Connect(cfg *config) (*mgo.Session, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{cfg.MongoURL + ":" + strconv.Itoa(cfg.MongoPort)},
		Timeout:  time.Duration(cfg.MongoConnectTimeout) * time.Millisecond,
		Database: cfg.MongoDatabase,
		Username: cfg.MongoUser,
		Password: cfg.MongoPass,
	}

	ms, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	ms.SetSocketTimeout(time.Duration(cfg.MongoSocketTimeout) * time.Millisecond)
	ms.SetMode(mgo.Monotonic, true)

	return ms, nil
}

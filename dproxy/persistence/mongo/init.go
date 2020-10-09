package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// DBName - DB name
// CollectionName - Collection name
const (
	dbName         = "datapace-dproxy"
	collectionName = "events"
)

// Connect creates a connection to the MongoDB instance. A non-nil error
// is returned to indicate failure.
func Connect(addr string, tout int, socketTout int, db string,
	user string, pass string) (*mgo.Session, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{addr},
		Timeout:  time.Duration(tout) * time.Millisecond,
		Database: db,
		Username: user,
		Password: pass,
	}
	ms, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	ms.SetSocketTimeout(time.Duration(socketTout) * time.Millisecond)
	ms.SetMode(mgo.Monotonic, true)

	session := ms.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	indices := []mgo.Index{
		mgo.Index{
			Name: "events",
			Key:  []string{"initiator"},
		},
	}
	for _, idx := range indices {
		collection.EnsureIndex(idx)
	}
	return ms, nil
}

package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// dbName - If name is empty, the database name provided in DialInfo is used instead
const (
	dbName     = ""
	collection = "access-requests"
)

var (
	indices = []mgo.Index{
		{
			Key: []string{
				"sender",
				"receiver",
			},
			Unique:     true,
			DropDups:   false,
			Background: false,
			Sparse:     true,
		},
	}
)

// Connect creates a connection to the MongoDB instance. A non-nil error
// is returned to indicate failure.
func Connect(addr string, tout int, socketTout int, db string, user string, pass string) (*mgo.Session, error) {
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

	return ms, nil
}

func EnsureIndices(ms *mgo.Session) error {
	mdb := ms.DB(dbName)
	coll := mdb.C(collection)
	for _, index := range indices {
		if err := coll.EnsureIndex(index); err != nil {
			return err
		}
	}
	return nil
}

package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

const (
	dbName             = "datapace-auth"
	usersCollection    = "users"
	policiesCollection = "policies"
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

	// Create unique constraint in mongoDB.
	session := ms.Copy()
	defer session.Close()
	usersColl := session.DB(dbName).C(usersCollection)

	usersIdx := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     true,
	}
	usersColl.EnsureIndex(usersIdx)

	policiesColl := session.DB(dbName).C(policiesCollection)

	// Use name and owner fields as a unique constraint.
	// This means that one user can't create many policies
	// with the same name.
	policiesIdx := mgo.Index{
		Key:        []string{"name", "owner"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     true,
	}
	policiesColl.EnsureIndex(policiesIdx)

	return ms, nil
}

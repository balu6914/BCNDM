package mongo

import (
	"fmt"
	access "github.com/datapace/datapace/access-control"
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ access.RequestRepository = (*accessRequestRepository)(nil)

type accessRequestRepository struct {
	db *mgo.Session
}

const (
	resultsPageLimit = 100
)

// NewAccessRequestRepository instantiates a Mongo implementation of access
// request repository.
func NewAccessRequestRepository(db *mgo.Session) access.RequestRepository {
	return &accessRequestRepository{db}
}

func (repo accessRequestRepository) RequestAccess(sender, receiver string) (string, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)
	q := bson.M{"sender": sender, "receiver": receiver}
	var ar accessRequest
	if err := collection.Find(q).One(&ar); err != nil {
		if err == mgo.ErrNotFound {
			return repo.insertNew(collection, sender, receiver)
		}
		return "", err
	}
	// otherwise - found, update the state
	u := bson.M{"$set": bson.M{"state": access.Pending}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return "", access.ErrNotFound
		}
		return "", err
	}
	return ar.ID.Hex(), nil
}

func (repo accessRequestRepository) insertNew(collection *mgo.Collection, sender, receiver string) (string, error) {
	ar := accessRequest{
		ID:       bson.NewObjectId(),
		Sender:   sender,
		Receiver: receiver,
		State:    access.Pending,
	}
	if err := collection.Insert(ar); err != nil {
		if mgo.IsDup(err) {
			return "", access.ErrConflict
		}

		return "", err
	}
	return ar.ID.Hex(), nil
}

func (repo accessRequestRepository) One(id string) (access.Request, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

	if !bson.IsObjectIdHex(id) {
		return access.Request{}, access.ErrMalformedEntity
	}

	var ar accessRequest
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&ar); err != nil {
		if err == mgo.ErrNotFound {
			return access.Request{}, access.ErrNotFound
		}
		return access.Request{}, err
	}

	req := access.Request{
		ID:       ar.ID.Hex(),
		Sender:   ar.Sender,
		Receiver: ar.Receiver,
		State:    ar.State,
	}

	return req, nil
}

func (repo accessRequestRepository) ListSent(sender string, state access.State) ([]access.Request, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

	q := bson.M{"sender": sender}
	if state != access.State("") {
		q["state"] = state
	}

	requests := []accessRequest{}
	if err := collection.Find(q).All(&requests); err != nil {
		return []access.Request{}, err
	}

	list := []access.Request{}
	for _, req := range requests {
		list = append(list, access.Request{
			ID:       req.ID.Hex(),
			Sender:   req.Sender,
			Receiver: req.Receiver,
			State:    req.State,
		})
	}

	return list, nil
}

func (repo accessRequestRepository) ListReceived(receiver string, state access.State) ([]access.Request, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

	q := bson.M{"receiver": receiver}
	if state != access.State("") {
		q["state"] = state
	}

	requests := []accessRequest{}
	if err := collection.Find(q).All(&requests); err != nil {
		return []access.Request{}, err
	}

	list := []access.Request{}
	for _, req := range requests {
		list = append(list, access.Request{
			ID:       req.ID.Hex(),
			Sender:   req.Sender,
			Receiver: req.Receiver,
			State:    req.State,
		})
	}

	return list, nil
}

func (repo accessRequestRepository) Approve(receiver, id string) error {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

	if !bson.IsObjectIdHex(id) {
		return access.ErrMalformedEntity
	}

	q := bson.M{"_id": bson.ObjectIdHex(id), "receiver": receiver, "state": access.Pending}
	u := bson.M{"$set": bson.M{"state": access.Approved}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return access.ErrNotFound
		}
		return err
	}

	return nil
}

func (repo accessRequestRepository) Revoke(receiver, id string) error {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

	if !bson.IsObjectIdHex(id) {
		return access.ErrMalformedEntity
	}

	q := bson.M{"_id": bson.ObjectIdHex(id), "receiver": receiver}
	u := bson.M{"$set": bson.M{"state": access.Revoked}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return fmt.Errorf("revoke failed, access request query: %s, cause: %w", q, access.ErrNotFound)
		}
		return fmt.Errorf("revoke failed, access request query: %s, cause: %w", q, err)
	}

	return nil
}

func (repo accessRequestRepository) GrantAccess(srcUserId string, dstUserId string) (string, error) {
	reqId, err := repo.RequestAccess(dstUserId, srcUserId)
	if err != nil {
		return "", err
	}
	return reqId, repo.Approve(srcUserId, reqId)
}

func (repo accessRequestRepository) MigrateSameSenderAndReceiver() (err error) {
	sess := repo.db.Copy()
	defer sess.Close()
	coll := sess.DB(dbName).C(collection)
	countDone := 0
	for {
		var acs []accessRequest
		err = coll.
			Find(bson.M{}).
			Limit(resultsPageLimit).
			Skip(countDone).
			All(&acs)
		if err != nil {
			return
		}
		if len(acs) == 0 {
			break // no more records
		}
		countDone += len(acs)
		for _, ac := range acs {
			err = repo.migrateToSingleRecord(coll, ac.Sender, ac.Receiver)
			if err != nil {
				return
			}
		}
	}
	return
}

func (repo accessRequestRepository) migrateToSingleRecord(coll *mgo.Collection, sender, receiver string) (err error) {
	q := bson.M{
		"sender":   sender,
		"receiver": receiver,
	}
	var dupAcs []accessRequest
	err = coll.
		Find(q).
		All(&dupAcs)
	if err != nil {
		return
	}
	if len(dupAcs) > 1 {
		// sort the duplicate access request by object id time in the ascending order (latest is last)
		sort.Slice(dupAcs, func(i, j int) bool {
			aci := dupAcs[i]
			acj := dupAcs[j]
			return aci.ID.Time().Before(acj.ID.Time())
		})
		latestId := dupAcs[len(dupAcs)-1].ID
		// remove all except the latest
		q["_id"] = bson.M{
			"$ne": latestId,
		}
		_, err = coll.RemoveAll(q)
	}
	return
}

type accessRequest struct {
	ID       bson.ObjectId `bson:"_id"`
	Sender   string        `bson:"sender"`
	Receiver string        `bson:"receiver"`
	State    access.State  `bson:"state"`
}

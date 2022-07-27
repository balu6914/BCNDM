package mongo

import (
	access "github.com/datapace/datapace/access-control"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ access.RequestRepository = (*accessRequestRepository)(nil)

type accessRequestRepository struct {
	db *mgo.Session
}

// NewAccessRequestRepository instantiates a Mongo implementation of access
// request repository.
func NewAccessRequestRepository(db *mgo.Session) access.RequestRepository {
	return &accessRequestRepository{db}
}

func (repo accessRequestRepository) RequestAccess(sender, receiver string) (string, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collection)

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

	q := bson.M{"_id": bson.ObjectIdHex(id), "receiver": receiver, "state": access.Pending}
	u := bson.M{"$set": bson.M{"state": access.Revoked}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return access.ErrNotFound
		}
		return err
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

type accessRequest struct {
	ID       bson.ObjectId `bson:"_id"`
	Sender   string        `bson:"sender"`
	Receiver string        `bson:"receiver"`
	State    access.State  `bson:"state"`
}

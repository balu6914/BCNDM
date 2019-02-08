package mongo

import (
	"datapace/auth"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ auth.AccessRequestRepository = (*accessRequestRepository)(nil)

type accessRequestRepository struct {
	db *mgo.Session
}

// NewAccessRequestRepository instantiates a Mongo implementation of access
// request repository.
func NewAccessRequestRepository(db *mgo.Session) auth.AccessRequestRepository {
	return &accessRequestRepository{db}
}

func (repo accessRequestRepository) RequestAccess(sender, receiver string) (string, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(accessRequestsCollection)

	ar := accessRequest{
		ID:       bson.NewObjectId(),
		Sender:   sender,
		Receiver: receiver,
		State:    auth.Pending,
	}
	if err := collection.Insert(ar); err != nil {
		if mgo.IsDup(err) {
			return "", auth.ErrConflict
		}

		return "", err
	}

	return ar.ID.Hex(), nil
}

func (repo accessRequestRepository) ListSentAccessRequests(sender string, state auth.State) ([]auth.AccessRequest, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(accessRequestsCollection)

	q := bson.M{"sender": sender, "state": state}
	requests := []accessRequest{}
	if err := collection.Find(q).All(&requests); err != nil {
		return []auth.AccessRequest{}, err
	}

	list := []auth.AccessRequest{}
	for _, req := range requests {
		list = append(list, auth.AccessRequest{
			ID:       req.ID.Hex(),
			Sender:   req.Sender,
			Receiver: req.Receiver,
			State:    req.State,
		})
	}

	return list, nil
}

func (repo accessRequestRepository) ListReceivedAccessRequests(receiver string, state auth.State) ([]auth.AccessRequest, error) {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(accessRequestsCollection)

	q := bson.M{"receiver": receiver, "state": state}
	requests := []accessRequest{}
	if err := collection.Find(q).All(&requests); err != nil {
		return []auth.AccessRequest{}, err
	}

	list := []auth.AccessRequest{}
	for _, req := range requests {
		list = append(list, auth.AccessRequest{
			ID:       req.ID.Hex(),
			Sender:   req.Sender,
			Receiver: req.Receiver,
			State:    req.State,
		})
	}

	return list, nil
}

func (repo accessRequestRepository) ApproveAccessRequest(receiver, id string) error {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(accessRequestsCollection)

	q := bson.M{"_id": id, "receiver": receiver, "state": auth.Pending}
	u := bson.M{"$set": bson.M{"state": auth.Approved}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	return nil
}

func (repo accessRequestRepository) RejectAccessRequest(receiver, id string) error {
	session := repo.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(accessRequestsCollection)

	q := bson.M{"_id": id, "receiver": receiver, "state": auth.Pending}
	u := bson.M{"$set": bson.M{"state": auth.Revoked}}
	if err := collection.Update(q, u); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	return nil
}

type accessRequest struct {
	ID       bson.ObjectId `bson:"_id"`
	Sender   string        `bson:"sender"`
	Receiver string        `bson:"receiver"`
	State    auth.State    `bson:"state"`
}

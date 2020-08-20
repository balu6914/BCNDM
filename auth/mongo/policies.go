package mongo

import (
	"github.com/datapace/datapace/auth"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ auth.PolicyRepository = (*policyRepository)(nil)

type policyRepository struct {
	db *mgo.Session
}

// NewPolicyRepository instantiates a Mongo implementation of user
// repository.
func NewPolicyRepository(db *mgo.Session) auth.PolicyRepository {
	return &policyRepository{
		db: db,
	}
}

func (pr *policyRepository) Save(policy auth.Policy) (string, error) {
	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := toMongoPolicy(policy)
	mp.ID = bson.NewObjectId()

	if err := collection.Insert(mp); err != nil {
		if mgo.IsDup(err) {
			return "", auth.ErrConflict
		}
		return "", err
	}

	return mp.ID.Hex(), nil
}

func (pr *policyRepository) OneByID(id string) (auth.Policy, error) {
	if !bson.IsObjectIdHex(id) {
		return auth.Policy{}, auth.ErrMalformedEntity
	}

	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := mongoPolicy{}

	oid := bson.ObjectIdHex(id)
	if err := collection.Find(bson.M{"_id": oid}).One(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return auth.Policy{}, auth.ErrNotFound
		}
		return auth.Policy{}, err
	}

	return toPolicy(mp), nil
}

func (pr *policyRepository) OneByName(owner, name string) (auth.Policy, error) {
	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := mongoPolicy{}
	if err := collection.Find(
		bson.M{"$and": []bson.M{{"name": name, "owner": owner}}}).One(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return auth.Policy{}, auth.ErrNotFound
		}
		return auth.Policy{}, err
	}

	return toPolicy(mp), nil
}

func (pr *policyRepository) List(owner string) ([]auth.Policy, error) {
	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := []mongoPolicy{}
	q := bson.M{
		"owner": owner,
	}

	if err := collection.Find(q).All(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return []auth.Policy{}, auth.ErrNotFound
		}
		return []auth.Policy{}, err
	}

	ret := []auth.Policy{}
	for _, p := range mp {
		ret = append(ret, toPolicy(p))
	}

	return ret, nil
}

func (pr *policyRepository) ListByIDs(ids []string) ([]auth.Policy, error) {
	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := []mongoPolicy{}
	mongoIds := []bson.ObjectId{}
	for _, id := range ids {
		if !bson.IsObjectIdHex(id) {
			return nil, auth.ErrMalformedEntity
		}
		oid := bson.ObjectIdHex(id)
		mongoIds = append(mongoIds, oid)
	}
	q := bson.M{
		"_id": bson.M{
			"$in": mongoIds,
		},
	}
	if err := collection.Find(q).All(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return []auth.Policy{}, auth.ErrNotFound
		}
		return []auth.Policy{}, err
	}

	ret := []auth.Policy{}
	for _, p := range mp {
		ret = append(ret, toPolicy(p))
	}

	return ret, nil
}

func (pr *policyRepository) Remove(id string) error {
	if !bson.IsObjectIdHex(id) {
		return auth.ErrMalformedEntity
	}

	session := pr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)
	oid := bson.ObjectIdHex(id)

	if err := collection.RemoveId(oid); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	uc := session.DB(dbName).C(usersCollection)
	return uc.Update(bson.M{}, bson.M{
		"$pull": bson.M{
			"policies": id,
		},
	})
}

func (pr *policyRepository) Attach(policyID, userID string) error {
	if !bson.IsObjectIdHex(policyID) || !bson.IsObjectIdHex(userID) {
		return auth.ErrMalformedEntity
	}
	session := pr.db.Copy()
	defer session.Close()

	users := session.DB(dbName).C(usersCollection)
	uid := bson.ObjectIdHex(userID)
	var mu mongoUser
	if err := users.Find(bson.M{"_id": uid}).One(&mu); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	policies := session.DB(dbName).C(policiesCollection)
	pid := bson.ObjectIdHex(policyID)
	var mp mongoPolicy
	if err := policies.Find(bson.M{"_id": pid}).One(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	mu.Policies = append(mu.Policies, policyID)
	update := bson.M{
		"$set": bson.M{
			"policies": mu.Policies,
		},
	}

	if err := users.UpdateId(mu.ID, update); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}
		return err
	}

	return nil
}

func (pr *policyRepository) Detach(policyID, userID string) error {
	if !bson.IsObjectIdHex(policyID) || !bson.IsObjectIdHex(userID) {
		return auth.ErrMalformedEntity
	}
	session := pr.db.Copy()
	defer session.Close()

	collection := session.DB(dbName).C(usersCollection)
	uid := bson.ObjectIdHex(userID)
	var mu mongoUser
	if err := collection.Find(bson.M{"_id": uid}).One(&mu); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}
	idx := -1
	for i, id := range mu.Policies {
		if id == policyID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil
	}
	mu.Policies = append(mu.Policies[:idx], mu.Policies[idx+1:len(mu.Policies)]...)

	update := bson.M{
		"$set": bson.M{
			"policies": mu.Policies,
		},
	}

	if err := collection.UpdateId(mu.ID, update); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}
		return err
	}

	return nil
}

type mongoPolicy struct {
	ID          bson.ObjectId     `bson:"_id,omitempty"`
	Version     string            `bson:"version,omitempty"`
	Name        string            `bson:"name,omitempty"`
	Owner       string            `bson:"owner,omitempty"`
	Rules       []mongoRule       `bson:"rules,omitempty"`
	Constraints []auth.Constraint `bson:"constraints,omitempty"`
}

type mongoRule struct {
	Action    auth.Action    `bson:"action,omitempty"`
	Type      string         `bson:"type,omitempty"`
	Condition mongoCondition `bson:"condition,omitempty"`
}

type mongoCondition struct {
	Key   string `bson:"key,omitempty"`
	Value string `bson:"value,omitempty"`
}

func toMongoPolicy(policy auth.Policy) mongoPolicy {
	id := bson.NewObjectId()
	mp := mongoPolicy{
		ID:          id,
		Version:     policy.Version,
		Owner:       policy.Owner,
		Name:        policy.Name,
		Rules:       []mongoRule{},
		Constraints: policy.Constraints,
	}
	for _, rule := range policy.Rules {
		mr := mongoRule{
			Action: rule.Action,
			Type:   rule.Type,
		}
		if con, ok := rule.Condition.(auth.SimpleCondition); ok {
			mr.Condition = mongoCondition{
				Key:   con.Key,
				Value: con.Value,
			}
		}
		mp.Rules = append(mp.Rules, mr)
	}
	return mp
}

func toPolicy(mp mongoPolicy) auth.Policy {
	p := auth.Policy{
		ID:          mp.ID.Hex(),
		Version:     mp.Version,
		Name:        mp.Name,
		Rules:       []auth.Rule{},
		Constraints: mp.Constraints,
	}

	for _, mr := range mp.Rules {
		rule := auth.Rule{
			Action: mr.Action,
			Type:   mr.Type,
		}
		if mr.Condition.Key != "" {
			rule.Condition = auth.SimpleCondition{
				Key:   mr.Condition.Key,
				Value: mr.Condition.Value,
			}
		}
		p.Rules = append(p.Rules, rule)
	}
	return p
}

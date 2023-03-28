package mongo

import (
	"github.com/datapace/datapace/executions"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ executions.ExecutionRepository = (*executionRepository)(nil)

type executionRepository struct {
	db *mgo.Session
}

// NewExecutionRepository returns new execution repository instance.
func NewExecutionRepository(db *mgo.Session) executions.ExecutionRepository {
	return executionRepository{
		db: db,
	}
}

func (er executionRepository) Create(exec executions.Execution) (string, error) {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(execCollection)
	id := bson.NewObjectId()

	e := execution{
		ID:         id,
		Name:       exec.Name,
		ExternalID: exec.ExternalID,
		Owner:      exec.Owner,
		Algo:       exec.Algo,
		Data:       exec.Data,
		Metadata:   exec.Metadata,
		State:      exec.State,
	}
	if err := c.Insert(e); err != nil {
		if mgo.IsDup(err) {
			return "", executions.ErrConflict
		}
		return "", err
	}

	return id.Hex(), nil
}

func (er executionRepository) Update(exec executions.Execution) error {
	s := er.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(exec.ID) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(execCollection)

	e := execution{
		ID:         bson.ObjectIdHex(exec.ID),
		Name:       exec.Name,
		ExternalID: exec.ExternalID,
		Owner:      exec.Owner,
		Algo:       exec.Algo,
		Data:       exec.Data,
		Metadata:   exec.Metadata,
		State:      exec.State,
	}
	if err := c.UpdateId(e.ID, e); err != nil {
		return err
	}

	return nil
}

func (er executionRepository) UpdateState(externalID string, state executions.State) error {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(execCollection)
	u := bson.M{"$set": bson.M{"state": state}}
	if err := c.Update(bson.M{"external_id": externalID}, u); err != nil {
		if err == mgo.ErrNotFound {
			return executions.ErrNotFound
		}
		return err
	}

	return nil
}

func (er executionRepository) Execution(owner, id string) (executions.Execution, error) {
	s := er.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return executions.Execution{}, executions.ErrMalformedData
	}

	c := s.DB(dbName).C(execCollection)
	q := bson.M{"_id": bson.ObjectIdHex(id), "owner": owner}

	var exec execution
	if err := c.Find(q).One(&exec); err != nil {
		if err == mgo.ErrNotFound {
			return executions.Execution{}, executions.ErrNotFound
		}
	}

	return executions.Execution{
		ID:       id,
		Owner:    owner,
		Algo:     exec.Algo,
		Data:     exec.Data,
		Metadata: exec.Metadata,
		State:    exec.State,
	}, nil
}

func (er executionRepository) List(owner string) ([]executions.Execution, error) {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(execCollection)
	q := bson.M{"owner": owner}

	var execs []execution
	if err := c.Find(q).All(&execs); err != nil {
		return []executions.Execution{}, err
	}

	es := []executions.Execution{}
	for _, exec := range execs {
		es = append(es, executions.Execution{
			ID:       exec.ID.Hex(),
			Owner:    owner,
			Algo:     exec.Algo,
			Data:     exec.Data,
			Metadata: exec.Metadata,
			State:    exec.State,
		})
	}

	return es, nil
}

type execution struct {
	ID         bson.ObjectId          `bson:"_id"`
	Name       string                 `bson:"name"`
	ExternalID string                 `bson:"external_id"`
	Owner      string                 `bson:"owner"`
	Algo       string                 `bson:"algo"`
	Data       string                 `bson:"data"`
	Metadata   map[string]interface{} `bson:"metadata"`
	State      executions.State       `bson:"state"`
}

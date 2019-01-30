package mongo

import (
	"datapace/executions"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ executions.ExecutionRepository = (*executionRepository)(nil)

type executionRepository struct {
	db *mgo.Session
}

// New returns new execution repository instance.
func New(db *mgo.Session) executions.ExecutionRepository {
	return executionRepository{
		db: db,
	}
}

func (er executionRepository) Create(owner, algo, data string, mode executions.JobMode) (string, error) {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(collection)
	id := bson.NewObjectId()

	exec := execution{
		ID:    id,
		Owner: owner,
		State: executions.Executing,
		Algo:  algo,
		Data:  data,
		Mode:  mode,
	}
	if err := c.Insert(exec); err != nil {
		if mgo.IsDup(err) {
			return "", executions.ErrConflict
		}
		return "", err
	}

	return id.Hex(), nil
}

func (er executionRepository) Finish(id string) error {
	s := er.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(collection)
	u := bson.M{"$set": bson.M{"state": executions.Done}}
	if err := c.UpdateId(bson.ObjectIdHex(id), u); err != nil {
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

	c := s.DB(dbName).C(collection)
	q := bson.M{"_id": bson.ObjectIdHex(id), "owner": owner}

	var exec execution
	if err := c.Find(q).One(&exec); err != nil {
		if err == mgo.ErrNotFound {
			return executions.Execution{}, executions.ErrNotFound
		}
	}

	return executions.Execution{
		ID:    id,
		Owner: owner,
		State: exec.State,
		Algo:  exec.Algo,
		Data:  exec.Data,
		Mode:  exec.Mode,
	}, nil
}

func (er executionRepository) List(owner string) ([]executions.Execution, error) {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(collection)
	q := bson.M{"owner": owner}

	var execs []execution
	if err := c.Find(q).All(&execs); err != nil {
		return []executions.Execution{}, err
	}

	es := []executions.Execution{}
	for _, exec := range execs {
		es = append(es, executions.Execution{
			ID:    exec.ID.Hex(),
			Owner: owner,
			State: exec.State,
			Algo:  exec.Algo,
			Data:  exec.Data,
			Mode:  exec.Mode,
		})
	}

	return es, nil
}

type execution struct {
	ID    bson.ObjectId      `bson:"_id"`
	Owner string             `bson:"owner"`
	State executions.State   `bson:"state"`
	Algo  string             `bson:"algo"`
	Data  string             `bson:"data"`
	Mode  executions.JobMode `bson:"mode"`
}

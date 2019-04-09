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
		ID:                       id,
		Owner:                    exec.Owner,
		Algo:                     exec.Algo,
		Data:                     exec.Data,
		AdditionalLocalJobArgs:   exec.AdditionalLocalJobArgs,
		Type:                     exec.Type,
		GlobalTimeout:            exec.GlobalTimeout,
		LocalTimeout:             exec.LocalTimeout,
		AdditionalPreprocessArgs: exec.AdditionalPreprocessArgs,
		Mode:                     exec.Mode,
		AdditionalGlobalJobArgs:  exec.AdditionalGlobalJobArgs,
		AdditionalFiles:          exec.AdditionalFiles,
		Token:                    exec.Token,
		State:                    exec.State,
	}
	if err := c.Insert(e); err != nil {
		if mgo.IsDup(err) {
			return "", executions.ErrConflict
		}
		return "", err
	}

	return id.Hex(), nil
}

func (er executionRepository) UpdateToken(id, token string) error {
	s := er.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(execCollection)
	u := bson.M{"$set": bson.M{"token": token}}
	if err := c.UpdateId(bson.ObjectIdHex(id), u); err != nil {
		if err == mgo.ErrNotFound {
			return executions.ErrNotFound
		}
		return err
	}

	return nil
}

func (er executionRepository) UpdateState(token string, state executions.State) error {
	s := er.db.Copy()
	defer s.Close()

	c := s.DB(dbName).C(execCollection)
	u := bson.M{"$set": bson.M{"state": state}}
	if err := c.Update(bson.M{"token": token}, u); err != nil {
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
		ID:                       id,
		Owner:                    owner,
		Algo:                     exec.Algo,
		Data:                     exec.Data,
		AdditionalLocalJobArgs:   exec.AdditionalLocalJobArgs,
		Type:                     exec.Type,
		GlobalTimeout:            exec.GlobalTimeout,
		LocalTimeout:             exec.LocalTimeout,
		AdditionalPreprocessArgs: exec.AdditionalPreprocessArgs,
		Mode:                     exec.Mode,
		AdditionalGlobalJobArgs:  exec.AdditionalGlobalJobArgs,
		AdditionalFiles:          exec.AdditionalFiles,
		State:                    exec.State,
		Token:                    exec.Token,
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
			ID:                       exec.ID.Hex(),
			Owner:                    owner,
			Algo:                     exec.Algo,
			Data:                     exec.Data,
			AdditionalLocalJobArgs:   exec.AdditionalLocalJobArgs,
			Type:                     exec.Type,
			GlobalTimeout:            exec.GlobalTimeout,
			LocalTimeout:             exec.LocalTimeout,
			AdditionalPreprocessArgs: exec.AdditionalPreprocessArgs,
			Mode:                     exec.Mode,
			AdditionalGlobalJobArgs:  exec.AdditionalGlobalJobArgs,
			AdditionalFiles:          exec.AdditionalFiles,
			State:                    exec.State,
			Token:                    exec.Token,
		})
	}

	return es, nil
}

type execution struct {
	ID                       bson.ObjectId      `bson:"_id"`
	Owner                    string             `bson:"owner"`
	Algo                     string             `bson:"algo"`
	Data                     string             `bson:"data"`
	AdditionalLocalJobArgs   []string           `bson:"additional_local_job_args"`
	Type                     string             `bson:"type"`
	GlobalTimeout            uint64             `bson:"global_timeout"`
	LocalTimeout             uint64             `bson:"local_timeout"`
	AdditionalPreprocessArgs []string           `bson:"additional_preprocess_args"`
	Mode                     executions.JobMode `bson:"mode"`
	AdditionalGlobalJobArgs  []string           `bson:"additional_global_job_args"`
	AdditionalFiles          []string           `bson:"additional_files"`
	State                    executions.State   `bson:"state"`
	Token                    string             `bson:"token"`
}

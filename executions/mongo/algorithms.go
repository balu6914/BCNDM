package mongo

import (
	"datapace/executions"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ executions.AlgorithmRepository = (*algorithmRepository)(nil)

type algorithmRepository struct {
	db *mgo.Session
}

// NewAlgorithmRepository returns new algorithm repository instance.
func NewAlgorithmRepository(db *mgo.Session) executions.AlgorithmRepository {
	return algorithmRepository{db: db}
}

func (repo algorithmRepository) Create(algo executions.Algorithm) error {
	s := repo.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(algo.ID) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(algoCollection)

	a := algorithm{
		ID:         bson.ObjectIdHex(algo.ID),
		Name:       algo.Name,
		Path:       algo.Path,
		ModelToken: algo.ModelToken,
		ModelName:  algo.ModelName,
	}

	if err := c.Insert(a); err != nil {
		if mgo.IsDup(err) {
			return executions.ErrConflict
		}
		return err
	}

	return nil
}

func (repo algorithmRepository) One(id string) (executions.Algorithm, error) {
	s := repo.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return executions.Algorithm{}, executions.ErrMalformedData
	}

	c := s.DB(dbName).C(algoCollection)

	var algo algorithm
	if err := c.FindId(bson.ObjectIdHex(id)).One(&algo); err != nil {
		if err == mgo.ErrNotFound {
			return executions.Algorithm{}, executions.ErrNotFound
		}
		return executions.Algorithm{}, err
	}

	return executions.Algorithm{
		ID:         algo.ID.Hex(),
		Name:       algo.Name,
		Path:       algo.Path,
		ModelToken: algo.ModelToken,
		ModelName:  algo.ModelName,
	}, nil
}

type algorithm struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"`
	Path       string        `bson:"path"`
	ModelToken string        `bson:"model_token,omitempty"`
	ModelName  string        `bson:"model_name,omitempty"`
}

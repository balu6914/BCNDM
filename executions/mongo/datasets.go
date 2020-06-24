package mongo

import (
	"github.com/datapace/datapace/executions"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ executions.DatasetRepository = (*datasetRepository)(nil)

type datasetRepository struct {
	db *mgo.Session
}

// NewDatasetRepository returns new dataset repository instance.
func NewDatasetRepository(db *mgo.Session) executions.DatasetRepository {
	return datasetRepository{db: db}
}

func (repo datasetRepository) Create(data executions.Dataset) error {
	s := repo.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(data.ID) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(dataCollection)

	d := dataset{
		ID:         bson.ObjectIdHex(data.ID),
		ExternalID: data.ExternalID,
		Metadata:   data.Metadata,
	}

	if err := c.Insert(d); err != nil {
		if mgo.IsDup(err) {
			return executions.ErrConflict
		}
		return err
	}

	return nil
}

func (repo datasetRepository) Update(data executions.Dataset) error {
	s := repo.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(data.ID) {
		return executions.ErrMalformedData
	}

	c := s.DB(dbName).C(dataCollection)

	d := dataset{
		ID:         bson.ObjectIdHex(data.ID),
		ExternalID: data.ExternalID,
		Metadata:   data.Metadata,
	}
	if err := c.UpdateId(d.ID, d); err != nil {
		return err
	}

	return nil
}

func (repo datasetRepository) One(id string) (executions.Dataset, error) {
	s := repo.db.Copy()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return executions.Dataset{}, executions.ErrMalformedData
	}

	c := s.DB(dbName).C(dataCollection)

	var data dataset
	if err := c.FindId(bson.ObjectIdHex(id)).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			return executions.Dataset{}, executions.ErrNotFound
		}
		return executions.Dataset{}, err
	}

	return executions.Dataset{
		ID:       data.ID.Hex(),
		Metadata: data.Metadata,
	}, nil
}

type dataset struct {
	ID         bson.ObjectId     `bson:"_id"`
	ExternalID string            `bson:"external_id"`
	Metadata   map[string]string `bson:"metadata"`
}

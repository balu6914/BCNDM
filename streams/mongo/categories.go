package mongo

import (
	"github.com/datapace/datapace/auth"
	"github.com/datapace/datapace/streams"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ streams.CategoryRepository = (*categoryRepository)(nil)

const (
	categoryCollection = "categories"
)

type categoryRepository struct {
	db *mgo.Session
}

// NewExecutionRepository returns new execution repository instance.
func NewCategoryRepo(db *mgo.Session) streams.CategoryRepository {
	c := db.DB(dbName).C(categoryCollection)
	indices := []mgo.Index{
		mgo.Index{
			Name: "names",
			Key:  []string{"name"},
		},
		mgo.Index{
			Name: "parentids",
			Key:  []string{"parentid"},
		},
	}
	for _, idx := range indices {
		c.EnsureIndex(idx)
	}

	return &categoryRepository{db}
}

type Category struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `bson:"name"`
	ParentID string        `bson:"parentid"`
}

func (ca categoryRepository) Save(categoryName string, subCategoryNames []string) (string, error) {
	s := ca.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(categoryCollection)

	dbs, _ := toMongoCategory(categoryName, "")

	if err := c.Insert(dbs); err != nil {
		if mgo.IsDup(err) {
			return "", streams.ErrConflict
		}
		return "", err
	}
	for _, subCategory := range subCategoryNames {
		subcat, _ := toMongoCategory(subCategory, string(dbs.ID.Hex()))
		if err := c.Insert(subcat); err != nil {
			return "", streams.ErrConflict
		}
	}

	return dbs.ID.Hex(), nil
}

func (ca *categoryRepository) OneByName(categoryName string) (Category, error) {
	session := ca.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(categoryCollection)

	ct := Category{}

	if err := collection.Find(bson.M{"name": categoryName}).One(&ct); err != nil {
		if err == mgo.ErrNotFound {
			return Category{}, auth.ErrNotFound
		}
		return Category{}, err
	}

	return ct, nil
}

func toMongoCategory(categoryName string, parentID string) (Category, error) {

	mcat := Category{
		ID:       bson.NewObjectId(),
		Name:     categoryName,
		ParentID: parentID,
	}

	return mcat, nil
}

func (ca *categoryRepository) List(key string) ([]streams.Category, error) {
	session := ca.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(categoryCollection)

	var catlist []streams.Category
	query := bson.M{}
	if err := collection.Find(query).All(&catlist); err != nil {
		if err == mgo.ErrNotFound {
			return catlist, auth.ErrNotFound
		}
		return catlist, err
	}

	return catlist, nil
}

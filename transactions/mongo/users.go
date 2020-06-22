package mongo

import (
	"github.com/datapace/transactions"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ transactions.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *mgo.Session
}

// NewUserRepository returns mongoDB specific instance of user repository.
func NewUserRepository(ms *mgo.Session) transactions.UserRepository {
	return userRepository{db: ms}
}

func (ur userRepository) Save(user transactions.User) error {
	mu, err := toMongoUser(user)
	if err != nil {
		return err
	}

	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	if err := collection.Insert(mu); err != nil {
		if mgo.IsDup(err) {
			return transactions.ErrConflict
		}

		return err
	}

	return nil
}

func (ur userRepository) Remove(id string) error {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	if !bson.IsObjectIdHex(id) {
		return transactions.ErrMalformedEntity
	}

	oid := bson.ObjectIdHex(id)
	if err := collection.Remove(bson.M{"_id": oid}); err != nil {
		if err == mgo.ErrNotFound {
			return transactions.ErrNotFound
		}
		return err
	}

	return nil
}

type mongoUser struct {
	ID     bson.ObjectId `bson:"_id"`
	Secret string        `bson:"secret"`
}

func toMongoUser(user transactions.User) (mongoUser, error) {
	if !bson.IsObjectIdHex(user.ID) {
		return mongoUser{}, transactions.ErrMalformedEntity
	}

	id := bson.ObjectIdHex(user.ID)
	return mongoUser{
		ID:     id,
		Secret: user.Secret,
	}, nil
}

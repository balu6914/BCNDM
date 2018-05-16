package mongo

import (
	"monetasa/auth"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ auth.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *mgo.Session
}

// NewUserRepository instantiates a PostgreSQL implementation of user
// repository.
func NewUserRepository(db *mgo.Session) auth.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Save(user auth.User) error {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Insert(user); err != nil {
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}

		return err
	}

	return nil
}

func (ur *userRepository) Update(user auth.User) error {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	query := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	if err := c.Update(query, update); err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) One(id string) (auth.User, error) {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	user := auth.User{}

	// Check if ID is an Email or an ObjectID and do proper query on mongo
	if bson.IsObjectIdHex(id) {
		_id := bson.ObjectIdHex(id)
		if err := c.Find(bson.M{"_id": _id}).One(&user); err != nil {
			if err == mgo.ErrNotFound {
				return user, auth.ErrNotFound
			}
		}
	}

	if govalidator.IsEmail(id) {
		if err := c.Find(bson.M{"email": id}).One(&user); err != nil {
			if err == mgo.ErrNotFound {
				return user, auth.ErrNotFound
			}
		}
	}

	return user, nil
}

func (ur *userRepository) Remove(id string) error {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	_id := bson.ObjectIdHex(id)
	if err := c.Remove(bson.M{"_id": _id}); err != nil {
		return err
	}

	return nil
}

package mongo

import (
	"monetasa/auth"

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

	query := bson.M{"email": user.Email}
	update := bson.M{"$set": user}
	if err := c.Update(query, update); err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) One(email string) (auth.User, error) {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	user := auth.User{}

	if err := c.Find(bson.M{"email": email}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return user, auth.ErrNotFound
		}

		return user, err
	}

	return user, nil
}

func (ur *userRepository) All() ([]auth.User, error) {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	users := []auth.User{}

	if err := c.Find(nil).All(&users); err != nil {
		return users, err
	}

	return users, nil
}

func (ur *userRepository) Remove(email string) error {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	if err := c.Remove(bson.M{"email": email}); err != nil {
		return err
	}

	return nil
}

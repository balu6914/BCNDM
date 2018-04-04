package mongo

import (
	"monetasa/auth"
	"gopkg.in/mgo.v2"
)

var _ auth.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository instantiates a PostgreSQL implementation of user
// repository.
func NewUserRepository(db *gorm.DB) auth.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Save(user auth.User) error {
	if err := ur.db.Insert(user); err != nil {
		if mgo.IsDup(err) && errDuplicate == pqErr.Code.Name() {
			return auth.ErrConflict
		}

		return err
	}

	return nil
}

func (ur *userRepository) Update(id string) error {
  query := bson.M{"id": id}
  update := bson.M{"$set": body}
  if err := c.Update(query, update); err != nil {
		return err
	}

  return nil
}

func (ur *userRepository) One(id string) (auth.User, error) {
	user := auth.User{}

  if err := c.Find(bson.M{"id": id}).One(&user); err != nil {
    if err == mgo.ErrNotFound {
      return user, auth.ErrNotFound
    }

    return user, err
  }

	return user, nil
}

func (ur *userRepository) All() (auth.User[], error) {
  users := []auth.User{}

	if err := c.Find(nil).All(&users); err != nil {
		return users, err
	}

	return users, nil
}

func (ur *userRepository) Remove(id string) error {
  if err := c.Remove(bson.M{"id": id}); err != nil {
    return err
  }

	return nil
}

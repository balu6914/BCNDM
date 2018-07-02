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

// NewUserRepository instantiates a Mongo implementation of user
// repository.
func NewUserRepository(db *mgo.Session) auth.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Save(user auth.User) error {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	// Verify if Email is already taken
	if count, _ := c.Find(bson.M{"email": user.Email}).Count(); count != 0 {
		return auth.ErrConflict
	}

	if err := collection.Insert(user); err != nil {
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}
		return err
	}

	return nil
}

func (ur *userRepository) Update(user auth.User) error {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	query := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	if err := collection.Update(query, update); err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) OneByID(id string) (auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	user := auth.User{}

	if bson.IsObjectIdHex(id) {
		objectID := bson.ObjectIdHex(id)
		if err := collection.Find(bson.M{"_id": objectID}).One(&user); err != nil {
			if err == mgo.ErrNotFound {
				return user, auth.ErrNotFound
			}
		}
	}

	return user, nil
}

func (ur *userRepository) OneByEmail(email string) (auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	user := auth.User{}

	if govalidator.IsEmail(email) {
		if err := collection.Find(bson.M{"email": email}).One(&user); err != nil {
			if err == mgo.ErrNotFound {
				return user, auth.ErrNotFound
			}
		}
	}

	return user, nil
}

func (ur *userRepository) Remove(id string) error {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(collectionName)

	if bson.IsObjectIdHex(id) {
		objectID := bson.ObjectIdHex(id)
		if err := collection.Remove(bson.M{"_id": objectID}); err != nil {
			return err
		}

		return nil
	}

	return auth.ErrUnauthorizedAccess
}

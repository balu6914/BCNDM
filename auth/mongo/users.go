package mongo

import (
	"fmt"
	"monetasa/auth"
	"monetasa/auth/fabric"

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

	if count, _ := c.Find(bson.M{"email": user.Email}).Count(); count != 0 {
		return auth.ErrConflict
	}

	if err := c.Insert(user); err != nil {
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}

		return err
	}

	// Find user structure to get generated mongoID
	if err := c.Find(bson.M{"email": user.Email}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrConflict
		}
	}

	// Create New user in Fabric network calling fabric-ca
	newUser, err := fabric.CreateUser(user.ID.Hex(), user.Password)
	if err != nil {
		fmt.Printf("Unable to create a user in the fabric-ca %v\n", err)
		return auth.ErrConflict
	}

	fmt.Printf("User created!: %v\n", newUser)
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

func (ur *userRepository) OneById(id string) (auth.User, error) {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	user := auth.User{}

	if bson.IsObjectIdHex(id) {
		_id := bson.ObjectIdHex(id)
		if err := c.Find(bson.M{"_id": _id}).One(&user); err != nil {
			if err == mgo.ErrNotFound {
				return user, auth.ErrNotFound
			}
		}
	}

	// Get user balance
	balance, err := fabric.Balance(user.ID.Hex())
	if err != nil {
		return auth.User{}, fmt.Errorf("Error fetching balance: %v\n", err)
	}
	user.Balance = balance

	return user, nil
}

func (ur *userRepository) OneByEmail(email string) (auth.User, error) {
	s := ur.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	user := auth.User{}

	if govalidator.IsEmail(email) {
		if err := c.Find(bson.M{"email": email}).One(&user); err != nil {
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

	if bson.IsObjectIdHex(id) {
		_id := bson.ObjectIdHex(id)
		if err := c.Remove(bson.M{"_id": _id}); err != nil {
			return err
		}

		return nil
	}

	return auth.ErrUnauthorizedAccess
}

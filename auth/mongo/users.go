package mongo

import (
	"datapace/auth"

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

func (ur *userRepository) Save(user auth.User) (string, error) {
	mu, err := toMongoUser(user)
	if err != nil {
		return "", err
	}

	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	if err := collection.Insert(mu); err != nil {
		if mgo.IsDup(err) {
			return "", auth.ErrConflict
		}
		return "", err
	}

	return mu.ID.Hex(), nil
}

func (ur *userRepository) Update(user auth.User) error {
	mu, err := toMongoUser(user)
	if err != nil {
		return err
	}

	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	if err := collection.UpdateId(mu.ID, bson.M{"$set": mu}); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		if mgo.IsDup(err) {
			return auth.ErrConflict
		}

		return err
	}

	return nil
}

func (ur *userRepository) OneByID(id string) (auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	mu := mongoUser{}
	if !bson.IsObjectIdHex(id) {
		return auth.User{}, auth.ErrMalformedEntity
	}

	oid := bson.ObjectIdHex(id)
	if err := collection.Find(bson.M{"_id": oid}).One(&mu); err != nil {
		if err == mgo.ErrNotFound {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, err
	}

	return mu.toUser(), nil
}

func (ur *userRepository) OneByEmail(email string) (auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	mu := mongoUser{}
	if err := collection.Find(bson.M{"email": email}).One(&mu); err != nil {
		if err == mgo.ErrNotFound {
			return auth.User{}, auth.ErrNotFound
		}
		return auth.User{}, err
	}

	return mu.toUser(), nil
}

func (ur *userRepository) Remove(id string) error {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	if !bson.IsObjectIdHex(id) {
		return auth.ErrMalformedEntity
	}

	oid := bson.ObjectIdHex(id)

	if err := collection.Remove(bson.M{"_id": oid}); err != nil {
		if err == mgo.ErrNotFound {
			return auth.ErrNotFound
		}
		return err
	}

	return nil
}

func (ur *userRepository) List() ([]auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	users := []auth.User{}
	if err := collection.Find(nil).All(&users); err != nil {
		return nil, err
	}

	return users, nil
}

type mongoUser struct {
	Email        string        `bson:"email,omitempty"`
	Password     string        `bson:"password,omitempty"`
	ContactEmail string        `bson:"contact_email,omitempty"`
	ID           bson.ObjectId `bson:"_id,omitempty"`
	FirstName    string        `bson:"first_name,omitempty"`
	LastName     string        `bson:"last_name,omitempty"`
}

func toMongoUser(user auth.User) (mongoUser, error) {
	if user.ID == "" {
		return mongoUser{
			Email:        user.Email,
			ContactEmail: user.ContactEmail,
			Password:     user.Password,
			ID:           bson.NewObjectId(),
			FirstName:    user.FirstName,
			LastName:     user.LastName,
		}, nil
	}

	if !bson.IsObjectIdHex(user.ID) {
		return mongoUser{}, auth.ErrMalformedEntity
	}

	id := bson.ObjectIdHex(user.ID)
	return mongoUser{
		Email:        user.Email,
		ContactEmail: user.ContactEmail,
		Password:     user.Password,
		ID:           id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}, nil
}

func (user mongoUser) toUser() auth.User {
	return auth.User{
		Email:        user.Email,
		ContactEmail: user.ContactEmail,
		Password:     user.Password,
		ID:           user.ID.Hex(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}
}

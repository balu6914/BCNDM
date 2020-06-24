package mongo

import (
	"github.com/datapace/datapace/auth"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ auth.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db     *mgo.Session
	cipher auth.Cipher
}

// NewUserRepository instantiates a Mongo implementation of user
// repository.
func NewUserRepository(db *mgo.Session, cipher auth.Cipher) auth.UserRepository {
	return &userRepository{
		db:     db,
		cipher: cipher,
	}
}

func (ur *userRepository) Save(user auth.User) (string, error) {
	var err error
	if user, err = ur.cipher.Encrypt(user); err != nil {
		return "", err
	}

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
	var err error
	if user, err = ur.cipher.Encrypt(user); err != nil {
		return err
	}

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

	return ur.cipher.Decrypt(mu.toUser())
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

	return ur.cipher.Decrypt(mu.toUser())
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

func (ur *userRepository) AllExcept(plist []string) ([]auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	ids := []bson.ObjectId{}
	for _, user := range plist {
		ids = append(ids, bson.ObjectIdHex(user))
	}

	q := bson.M{
		"_id": bson.M{
			"$nin": ids,
		},
	}
	mu := []mongoUser{}
	if err := collection.Find(q).All(&mu); err != nil {
		return nil, err
	}

	users := []auth.User{}
	for _, u := range mu {
		users = append(users, u.toUser())
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
	Company      string        `bson:"company,omitempty"`
	Address      string        `bson:"address,omitempty"`
	Phone        string        `bson:"phone,omitempty"`
	Roles        []string      `bson:"roles,omitempty"`
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
			Company:      user.Company,
			Address:      user.Address,
			Phone:        user.Phone,
			Roles:        user.Roles,
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
		Company:      user.Company,
		Address:      user.Address,
		Phone:        user.Phone,
		Roles:        user.Roles,
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
		Company:      user.Company,
		Address:      user.Address,
		Phone:        user.Phone,
		Roles:        user.Roles,
	}
}

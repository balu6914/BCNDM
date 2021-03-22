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

	policies, err := ur.listPolicies(mu.Policies)
	if err != nil {
		return auth.User{}, err
	}

	u, err := ur.cipher.Decrypt(mu.toUser())
	if err != nil {
		return auth.User{}, err
	}

	u.Policies = policies
	return u, nil
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

	policies, err := ur.listPolicies(mu.Policies)
	if err != nil {
		return auth.User{}, err
	}

	u, err := ur.cipher.Decrypt(mu.toUser())
	if err != nil {
		return auth.User{}, err
	}

	u.Policies = policies
	return u, nil
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

func (ur *userRepository) RetrieveAll(filters auth.AdminFilters) ([]auth.User, error) {
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(usersCollection)

	// Create roles filters
	var mRoles []bson.M
	for _, r := range filters.Roles {
		role := bson.M{
			"role": r,
		}
		mRoles = append(mRoles, role)
	}

	// Create status filters
	var mStatus []bson.M
	if !filters.Locked {
		locked := bson.M{
			"locked": bson.M{"$eq": false},
		}
		mStatus = append(mStatus, locked)
	}
	if !filters.Disabled {
		disabled := bson.M{
			"disabled": bson.M{"$eq": false},
		}
		mStatus = append(mStatus, disabled)
	}

	// Create final Mongo filter with roles and status
	q := bson.M{
		"$or": mRoles,
	}
	if !filters.Locked || !filters.Disabled {
		q = bson.M{
			"$or":  mRoles,
			"$and": mStatus,
		}
	}

	mu := []mongoUser{}
	if err := collection.Find(q).All(&mu); err != nil {
		return nil, err
	}

	users := []auth.User{}
	for _, u := range mu {
		du, err := ur.cipher.Decrypt(u.toUser())
		if err != nil {
			return nil, err
		}
		users = append(users, du)
	}

	return users, nil
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
		du, err := ur.cipher.Decrypt(u.toUser())
		if err != nil {
			return nil, err
		}
		users = append(users, du)
	}

	return users, nil
}

func (ur *userRepository) listPolicies(ids []string) ([]auth.Policy, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	session := ur.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(policiesCollection)

	mp := []mongoPolicy{}
	mongoIds := []bson.ObjectId{}
	for _, id := range ids {
		if !bson.IsObjectIdHex(id) {
			return nil, auth.ErrMalformedEntity
		}
		oid := bson.ObjectIdHex(id)
		mongoIds = append(mongoIds, oid)
	}
	q := bson.M{
		"_id": bson.M{
			"$in": mongoIds,
		},
	}
	if err := collection.Find(q).All(&mp); err != nil {
		if err == mgo.ErrNotFound {
			return []auth.Policy{}, auth.ErrNotFound
		}
		return []auth.Policy{}, err
	}

	ret := []auth.Policy{}
	for _, p := range mp {
		ret = append(ret, toPolicy(p))
	}

	return ret, nil
}

type mongoUser struct {
	Email           string        `bson:"email,omitempty"`
	Password        string        `bson:"password,omitempty"`
	ContactEmail    string        `bson:"contact_email,omitempty"`
	ID              bson.ObjectId `bson:"_id,omitempty"`
	FirstName       string        `bson:"first_name,omitempty"`
	LastName        string        `bson:"last_name,omitempty"`
	Company         string        `bson:"company,omitempty"`
	Address         string        `bson:"address,omitempty"`
	Phone           string        `bson:"phone,omitempty"`
	Role            string        `bson:"role,omitempty"`
	Disabled        *bool         `bson:"disabled,omitempty"`
	Policies        []string      `bson:"policies,omitempty"`
	Attempt         int           `bson:"attempt,omitempty"`
	Locked          *bool         `bson:"locked,omitempty"`
	PasswordHistory []string      `bson:"password_history,omitempty"`
}

func toMongoUser(user auth.User) (mongoUser, error) {
	var policies []string
	n := len(user.Policies)
	if n > 0 {
		policies = make([]string, n)
		for i, p := range user.Policies {
			policies[i] = p.ID
		}
	}

	mu := mongoUser{
		Email:           user.Email,
		ContactEmail:    user.ContactEmail,
		Password:        user.Password,
		ID:              bson.NewObjectId(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Company:         user.Company,
		Address:         user.Address,
		Phone:           user.Phone,
		Role:            user.Role,
		Disabled:        &user.Disabled,
		Policies:        policies,
		Locked:          &user.Locked,
		Attempt:         user.Attempt,
		PasswordHistory: user.PasswordHistory,
	}
	if user.ID == "" {
		mu.ID = bson.NewObjectId()
		return mu, nil
	}

	if !bson.IsObjectIdHex(user.ID) {
		return mongoUser{}, auth.ErrMalformedEntity
	}

	mu.ID = bson.ObjectIdHex(user.ID)
	return mu, nil
}

func (user mongoUser) toUser() auth.User {
	disabled := false
	if user.Disabled != nil {
		disabled = *user.Disabled
	}
	locked := false
	if user.Locked != nil {
		locked = *user.Locked
	}

	return auth.User{
		Email:           user.Email,
		ContactEmail:    user.ContactEmail,
		Password:        user.Password,
		ID:              user.ID.Hex(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Company:         user.Company,
		Address:         user.Address,
		Phone:           user.Phone,
		Role:            user.Role,
		Disabled:        disabled,
		Locked:          locked,
		Attempt:         user.Attempt,
		PasswordHistory: user.PasswordHistory,
	}
}

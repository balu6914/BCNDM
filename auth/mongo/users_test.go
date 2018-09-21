package mongo_test

import (
	"fmt"
	"monetasa/auth"
	"monetasa/auth/mongo"
	"testing"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const invalid = "invalid"

var db *mgo.Session

func TestSave(t *testing.T) {
	repo := mongo.NewUserRepository(db)
	user := auth.User{
		Email:        "john.doe@email.com",
		ContactEmail: "john.doe@email.com",
		Password:     "pass",
		FirstName:    "John",
		LastName:     "Doe",
	}

	cases := []struct {
		desc string
		user auth.User
		err  error
	}{
		{
			desc: "save new user",
			user: user,
			err:  nil,
		},
		{
			desc: "save existing user",
			user: user,
			err:  auth.ErrConflict,
		},
		{
			desc: "save invalid user",
			user: auth.User{ID: invalid},
			err:  auth.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		_, err := repo.Save(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	repo := mongo.NewUserRepository(db)

	user := auth.User{
		Email:        "john.doe1@email.com",
		ContactEmail: "john.doe1@email.com",
		Password:     "pass",
		FirstName:    "John",
		LastName:     "Doe",
	}
	id, _ := repo.Save(user)
	otherUser := user
	user.ID = id

	otherUser.Email = "john.doe2@email.com"
	otherID, _ := repo.Save(otherUser)
	otherUser.ID = otherID
	otherUser.Email = "john.doe1@email.com"
	otherUser.ContactEmail = "john.doe1@email.com"
	otherUser.FirstName = "john1"
	otherUser.LastName = "doe2"

	cases := []struct {
		desc string
		user auth.User
		err  error
	}{
		{
			desc: "update existing user",
			user: user,
			err:  nil,
		},
		{
			desc: "update non-existent user",
			user: auth.User{
				ID:       bson.NewObjectId().Hex(),
				Email:    "non.existent@email.com",
				Password: "pass",
			},
			err: auth.ErrNotFound,
		},
		{
			desc: "update user with invalid id",
			user: auth.User{ID: invalid},
			err:  auth.ErrMalformedEntity,
		},
		{
			desc: "update existing user with conflict data",
			user: otherUser,
			err:  auth.ErrConflict,
		},
	}

	for _, tc := range cases {
		err := repo.Update(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestOneByID(t *testing.T) {
	repo := mongo.NewUserRepository(db)

	user := auth.User{
		Email:     "john.doe4@email.com",
		Password:  "pass",
		FirstName: "John",
		LastName:  "Doe",
	}
	id, _ := repo.Save(user)
	user.ID = id

	cases := map[string]struct {
		id   string
		user auth.User
		err  error
	}{
		"get existing user by id": {
			id:   user.ID,
			user: user,
			err:  nil,
		},
		"get non-existent user by id": {
			id:   bson.NewObjectId().Hex(),
			user: auth.User{},
			err:  auth.ErrNotFound,
		},
		"get user by invalid id": {
			id:   invalid,
			user: auth.User{},
			err:  auth.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		u, err := repo.OneByID(tc.id)
		assert.Equal(t, tc.user, u, fmt.Sprintf("%s: expected %v got %v", desc, tc.user, u))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestOneByEmail(t *testing.T) {
	repo := mongo.NewUserRepository(db)

	user := auth.User{
		Email:     "john.doe5@email.com",
		Password:  "pass",
		FirstName: "John",
		LastName:  "Doe",
	}
	id, _ := repo.Save(user)
	user.ID = id

	cases := map[string]struct {
		email string
		user  auth.User
		err   error
	}{
		"get existing user by email": {
			email: user.Email,
			user:  user,
			err:   nil,
		},
		"get non-existent user by email": {
			email: "",
			user:  auth.User{},
			err:   auth.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		u, err := repo.OneByEmail(tc.email)
		assert.Equal(t, tc.user, u, fmt.Sprintf("%s: expected %v got %v", desc, tc.user, u))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestRemove(t *testing.T) {
	repo := mongo.NewUserRepository(db)

	user := auth.User{
		Email:        "john.doe6@email.com",
		ContactEmail: "john.doe6@email.com",
		Password:     "pass",
		FirstName:    "John",
		LastName:     "Doe",
	}
	id, _ := repo.Save(user)
	user.ID = id

	cases := map[string]struct {
		id  string
		err error
	}{
		"remove existing user by id": {
			id:  user.ID,
			err: nil,
		},
		"remove non-existent user by id": {
			id:  bson.NewObjectId().Hex(),
			err: auth.ErrNotFound,
		},
		"remove user with invalid id": {
			id:  invalid,
			err: auth.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		err := repo.Remove(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

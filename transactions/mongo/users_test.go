package mongo_test

import (
	"fmt"
	"testing"

	"github.com/datapace/datapace/transactions"
	"github.com/datapace/datapace/transactions/mongo"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const invalid = "invalid"

var db *mgo.Session

func TestSave(t *testing.T) {
	repo := mongo.NewUserRepository(db)
	user := transactions.User{
		ID:     bson.NewObjectId().Hex(),
		Secret: "secret",
	}

	cases := []struct {
		desc string
		user transactions.User
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
			err:  transactions.ErrConflict,
		},
		{
			desc: "save invalid user",
			user: transactions.User{ID: invalid},
			err:  transactions.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		err := repo.Save(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestRemove(t *testing.T) {
	repo := mongo.NewUserRepository(db)

	user := transactions.User{
		ID:     bson.NewObjectId().Hex(),
		Secret: "secret",
	}
	repo.Save(user)

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
			err: transactions.ErrNotFound,
		},
		"remove user with invalid id": {
			id:  invalid,
			err: transactions.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		err := repo.Remove(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

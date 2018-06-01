package auth_test

import (
	"fmt"
	"monetasa/auth"
	"monetasa/auth/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const wrong string = "wrong-value"

var user auth.User = auth.User{"user@example.com", "password", bson.NewObjectId(), 0}

func newService() auth.Service {
	users := mocks.NewUserRepository()
	hasher := mocks.NewHasher()
	idp := mocks.NewIdentityProvider()

	return auth.New(users, hasher, idp)
}

func TestRegister(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc string
		user auth.User
		err  error
	}{
		{"register new user", user, nil},
		{"register existing user", user, auth.ErrConflict},
	}

	for _, tc := range cases {
		err := svc.Register(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestView(t *testing.T) {
	svc := newService()
	svc.Register(user)
	key, _ := svc.Login(user)

	cases := map[string]struct {
		key string
		err error
	}{
		"View existing user":     {key, nil},
		"View non-existing user": {wrong, auth.ErrUnauthorizedAccess},
	}

	for desc, tc := range cases {
		_, err := svc.View(tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()
	svc.Register(user)
	key, _ := svc.Login(user)

	user.Password = "newPassword"
	user.Email = "new@email.com"
	cases := map[string]struct {
		key  string
		user auth.User
		err  error
	}{
		"Update user email and password":     {key, user, nil},
		"Update user with wrong credentials": {wrong, user, auth.ErrUnauthorizedAccess},
	}

	for desc, tc := range cases {
		err := svc.Update(tc.key, tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestDelete(t *testing.T) {
	svc := newService()
	svc.Register(user)
	key, _ := svc.Login(user)

	cases := map[string]struct {
		key string
		err error
	}{
		"Delete user":                        {key, nil},
		"Delete user with wrong credentials": {wrong, auth.ErrUnauthorizedAccess},
	}

	for desc, tc := range cases {
		err := svc.Delete(tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestLogin(t *testing.T) {
	svc := newService()
	svc.Register(user)

	cases := map[string]struct {
		user auth.User
		err  error
	}{
		"login with good credentials": {user, nil},
		"login with wrong e-mail":     {auth.User{wrong, user.Password, user.ID, user.Balance}, auth.ErrUnauthorizedAccess},
		"login with wrong password":   {auth.User{user.Email, wrong, user.ID, user.Balance}, auth.ErrUnauthorizedAccess},
	}

	for desc, tc := range cases {
		_, err := svc.Login(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestIdentity(t *testing.T) {
	svc := newService()
	svc.Register(user)
	key, _ := svc.Login(user)

	cases := map[string]struct {
		key string
		err error
	}{
		"valid token's identity":   {key, nil},
		"invalid token's identity": {"", auth.ErrUnauthorizedAccess},
	}

	for desc, tc := range cases {
		id, err := svc.Identity(tc.key)
		if id != user.ID.Hex() {
			err = auth.ErrUnauthorizedAccess
		}
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

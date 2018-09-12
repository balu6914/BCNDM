package auth_test

import (
	"fmt"
	"monetasa/auth"
	"monetasa/auth/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

const wrong string = "wrong-value"

var user = auth.User{
	Email:     "user@example.com",
	Password:  "password",
	ID:        "1",
	FirstName: "first",
	LastName:  "last",
}

func newService() auth.Service {
	users := mocks.NewUserRepository()
	hasher := mocks.NewHasher()
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()

	return auth.New(users, hasher, idp, ts)
}

func TestRegister(t *testing.T) {
	svc := newService()
	invalidUser := user
	invalidUser.Password = ""

	cases := []struct {
		desc string
		user auth.User
		err  error
	}{
		{
			desc: "register new user",
			user: user,
			err:  nil,
		},
		{
			desc: "register user with invalid data",
			user: invalidUser,
			err:  auth.ErrMalformedEntity,
		},
		{
			desc: "register existing user",
			user: user,
			err:  auth.ErrConflict,
		},
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
		"view existing user": {
			key: key,
			err: nil,
		},
		"view non-existing user": {
			key: wrong,
			err: auth.ErrUnauthorizedAccess,
		},
		"view user with empty key": {
			key: "",
			err: auth.ErrUnauthorizedAccess,
		},
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
	invalidUser := user
	invalidUser.Password = ""

	cases := []struct {
		desc string
		key  string
		user auth.User
		err  error
	}{
		{
			desc: "update user email and password",
			key:  key,
			user: user,
			err:  nil,
		},
		{
			desc: "update user with invalid data",
			key:  key,
			user: invalidUser,
			err:  auth.ErrMalformedEntity,
		},
		{
			desc: "update user with invalid credentials",
			key:  "",
			user: user,
			err:  auth.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		err := svc.Update(tc.key, tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestLogin(t *testing.T) {
	svc := newService()
	svc.Register(user)

	user2 := user
	user2.Email = wrong

	user3 := user
	user3.Password = wrong

	cases := map[string]struct {
		user auth.User
		err  error
	}{
		"login with good credentials": {
			user: user,
			err:  nil,
		},
		"login with wrong e-mail": {
			user: user2,
			err:  auth.ErrUnauthorizedAccess,
		},
		"login with wrong password": {
			user: user3,
			err:  auth.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		_, err := svc.Login(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestIdentify(t *testing.T) {
	svc := newService()
	svc.Register(user)
	key, _ := svc.Login(user)

	cases := map[string]struct {
		key string
		err error
	}{
		"valid token's identity": {
			key: key,
			err: nil,
		},
		"invalid token's identity": {
			key: "",
			err: auth.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		_, err := svc.Identify(tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

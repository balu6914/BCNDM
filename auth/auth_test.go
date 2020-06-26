package auth_test

import (
	"fmt"
	"testing"

	"github.com/datapace/datapace/auth"
	"github.com/datapace/datapace/auth/mocks"

	"github.com/stretchr/testify/assert"
)

const wrong string = "wrong-value"

var user = auth.User{
	Email:        "user@example.com",
	ContactEmail: "user@example.com",
	Password:     "password",
	ID:           "1",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
}

var admin = auth.User{
	Email:        "admin@example.com",
	ContactEmail: "admin@example.com",
	Password:     "password",
	ID:           "admin",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
	Roles:        []string{"admin"},
}

var nonadmin = auth.User{
	Email:        "nonadmin@example.com",
	ContactEmail: "nonadmin@example.com",
	Password:     "password",
	ID:           "nonadmin",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
}

func newService() (auth.Service, string) {
	svc, key, _ := newServiceWithAdmin()
	return svc, key
}

func newServiceWithAdmin() (auth.Service, string, auth.User) {
	hasher := mocks.NewHasher()
	users := mocks.NewUserRepository(hasher, admin)
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()
	ac := mocks.NewAccessControl()
	svc := auth.New(users, hasher, idp, ts, ac)
	key, _ := svc.Login(auth.User{
		Email:    admin.Email,
		Password: admin.Password,
	})
	return svc, key, admin
}

func TestRegister(t *testing.T) {
	svc, key, _ := newServiceWithAdmin()
	invalidUser := user
	invalidUser.Password = ""
	_, err := svc.Register(key, nonadmin)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, nonadmin.ID))
	nonadminkey, _ := svc.Login(auth.User{
		Email:    nonadmin.Email,
		Password: nonadmin.Password,
	})
	cases := []struct {
		desc string
		key  string
		user auth.User
		err  error
	}{
		{
			desc: "register new user by nonadmin",
			key:  nonadminkey,
			user: user,
			err:  auth.ErrUnauthorizedAccess,
		},
		{
			desc: "register new user",
			key:  key,
			user: user,
			err:  nil,
		},
		{
			desc: "register user with invalid data",
			key:  key,
			user: invalidUser,
			err:  auth.ErrMalformedEntity,
		},
		{
			desc: "register existing user",
			key:  key,
			user: user,
			err:  auth.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := svc.Register(tc.key, tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestView(t *testing.T) {
	svc, k := newService()
	uv := auth.User{
		ID:       "testv",
		Email:    "testview@example.com",
		Password: "testpass",
	}
	uv2 := auth.User{
		ID:       "testv2",
		Email:    "testview2@example.com",
		Password: "testpass2",
	}
	id, err := svc.Register(k, uv)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, uv.ID))
	_, err = svc.Register(k, uv2)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, uv2.ID))
	key, err := svc.Login(uv)
	key2, err := svc.Login(uv2)
	_ = err

	cases := map[string]struct {
		key string
		id  string
		err error
	}{
		"view existing user as self": {
			key: key,
			id:  id,
			err: nil,
		},
		"view existing user as another": {
			key: key2,
			id:  id,
			err: auth.ErrUnauthorizedAccess,
		},
		"view existing user as admin": {
			key: k,
			id:  id,
			err: nil,
		},
		"view non-existing user": {
			key: wrong,
			id:  id,
			err: auth.ErrUnauthorizedAccess,
		},
		"view user with empty key": {
			key: "",
			id:  id,
			err: auth.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		_, err := svc.View(tc.key, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc, k := newService()
	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
	key, _ := svc.Login(user)
	user.ContactEmail = "new@email.com"

	cases := []struct {
		desc string
		key  string
		user auth.User
		err  error
	}{
		{
			desc: "update user contact email",
			key:  key,
			user: user,
			err:  nil,
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

func TestUpdatePassword(t *testing.T) {
	svc, k := newService()
	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
	key, _ := svc.Login(user)
	user.Password = "newpassword"

	cases := []struct {
		desc string
		key  string
		user auth.User
		err  error
	}{
		{
			desc: "update user password",
			key:  key,
			user: user,
			err:  nil,
		},
		{
			desc: "update user password invalid credentials",
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
	svc, k := newService()
	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
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
	svc, k := newService()
	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
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

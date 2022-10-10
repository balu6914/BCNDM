package auth_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/datapace/datapace/auth"
	"github.com/datapace/datapace/auth/mocks"

	"github.com/stretchr/testify/assert"
)

const wrong = "wrong-value"

var policies = map[string]auth.Policy{
	"admin": {
		Name:    "admin",
		Owner:   "admin",
		Version: "1.0.0",
		Rules: []auth.Rule{
			{
				Action: auth.Any,
				Type:   "user",
			},
			{
				Action: auth.Any,
				Type:   "stream",
			},
			{
				Action: auth.Any,
				Type:   "subscription",
			},
			{
				Action: auth.Any,
				Type:   "policy",
			},
			{
				Action: auth.Any,
				Type:   "contract",
			},
		},
	},
	"user": {
		Name:    "user",
		Owner:   "admin",
		Version: "1.0.0",
		Rules: []auth.Rule{
			{
				Action: auth.CreateBulk,
				Type:   "stream",
			},
			{
				Action: auth.List,
				Type:   "stream",
			},
			{
				Action: auth.List,
				Type:   "user",
			},
			{
				Action: auth.Any,
				Type:   "stream",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.Any,
				Type:   "contract",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.List,
				Type:   "subscription",
			},
			{
				Action: auth.Any,
				Type:   "subscription",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.Any,
				Type:   "user",
				Condition: auth.SimpleCondition{
					Key: "id",
				},
			},
			{
				Action: auth.Any,
				Type:   "token",
			},
		},
	},
}

var policiesMu sync.Mutex

var user = auth.User{
	Email:        "user@example.com",
	ContactEmail: "user@example.com",
	Password:     "Pass1234!",
	ID:           "1",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
	Policies:     []auth.Policy{policies["user"]},
}

var admin = auth.User{
	Email:        "admin@example.com",
	ContactEmail: "admin@example.com",
	Password:     "Pass1234!",
	ID:           "admin",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
	Role:         "admin",
	Policies:     []auth.Policy{policies["admin"]},
}

var noAdmin = auth.User{
	Email:        "noadmin@example.com",
	ContactEmail: "nonadmin@example.com",
	Password:     "Pass1234!",
	ID:           "nonadmin",
	FirstName:    "first",
	LastName:     "last",
	Company:      "company",
	Address:      "address",
	Phone:        "+1234567890",
	Policies:     []auth.Policy{policies["user"]},
}

func newService() (auth.Service, string) {
	svc, key, _ := newServiceWithAdmin()
	return svc, key
}

func newServiceWithAdmin() (auth.Service, string, auth.User) {
	hasher := mocks.NewHasher()
	users := mocks.NewUserRepository(hasher, admin, policies, &policiesMu)
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()
	ac := mocks.NewAccessControl()
	policies := mocks.NewPolicyRepository(policies, &policiesMu)
	svc := auth.New(users, policies, hasher, idp, ts, ac)
	key, _ := svc.Login(auth.User{
		Email:    admin.Email,
		Password: admin.Password,
	})
	return svc, key, admin
}

func TestRegister(t *testing.T) {
	svc, key, _ := newServiceWithAdmin()
	invalidUser := user
	invalidUser.Password = "wrong"
	_, err := svc.Register(key, noAdmin)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, noAdmin.ID))
	nonadminkey, _ := svc.Login(auth.User{
		Email:    noAdmin.Email,
		Password: noAdmin.Password,
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
			err:  auth.ErrPassLength,
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
		Password: "Pass1234!",
		Role:     auth.UserRole,
	}
	uv2 := auth.User{
		ID:       "testv2",
		Email:    "testview2@example.com",
		Password: "Pass1234!",
		Role:     auth.UserRole,
	}
	id, err := svc.Register(k, uv)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, uv.ID))
	_, err = svc.Register(k, uv2)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, uv2.ID))
	key, err := svc.Login(uv)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while login user %s", err, uv2.ID))
	key2, err := svc.Login(uv2)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while login user %s", err, uv2.ID))

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
		_, err := svc.ViewUser(tc.key, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc, k := newService()
	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
	key, _ := svc.Login(user)
	user2 := user
	user2.ContactEmail = "new@email.com"
	user2.Password = ""

	cases := []struct {
		desc string
		key  string
		user auth.User
		err  error
	}{
		{
			desc: "update user contact email",
			key:  key,
			user: user2,
			err:  nil,
		},
		{
			desc: "update user with invalid credentials",
			key:  "",
			user: user2,
			err:  auth.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateUser(tc.key, tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdatePassword(t *testing.T) {
	svc, k := newService()

	_, err := svc.Register(k, user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while registering user %s", err, user.ID))
	key, err := svc.Login(user)
	assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while login user %s", err, user.ID))
	user.Password = "Pass2222!"

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
		err := svc.UpdateUser(tc.key, tc.user)
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

func TestListUsers(t *testing.T) {
	svc, adminKey := newService()
	_, err := svc.Register(adminKey, noAdmin)
	userKey, _ := svc.Login(noAdmin)
	users, err := svc.ListUsers(userKey)
	assert.Nil(t, err)
	assert.Empty(t, users)
}

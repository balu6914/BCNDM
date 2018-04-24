package auth_test

import (
	"fmt"
	"monetasa/auth"
	"monetasa/auth/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

var user auth.User = auth.User{"user@example.com", "password"}

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

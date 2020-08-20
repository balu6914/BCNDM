package mocks

import (
	authproto "github.com/datapace/datapace/proto/auth"
	"github.com/datapace/datapace/streams"
)

type authorization struct {
	// Users are represented by its ids.
	users []string
}

// NewAuth returns mock auth service.
func NewAuth(users []string) streams.Authorization {
	return authorization{
		users: users,
	}
}

func (a authorization) Authorize(r *authproto.AuthRequest) (string, error) {
	for _, id := range a.users {
		if id == r.Token {
			return id, nil
		}
	}
	return "", streams.ErrUnauthorizedAccess
}

func (a authorization) Email(token string) (*authproto.UserEmail, error) {
	return &authproto.UserEmail{}, nil
}

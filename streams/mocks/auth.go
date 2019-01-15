package mocks

import (
	"datapace"
	"datapace/streams"
	"net/http"
)

type authorization struct {
	// Users are represented by its ids.
	users []string
}

func (a authorization) Authorize(r *http.Request) (string, error) {
	key := r.Header.Get("Authorization")
	for _, id := range a.users {
		if id == key {
			return id, nil
		}
	}
	return "", streams.ErrUnauthorizedAccess
}

func (a authorization) Email(token string) (datapace.UserEmail, error) {
	return datapace.UserEmail{}, nil
}

// NewAuth returns mock auth service.
func NewAuth(users []string) streams.Authorization {
	return authorization{
		users: users,
	}
}

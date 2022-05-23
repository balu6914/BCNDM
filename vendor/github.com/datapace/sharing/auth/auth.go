// Package auth provides API for working with user entity.
package auth

import (
	"context"
	"errors"
	"time"

	authproto "github.com/datapace/datapace/proto/auth"
)

// AuthService contains API for fetching user related data.
type AuthService interface {

	// Identifies user by the provided key.
	Identify(string) (string, error)
}

var _ AuthService = (*authService)(nil)

type authService struct {
	auth authproto.AuthServiceClient
}

// New returns new auth service implementation instance.
func New(auth authproto.AuthServiceClient) AuthService {
	return authService{auth: auth}
}

func (as authService) Identify(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := as.auth.Identify(ctx, &authproto.Token{Value: key})
	if err != nil {
		return "", errors.New("unauthorized")
	}
	return res.GetValue(), nil
}

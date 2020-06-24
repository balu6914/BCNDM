// Package auth provides API for working with user entity.
package auth

import (
	"context"
	"time"

	"github.com/datapace/datapace"

	access "github.com/datapace/datapace/access-control"
)

var _ access.AuthService = (*authService)(nil)

type authService struct {
	auth datapace.AuthServiceClient
}

// New returns new auth service implementation instance.
func New(auth datapace.AuthServiceClient) access.AuthService {
	return authService{auth: auth}
}

func (as authService) Identify(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := as.auth.Identify(ctx, &datapace.Token{Value: key})
	if err != nil {
		return "", access.ErrUnauthorizedAccess
	}

	return res.GetValue(), nil
}

func (as authService) Exists(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := as.auth.Exists(ctx, &datapace.ID{Value: id}); err != nil {
		return access.ErrNotFound
	}

	return nil
}

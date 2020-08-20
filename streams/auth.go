package streams

import (
	"context"
	"fmt"
	"time"

	log "github.com/datapace/datapace/logger"
	authproto "github.com/datapace/datapace/proto/auth"
)

var _ Authorization = (*authService)(nil)

// Authorization represents an authorization utility to
// be used in transport layer.
type Authorization interface {
	// Authorize method is used to authorize http request.
	Authorize(r *authproto.AuthRequest) (string, error)
	// Email method is used to fetch email and contactEmail for the user.
	Email(token string) (*authproto.UserEmail, error)
}

type authService struct {
	auth   authproto.AuthServiceClient
	logger log.Logger
}

// NewAuthorization method instantiates new security service.
func NewAuthorization(auth authproto.AuthServiceClient, logger log.Logger) Authorization {
	return authService{
		auth:   auth,
		logger: logger,
	}
}

func (as authService) Authorize(r *authproto.AuthRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := as.auth.Authorize(ctx, r)
	if err != nil {
		as.logger.Error(fmt.Sprintf("failed to authorize request: %s", err))
		return "", ErrUnauthorizedAccess
	}
	return res.GetValue(), nil
}

func (as authService) Email(token string) (*authproto.UserEmail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := as.auth.Email(ctx, &authproto.Token{Value: token})
	if err != nil {
		as.logger.Error(fmt.Sprintf("failed to fetch users emails: %s", err))
		return &authproto.UserEmail{}, ErrUnauthorizedAccess
	}

	return res, nil
}

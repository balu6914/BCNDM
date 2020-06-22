package streams

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/datapace"

	log "github.com/datapace/logger"
)

var _ Authorization = (*authService)(nil)

// Authorization represents an authorization utility to
// be used in transport layer.
type Authorization interface {
	// Authorize method is used to authorize http request.
	Authorize(r *http.Request) (string, error)
	// Email method is used to fetch email and contactEmail for the user.
	Email(token string) (datapace.UserEmail, error)
}

type authService struct {
	auth   datapace.AuthServiceClient
	logger log.Logger
}

// NewAuthorization method instantiates new security service.
func NewAuthorization(auth datapace.AuthServiceClient, logger log.Logger) Authorization {
	return authService{
		auth:   auth,
		logger: logger,
	}
}

func (as authService) Authorize(r *http.Request) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	key := r.Header.Get("Authorization")
	res, err := as.auth.Identify(ctx, &datapace.Token{Value: key})
	if err != nil {
		as.logger.Error(fmt.Sprintf("failed to authorize request: %s", err))
		return "", ErrUnauthorizedAccess
	}

	return res.GetValue(), nil
}

func (as authService) Email(token string) (datapace.UserEmail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := as.auth.Email(ctx, &datapace.Token{Value: token})
	if err != nil {
		as.logger.Error(fmt.Sprintf("failed to fetch users emails: %s", err))
		return datapace.UserEmail{}, ErrUnauthorizedAccess
	}

	return *res, nil
}

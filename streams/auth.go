package streams

import (
	"context"
	"fmt"
	"monetasa"
	log "monetasa/logger"
	"net/http"
	"time"
)

var _ Authorization = (*authService)(nil)

// Authorization represents an authorization utility to
// be used in transport layer.
type Authorization interface {
	// Authorize method is used to authorize http request.
	Authorize(r *http.Request) (string, error)
}

type authService struct {
	auth   monetasa.AuthServiceClient
	logger log.Logger
}

// NewAuthorization method instantiates new security service.
func NewAuthorization(auth monetasa.AuthServiceClient, logger log.Logger) Authorization {
	return authService{
		auth:   auth,
		logger: logger,
	}
}

func (ss authService) Authorize(r *http.Request) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	key := r.Header.Get("Authorization")
	res, err := ss.auth.Identify(ctx, &monetasa.Token{Value: key})
	if err != nil {
		ss.logger.Error(fmt.Sprintf("failed to authorize request: %s", err))
		return "", ErrUnauthorizedAccess
	}

	return res.GetValue(), nil
}

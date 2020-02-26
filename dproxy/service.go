package dproxy

import (
	"errors"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrTokenParsingFailed = errors.New("token parsing failed")
	ErrMalformedEntity    = errors.New("malformed entity specification")
)

type Service interface {
	CreateToken(string, int) (string, error)
	GetTargetURL(string) (string, error)
}

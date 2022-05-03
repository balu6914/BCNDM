package auth

import (
	"errors"
)

type authServiceMock struct {
}

var _ AuthService = (*authServiceMock)(nil)

func NewAuthServiceMock() AuthService {
	return &authServiceMock{}
}

func (svc authServiceMock) Identify(s string) (string, error) {
	if s == "unauthorized" {
		return "", errors.New("unauthorized")
	}
	return s, nil
}

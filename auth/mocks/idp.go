package mocks

import (
	"strings"

	"github.com/datapace/datapace/auth"
)

var _ auth.IdentityProvider = (*identityProviderMock)(nil)

type identityProviderMock struct{}

// NewIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewIdentityProvider() auth.IdentityProvider {
	return &identityProviderMock{}
}

func (idp *identityProviderMock) TemporaryKey(id string, role string) (string, error) {
	if id == "" {
		return "", auth.ErrUnauthorizedAccess
	}
	if role != "" {
		return id + "|" + role, nil
	}
	return id, nil

}

func (idp *identityProviderMock) PermanentKey(id string) (string, error) {
	return idp.TemporaryKey(id, "")
}

func (idp *identityProviderMock) Identity(key string) (string, error) {
	if key == "invalid" {
		return "", auth.ErrUnauthorizedAccess
	}
	parts := strings.Split(key, "|")
	return idp.TemporaryKey(parts[0], "")
}

func (idp *identityProviderMock) Role(key string) (string, error) {
	parts := strings.Split(key, "|")
	l := len(parts)
	var role string
	if l > 1 {
		role = parts[1]
	}
	return role, nil
}

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

func (idp *identityProviderMock) TemporaryKey(id string, roles []string) (string, error) {
	if id == "" {
		return "", auth.ErrUnauthorizedAccess
	}
	if len(roles) > 0 {
		return id + "|" + strings.Join(roles, "|"), nil
	}
	return id, nil

}

func (idp *identityProviderMock) PermanentKey(id string) (string, error) {
	return idp.TemporaryKey(id, []string{})
}

func (idp *identityProviderMock) Identity(key string) (string, error) {
	if key == "invalid" {
		return "", auth.ErrUnauthorizedAccess
	}
	parts := strings.Split(key, "|")
	l := len(parts)
	var roles []string
	if l > 1 {
		roles = parts[1:]
	}
	return idp.TemporaryKey(parts[0], roles)
}

func (idp *identityProviderMock) Roles(key string) ([]string, error) {
	parts := strings.Split(key, "|")
	l := len(parts)
	var roles []string
	if l > 1 {
		roles = parts[1:]
	}
	return roles, nil
}

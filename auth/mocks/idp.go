package mocks

import "datapace/auth"

var _ auth.IdentityProvider = (*identityProviderMock)(nil)

type identityProviderMock struct{}

// NewIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewIdentityProvider() auth.IdentityProvider {
	return &identityProviderMock{}
}

func (idp *identityProviderMock) TemporaryKey(id string) (string, error) {
	if id == "" {
		return "", auth.ErrUnauthorizedAccess
	}

	return id, nil
}

func (idp *identityProviderMock) PermanentKey(id string) (string, error) {
	return idp.TemporaryKey(id)
}

func (idp *identityProviderMock) Identity(key string) (string, error) {
	if key == "invalid" {
		return "", auth.ErrUnauthorizedAccess
	}

	return idp.TemporaryKey(key)
}

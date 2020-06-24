package mocks

import "github.com/datapace/datapace/auth"

var _ auth.Hasher = (*hasherMock)(nil)

type hasherMock struct{}

// NewHasher creates "no-op" hasher for test purposes. This implementation will
// return secrets without changing them.
func NewHasher() auth.Hasher {
	return &hasherMock{}
}

func (hm *hasherMock) Hash(pwd string) (string, error) {
	if pwd == "" {
		return "", auth.ErrMalformedEntity
	}

	return pwd, nil
}

func (hm *hasherMock) Compare(plain, hashed string) error {
	if plain != hashed {
		return auth.ErrUnauthorizedAccess
	}

	return nil
}

// Package bcrypt provides a hasher implementation utilising bcrypt.
package bcrypt

import (
	"github.com/datapace/auth"

	"golang.org/x/crypto/bcrypt"
)

const cost = 10

var _ auth.Hasher = (*bcryptHasher)(nil)

type bcryptHasher struct{}

// New instantiates a bcrypt-based hasher implementation.
func New() auth.Hasher {
	return &bcryptHasher{}
}

func (bh *bcryptHasher) Hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (bh *bcryptHasher) Compare(plain, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

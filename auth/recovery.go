package auth

import "github.com/dgrijalva/jwt-go"

// RecoveryTokenProvider specifies an API for password recovery via security
// tokens.
type RecoveryTokenProvider interface {
	// CreateTokenString generates the recovery token string.
	CreateTokenString(string, string) (string, error)

	// ParseToken validates recovery token string and extracts the entity identifier given its secret key.
	ParseToken(string, string) (*jwt.Token, error)
}

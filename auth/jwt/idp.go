// Package jwt provides a JWT identity provider.
package jwt

import (
	"time"

	"datapace/auth"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	issuer   = "datapace"
	duration = 10 * time.Hour
)

var _ auth.IdentityProvider = (*jwtIdentityProvider)(nil)

type jwtIdentityProvider struct {
	secret string
}

// New instantiates a JWT identity provider.
func New(secret string) auth.IdentityProvider {
	return &jwtIdentityProvider{}
}

func (idp *jwtIdentityProvider) TemporaryKey(id string) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(duration)

	claims := jwt.StandardClaims{
		Subject:   id,
		Issuer:    issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: exp.Unix(),
	}

	return idp.jwt(claims)
}

func (idp *jwtIdentityProvider) Identity(key string) (string, error) {
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnauthorizedAccess
		}

		return []byte(idp.secret), nil
	})
	if err != nil || !token.Valid {
		return "", auth.ErrUnauthorizedAccess
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", auth.ErrUnauthorizedAccess
	}

	return claims["sub"].(string), nil
}

func (idp *jwtIdentityProvider) jwt(claims jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(idp.secret))
}

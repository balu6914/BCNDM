package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"monetasa/auth"
)

const (
	issuer   string        = "monetasa"
	duration time.Duration = 10 * time.Hour
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

func (idp *jwtIdentityProvider) jwt(claims jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(idp.secret))
}

func (idp *jwtIdentityProvider) Identity(key string) (string, error) {
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnauthorizedAccess
		}

		return []byte(idp.secret), nil
	})

	if err != nil {
		return "", auth.ErrUnauthorizedAccess
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}

	return "", auth.ErrUnauthorizedAccess
}

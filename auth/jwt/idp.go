// Package jwt provides a JWT identity provider.
package jwt

import (
	"time"

	"github.com/datapace/auth"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	issuer   = "github.com/datapace"
	duration = 10 * time.Hour
)

var _ auth.IdentityProvider = (*jwtIdentityProvider)(nil)

type jwtIdentityProvider struct {
	secret string
}

type CustomClaims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

// New instantiates a JWT identity provider.
func New(secret string) auth.IdentityProvider {
	return &jwtIdentityProvider{}
}

func (idp *jwtIdentityProvider) TemporaryKey(id string, roles []string) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(duration)

	claims := CustomClaims{
		roles,
		jwt.StandardClaims{
			Subject:   id,
			Issuer:    issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: exp.Unix(),
		},
	}

	return idp.jwt(claims)
}

func (idp *jwtIdentityProvider) Roles(key string) ([]string, error) {
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnauthorizedAccess
		}

		return []byte(idp.secret), nil
	})
	var roles []string
	if err != nil || !token.Valid {
		return roles, auth.ErrUnauthorizedAccess
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return roles, auth.ErrUnauthorizedAccess
	}
	if _, ok := claims["roles"]; !ok {
		return roles, nil
	}
	var rr []interface{}
	if rr, ok = claims["roles"].([]interface{}); !ok {
		return roles, nil
	}
	for _, v := range rr {
		if r, ok := v.(string); ok {
			roles = append(roles, r)
		}
	}
	return roles, nil
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

func (idp *jwtIdentityProvider) jwt(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(idp.secret))
}

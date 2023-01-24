package jwt

import "github.com/dgrijalva/jwt-go"

type DproxyClaims struct {
	StandardClaims jwt.StandardClaims `json:"std"`
	URL            string             `json:"url"`
	MaxCalls       int                `json:"max_calls"`
	MaxUnit        string             `json:"max_unit"`
}

func (d DproxyClaims) Valid() error {
	return d.StandardClaims.Valid()
}

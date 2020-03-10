package jwt

import (
	"datapace/dproxy"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

type jwtService struct {
	jwtSecret string
}

var _ dproxy.TokenService = (*jwtService)(nil)

func NewService(jwtSecret string) dproxy.TokenService {
	return &jwtService{jwtSecret: jwtSecret}
}

func (d *jwtService) Create(url string, ttl int, maxCalls int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, DproxyClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ttl)).Unix(),
		},
		URL:      url,
		MaxCalls: maxCalls,
	})
	tokenString, err := token.SignedString([]byte(d.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (d *jwtService) Parse(tokenString string) (dproxy.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DproxyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, dproxy.ErrInvalidToken
		}
		return []byte(d.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*DproxyClaims); ok && token.Valid {
		return NewToken(claims.StandardClaims.Id, claims.URL, claims.MaxCalls), nil
	}
	return nil, dproxy.ErrTokenParsingFailed
}

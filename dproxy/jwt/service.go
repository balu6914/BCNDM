package jwt

import (
	"datapace/dproxy"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtService struct {
	jwtSecret string
}

var _ dproxy.Service = (*jwtService)(nil)

func NewService(jwtSecret string) dproxy.Service {
	return &jwtService{jwtSecret: jwtSecret}
}

func (d *jwtService) CreateToken(url string, ttl int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now(),
		"exp": time.Now().Add(time.Second * time.Duration(ttl)),
		"url": url,
	})
	tokenString, err := token.SignedString([]byte(d.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func (d *jwtService) GetTargetURL(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, dproxy.ErrInvalidToken
		}
		return []byte(d.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["url"].(string), nil
	}
	return "", dproxy.ErrTokenParsingFailed

}

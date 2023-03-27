package jwt

import (
	"time"

	"github.com/datapace/datapace/dproxy"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type jwtService struct {
	jwtSecret string
}

var _ dproxy.TokenService = (*jwtService)(nil)

func NewService(jwtSecret string) dproxy.TokenService {
	return &jwtService{jwtSecret: jwtSecret}
}

func (d *jwtService) Create(url string, ttl int, maxCalls int, maxUnit, subID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, DproxyClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ttl)).Unix(),
		},
		URL:      url,
		SubID:    subID,
		MaxCalls: maxCalls,
		MaxUnit:  maxUnit,
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

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, dproxy.ErrInvalidToken
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, dproxy.ErrTokenExpired
		} else {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*DproxyClaims); ok && token.Valid {
		return NewToken(claims.StandardClaims.Id, claims.URL, claims.MaxCalls, claims.MaxUnit, claims.SubID), nil
	}
	return nil, dproxy.ErrTokenParsingFailed
}

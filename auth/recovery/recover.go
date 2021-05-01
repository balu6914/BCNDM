//Package implements service that generates JWT password recovery tokens
package recovery

import (
	"errors"
	"github.com/datapace/datapace/auth"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	issuer   = "github.com/datapace/datapace"
	duration = 15 * time.Minute
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenParsingFailed = errors.New("token parsing failed")
)

var _ auth.RecoveryTokenProvider = (*RecoverTokenService)(nil)

type RecoverTokenService struct {
}

func New() auth.RecoveryTokenProvider {
	return &RecoverTokenService{}
}

type ResetPasswordClaims struct {
	StandardClaims jwt.StandardClaims
	ID             string
}

func (r ResetPasswordClaims) Valid() error {
	return r.StandardClaims.Valid()
}

func (r RecoverTokenService) CreateTokenString(id string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ResetPasswordClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    issuer,
		},
		ID: id,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (r RecoverTokenService) ParseToken(tokenString string, storedSecret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ResetPasswordClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(storedSecret), nil
	})
	if err != nil {
		return nil, ErrTokenParsingFailed
	}
	return token, nil
}

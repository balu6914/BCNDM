package dproxy

import (
	"errors"
	"time"

	"github.com/datapace/datapace/dproxy/persistence"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrTokenParsingFailed = errors.New("token parsing failed")
	ErrMalformedEntity    = errors.New("malformed entity specification")
	ErrQuotaExceeded      = errors.New("quota exceeded")
	ErrConflict           = errors.New("entity already exists")
)

type Service interface {
	CreateToken(string, int, int) (string, error)
	GetTargetURL(string) (string, error)
}

type Token interface {
	Url() string
	Uid() string
	MaxCalls() int
}

type TokenService interface {
	Create(string, int, int) (string, error)
	Parse(string) (Token, error)
}

type dService struct {
	aesKey       []byte
	tokenService TokenService
	eventsRepo   persistence.EventRepository
}

var _ Service = (*dService)(nil)

func NewService(tokenService TokenService, eventsRepo persistence.EventRepository, key []byte) Service {
	return &dService{
		tokenService: tokenService,
		eventsRepo:   eventsRepo,
		aesKey:       key,
	}
}

func (d *dService) CreateToken(url string, ttl, maxCalls int) (string, error) {
	url, err := encrypt(d.aesKey, url)
	if err != nil {
		return "", err
	}
	return d.tokenService.Create(url, ttl, maxCalls)
}

func (d *dService) GetTargetURL(tokenString string) (string, error) {
	t, err := d.tokenService.Parse(tokenString)

	calls, err := d.eventsRepo.Accumulate(persistence.Event{Time: time.Now(), Initiator: t.Uid()})
	if err != nil {
		return "", err
	}
	if t.MaxCalls() != 0 && calls > t.MaxCalls() {
		return "", ErrQuotaExceeded
	}
	url, err := decrypt(d.aesKey, t.Url())
	if err != nil {
		return "", err
	}

	return url, err
}

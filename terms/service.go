package terms

import (
	"crypto/sha256"
	"errors"
	"fmt"
	authproto "github.com/datapace/datapace/proto/auth"
	"io"
	"net/http"
)

var (
	// ErrConflict indicates usage of the existing terms id for the new terms.
	ErrConflict = errors.New("Terms ID already taken")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrFailedCreateTerms indicates that creation of terms failed.
	ErrFailedCreateTerms = errors.New("failed to create terms")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateTerms creates new Terms.
	CreateTerms(Terms) (string, error)
}

var _ Service = (*termsService)(nil)

type termsService struct {
	auth  authproto.AuthServiceClient
	terms TermsRepository
}

func (ts termsService) CreateTerms(t Terms) (string, error) {
	hash, err := makeHash(t)
	if err != nil {
		return "", err
	}
	t.TermsHash = hash
	_, err = ts.terms.Save(t)
	return "", err
}

func New(auth authproto.AuthServiceClient, terms TermsRepository) Service {
	return &termsService{
		auth:  auth,
		terms: terms,
	}
}

func makeHash(t Terms) (string, error) {
	resp, err := http.Get(t.TermsURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	hash := sha256.New()
	io.Copy(hash, resp.Body)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

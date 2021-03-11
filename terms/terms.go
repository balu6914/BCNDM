package terms

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

var _ Service = (*termsService)(nil)

type termsService struct {
	terms  TermsRepository
	ledger TermsLedger
}

func New(terms TermsRepository, ledger TermsLedger) Service {
	return &termsService{
		terms:  terms,
		ledger: ledger,
	}
}

func (ts termsService) CreateTerms(t Terms) (string, error) {
	hash, err := makeHash(t)
	if err != nil {
		return "", err
	}

	t.TermsHash = hash
	_, err = ts.terms.Save(t)
	if err != nil {
		return "", err
	}

	err = ts.ledger.CreateTerms(t)
	if err != nil {
		return "", err
	}

	return hash, err
}

func (ts termsService) ValidateTerms(t Terms) (bool, error) {
	return ts.ledger.ValidateTerms(t)
}

func makeHash(t Terms) (string, error) {
	resp, err := http.Get(t.TermsURL)
	if err != nil {
		return "", ErrFailedFetchTermsURL
	}
	if resp.StatusCode != http.StatusOK {
		return "", ErrNotFound
	}
	defer resp.Body.Close()

	hash := sha256.New()
	io.Copy(hash, resp.Body)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

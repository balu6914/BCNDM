package mocks

import (
	trs "github.com/datapace/datapace/terms"
	uuid "github.com/satori/go.uuid"
)

type mockTermsRepository struct {
	terms map[string]trs.Terms
}

func NewTermsRepository() trs.TermsRepository {
	return &mockTermsRepository{terms: map[string]trs.Terms{}}
}

func (tr mockTermsRepository) Save(terms trs.Terms) (string, error) {
	uuid := uuid.NewV4().String()
	tr.terms[uuid] = terms
	return uuid, nil
}

func (tr mockTermsRepository) One(id string) (trs.Terms, error) {
	return tr.terms[id], nil
}

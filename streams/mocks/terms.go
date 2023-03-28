package mocks

import "github.com/datapace/datapace/streams"

var _ streams.TermsService = (*termsServiceMock)(nil)

type termsServiceMock struct{}

// NewTermsService returns mock instance of Terms service.
func NewTermsService() streams.TermsService {
	return termsServiceMock{}
}

func (terms termsServiceMock) CreateTerms(_ streams.Stream) error {
	return nil
}

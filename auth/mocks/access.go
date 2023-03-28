package mocks

import "github.com/datapace/datapace/auth"

var _ auth.AccessControl = (*accessControlMock)(nil)

type accessControlMock struct{}

// NewAccessControl returns mock implementation of access control.
func NewAccessControl() auth.AccessControl {
	return accessControlMock{}
}

func (acm accessControlMock) PotentialPartners(string) ([]string, error) {
	return nil, nil
}

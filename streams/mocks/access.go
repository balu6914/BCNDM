package mocks

import "github.com/datapace/datapace/streams"

var _ streams.AccessControl = (*accessControlMock)(nil)

type accessControlMock struct {
	partners []string
}

// NewAccessControl returns mock access control instance.
func NewAccessControl(partners []string) streams.AccessControl {
	return accessControlMock{partners: partners}
}

func (ac accessControlMock) Partners(partner string) ([]string, error) {
	return ac.partners, nil
}

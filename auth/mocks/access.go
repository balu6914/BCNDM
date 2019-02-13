package mocks

import "datapace/auth"

var _ auth.AccessRequestRepository = (*accessRequestRepoMock)(nil)

type accessRequestRepoMock struct{}

// NewAccessRequestRepository returns mock instance of access request repository.
func NewAccessRequestRepository() auth.AccessRequestRepository {
	return accessRequestRepoMock{}
}

func (repo accessRequestRepoMock) RequestAccess(string, string) (string, error) {
	return "", nil
}

func (repo accessRequestRepoMock) ListSentAccessRequests(string, auth.State) ([]auth.AccessRequest, error) {
	return nil, nil
}

func (repo accessRequestRepoMock) ListReceivedAccessRequests(string, auth.State) ([]auth.AccessRequest, error) {
	return nil, nil
}

func (repo accessRequestRepoMock) ApproveAccessRequest(string, string) error {
	return nil
}

func (repo accessRequestRepoMock) RejectAccessRequest(string, string) error {
	return nil
}

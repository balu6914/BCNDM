package mocks

import "datapace/subscriptions"

var _ subscriptions.TransactionsService = (*transactionsServiceMock)(nil)

type transactionsServiceMock struct {
	from uint64
}

// NewTransactionsService returns mock transactions service instance.
func NewTransactionsService(from uint64) subscriptions.TransactionsService {
	return transactionsServiceMock{from: from}
}

func (svc transactionsServiceMock) Transfer(stream, from, to string, value uint64) error {
	if svc.from < value {
		return subscriptions.ErrFailedTransfer
	}

	return nil
}

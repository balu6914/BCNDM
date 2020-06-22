package mocks

import "github.com/datapace/transactions"

var _ transactions.StreamsService = (*streamsServiceMock)(nil)

type streamsServiceMock struct {
	streams map[string]transactions.Stream
}

// NewStreamsService returns mock streams service instance.
func NewStreamsService(streams map[string]transactions.Stream) transactions.StreamsService {
	return streamsServiceMock{streams: streams}
}

func (svc streamsServiceMock) One(id string) (transactions.Stream, error) {
	stream, ok := svc.streams[id]
	if !ok {
		return transactions.Stream{}, transactions.ErrNotFound
	}

	return stream, nil
}

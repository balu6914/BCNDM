package mocks

import "github.com/datapace/subscriptions"

var _ subscriptions.StreamsService = (*streamsServiceMock)(nil)

type streamsServiceMock struct {
	streams map[string]subscriptions.Stream
}

// NewStreamsService returns mock streams service instance.
func NewStreamsService(streams map[string]subscriptions.Stream) subscriptions.StreamsService {
	return streamsServiceMock{streams: streams}
}

func (svc streamsServiceMock) One(id string) (subscriptions.Stream, error) {
	stream, ok := svc.streams[id]
	if !ok {
		return subscriptions.Stream{}, subscriptions.ErrNotFound
	}

	return stream, nil
}

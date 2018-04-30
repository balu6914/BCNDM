package mocks

import (
	"sync"

	"monetasa/dapp"
)

var _ dapp.StreamRepository = (*streamRepositoryMock)(nil)

const strName = "stream_name_"

type streamRepositoryMock struct {
	mu      sync.Mutex
	counter int
	streams map[string]dapp.Stream
}

// NewStreamRepository creates in-memory stream repository.
func NewStreamRepository() dapp.StreamRepository {
	return &streamRepositoryMock{
		streams: make(map[string]dapp.Stream),
	}
}

func (srm *streamRepositoryMock) Save(stream dapp.Stream) (string, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if _, ok := srm.streams[stream.ID.Hex()]; ok {
		return "", dapp.ErrConflict
	}
	srm.streams[stream.ID.Hex()] = stream

	return stream.ID.Hex(), nil
}

func (srm *streamRepositoryMock) Update(id string, stream dapp.Stream) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if _, ok := srm.streams[id]; !ok {
		return dapp.ErrNotFound
	}

	srm.streams[id] = stream

	return nil
}

func (srm *streamRepositoryMock) One(id string) (dapp.Stream, error) {
	if s, ok := srm.streams[id]; ok {
		return s, nil
	}

	return dapp.Stream{}, dapp.ErrNotFound
}

func (srm *streamRepositoryMock) Search([][]float64) ([]dapp.Stream, error) {
	// TODO: implement a geolocation search mock

	return nil, dapp.ErrNotFound
}

func (srm *streamRepositoryMock) Remove(id string) error {
	delete(srm.streams, id)
	return nil
}

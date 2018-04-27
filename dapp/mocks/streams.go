package mocks

import (
	"fmt"
	"sync"

	"monetasa/dapp"
)

var _ dapp.StreamRepository = (*streamRepositoryMock)(nil)

const strId = "123e4567-e89b-12d3-a456-"

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

func (srm *streamRepositoryMock) Id() string {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	srm.counter += 1
	return fmt.Sprintf("%s%012d", strId, srm.counter)
}

// TODO: Generate a dbKey with func key(email, id string) string in the commons.go

func (srm *streamRepositoryMock) Save(stream dapp.Stream) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	srm.streams[srm.Id()] = stream

	return nil
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
	if c, ok := srm.streams[id]; ok {
		return c, nil
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

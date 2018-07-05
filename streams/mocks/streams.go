package mocks

import (
	"strings"
	"sync"

	"monetasa/streams"
)

var _ streams.StreamRepository = (*streamRepositoryMock)(nil)

type streamRepositoryMock struct {
	mu      sync.Mutex
	counter int
	streams map[string]streams.Stream
}

// NewStreamRepository creates in-memory stream repository.
func NewStreamRepository() streams.StreamRepository {
	return &streamRepositoryMock{
		streams: make(map[string]streams.Stream),
	}
}

func (srm *streamRepositoryMock) Save(stream streams.Stream) (string, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	dbKey := key(stream.Owner, stream.ID.Hex())

	if _, ok := srm.streams[dbKey]; ok {
		return "", streams.ErrConflict
	}

	srm.streams[dbKey] = stream

	return stream.ID.Hex(), nil
}

func (srm *streamRepositoryMock) SaveAll(bulk []streams.Stream) error {
	for _, stream := range bulk {
		if _, err := srm.Save(stream); err != nil {
			return err
		}
	}

	return nil
}

func (srm *streamRepositoryMock) Update(stream streams.Stream) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	dbKey := key(stream.Owner, stream.ID.Hex())

	if _, ok := srm.streams[dbKey]; !ok {
		return streams.ErrNotFound
	}

	srm.streams[dbKey] = stream

	return nil
}

func (srm *streamRepositoryMock) One(id string) (streams.Stream, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	for k, v := range srm.streams {
		if strings.HasSuffix(k, id) {
			return v, nil
		}
	}

	return streams.Stream{}, streams.ErrNotFound
}

func (srm *streamRepositoryMock) Search(coords [][]float64) ([]streams.Stream, error) {
	// Geolocation search mock is not used.
	return nil, streams.ErrNotFound
}

func (srm *streamRepositoryMock) Remove(owner, id string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	dbKey := key(owner, id)

	if _, ok := srm.streams[dbKey]; !ok {
		return streams.ErrNotFound
	}

	delete(srm.streams, dbKey)

	return nil
}

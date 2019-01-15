package mocks

import (
	"datapace/streams"
	"strings"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

var _ streams.StreamRepository = (*streamRepositoryMock)(nil)

type streamRepositoryMock struct {
	mu      sync.Mutex
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

	dbKey := stream.ID.Hex()

	if _, ok := srm.streams[dbKey]; ok {
		return "", streams.ErrConflict
	}

	for _, s := range srm.streams {
		if s.URL == stream.URL {
			return "", streams.ErrConflict
		}
	}

	srm.streams[dbKey] = stream

	return stream.ID.Hex(), nil
}

func (srm *streamRepositoryMock) SaveAll(bulk []streams.Stream) error {

	bulkErr := streams.ErrBulkConflict{
		Message:   "Mock error: unique URL violation.",
		Conflicts: []string{},
	}

	for _, stream := range bulk {
		stream.ID = bson.NewObjectId()
		if _, err := srm.Save(stream); err != nil {
			if _, ok := err.(streams.ErrBulkConflict); ok {
				bulkErr.Conflicts = append(bulkErr.Conflicts, stream.URL)
				continue
			}
			return err
		}
	}

	if len(bulkErr.Conflicts) > 0 {
		return bulkErr
	}

	return nil
}

func (srm *streamRepositoryMock) Update(stream streams.Stream) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	dbKey := stream.ID.Hex()

	if v, ok := srm.streams[dbKey]; !ok || v.Owner != stream.Owner {
		return streams.ErrNotFound
	}

	srm.streams[dbKey] = stream

	return nil
}

func (srm *streamRepositoryMock) One(id string) (streams.Stream, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	for k, v := range srm.streams {
		if k == id {
			return v, nil
		}
	}

	return streams.Stream{}, streams.ErrNotFound
}

func inRange(price uint64, min, max *uint64) bool {
	if min != nil && max != nil {
		return price >= *min && price < *max
	}
	if min != nil {
		return price >= *min
	}
	if max != nil {
		return price < *max
	}
	return true
}

func contains(value, search string) bool {
	if strings.HasPrefix(search, "-") {
		search = search[1:]
		if !strings.HasPrefix(search, "-") {
			return !strings.Contains(value, search)
		}
	}
	return strings.Contains(value, search)
}

func (srm *streamRepositoryMock) Search(query streams.Query) (streams.Page, error) {
	ret := []streams.Stream{}
	for _, stream := range srm.streams {
		if contains(stream.Name, query.Name) && contains(stream.Type, query.StreamType) &&
			inRange(stream.Price, query.MinPrice, query.MaxPrice) {
			if query.Owner == "" {
				ret = append(ret, stream)
				continue
			}
			owner := query.Owner
			if strings.HasPrefix(query.Owner, "-") {
				owner := query.Owner[1:]
				if !strings.HasPrefix(owner, "-") && stream.Owner != owner {
					ret = append(ret, stream)
					continue
				}
			}
			if stream.Owner == owner {
				ret = append(ret, stream)
			}
		}
	}

	start := query.Page * query.Limit
	end := start + query.Limit
	page := streams.Page{
		Total:   uint64(len(ret)),
		Limit:   query.Limit,
		Page:    query.Page,
		Content: []streams.Stream{},
	}

	n := uint64(len(ret))
	if start >= n {
		return page, nil
	}
	if end >= n {
		end = n
	}
	ret = ret[start:end]
	page.Content = ret

	// Geolocation search mock is not used.
	return page, nil
}

func (srm *streamRepositoryMock) Remove(owner, id string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	delete(srm.streams, id)

	return nil
}

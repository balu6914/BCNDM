package mocks

import (
	"strings"
	"sync"
	"time"

	"github.com/datapace/datapace/streams"

	"gopkg.in/mgo.v2/bson"
)

var _ streams.StreamRepository = (*streamRepositoryMock)(nil)
var _ streams.CategoryRepository = (*categoryRepositoryMock)(nil)

type streamRepositoryMock struct {
	mu      sync.Mutex
	streams map[string]streams.Stream
}

type categoryRepositoryMock struct {
	mu sync.Mutex
}

// NewStreamRepository creates in-memory stream repository.
func NewStreamRepository() streams.StreamRepository {
	return &streamRepositoryMock{
		streams: make(map[string]streams.Stream),
	}
}

func NewCategoryRepository() streams.StreamRepository {
	return &streamRepositoryMock{
		streams: make(map[string]streams.Stream),
	}
}
func (srm *streamRepositoryMock) Save(stream streams.Stream) (string, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	dbKey := stream.ID

	if _, ok := srm.streams[dbKey]; ok {
		return "", streams.ErrConflict
	}

	for _, s := range srm.streams {
		if s.URL == stream.URL {
			return "", streams.ErrConflict
		}
	}
	now := time.Now().UTC().Round(time.Hour)
	stream.StartDate = &now
	srm.streams[dbKey] = stream

	return stream.ID, nil
}

func (srm *streamRepositoryMock) SaveAll(bulk []streams.Stream) error {

	bulkErr := streams.ErrBulkConflict{
		Message:   "Mock error: unique URL violation.",
		Conflicts: []string{},
	}

	for _, stream := range bulk {
		stream.ID = bson.NewObjectId().Hex()
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

	dbKey := stream.ID

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
			inRange(stream.Price, query.MinPrice, query.MaxPrice) && metadataMatches(stream.Metadata, query.Metadata) &&
			isVisible(stream, query) {
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

func metadataMatches(md map[string]interface{}, constraint map[string]interface{}) bool {
	for k, v := range constraint {
		mdv, present := md[k]
		if !present {
			return false
		}
		mv, constraintValIsMap := v.(map[string]interface{})
		mmdv, mdValIsMap := mdv.(map[string]interface{})
		if constraintValIsMap {
			if !mdValIsMap {
				return false
			}
			if metadataMatches(mmdv, mv) {
				return false
			} else {
				continue
			}
		}
		if mdv != v {
			return false
		}
	}
	return true
}

func isVisible(stream streams.Stream, query streams.Query) bool {
	switch stream.Visibility {
	case streams.Public:
		return true
	case streams.Protected:
		for _, p := range query.Partners {
			if stream.Owner == p {
				return true
			}
		}
		return query.Shared[stream.ID]
	case streams.Private:
		requester := query.Partners[len(query.Partners)-1]
		return stream.Owner == requester
	}
	return false
}

func (srm *streamRepositoryMock) Remove(owner, id string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	delete(srm.streams, id)

	return nil
}

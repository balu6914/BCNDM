package streams_test

import (
	"fmt"
	"github.com/datapace/datapace/streams/groups"
	"github.com/datapace/datapace/streams/sharing"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"

	"github.com/datapace/datapace/streams"
	"github.com/datapace/datapace/streams/mocks"
)

const (
	nameLen = 8
	typeLen = 4
	descLen = 12
	urlLen  = 6
	maxLong = 180
	maxLat  = 90
	maxPage = uint64(100)
	limit   = uint64(3)
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key         = bson.NewObjectId().Hex()
)

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func stream() streams.Stream {
	return streams.Stream{
		Owner:       bson.NewObjectId().Hex(),
		ID:          bson.NewObjectId().Hex(),
		Name:        randomString(nameLen),
		Type:        randomString(typeLen),
		Description: randomString(descLen),
		Snippet: `{
				"sensor_id": "8746",
				"sensor_type": "DHT22",
				"location": "4409",
				"lat": "50.873",
				"lon": "4.698",
				"timestamp": "2018-03-09T00:02:09",
				"temperature": "5.20"
			}`,
		Price: rand.Uint64(),
		URL:   fmt.Sprintf("http://%s.com", randomString(urlLen)),
		Location: streams.Location{
			Type: "Point",
			Coordinates: [2]float64{
				rand.Float64() * (float64)(rand.Intn(maxLat*2)-maxLat),
				rand.Float64() * (float64)(rand.Intn(maxLong*2)-maxLong),
			},
		},
		Visibility: streams.Public,
	}
}

func newService(partners ...string) streams.Service {
	repo := mocks.NewStreamRepository()
	ac := mocks.NewAccessControl(partners)
	ai := mocks.NewAIService()
	terms := mocks.NewTermsService()
	groupsSvc := groups.NewServiceMock()
	sharingSvc := sharing.NewServiceMock()

	return streams.NewService(repo, ac, ai, terms, groupsSvc, sharingSvc)
}

func pointer(val uint64) *uint64 {
	return &val
}

func TestAddStream(t *testing.T) {
	svc := newService()
	s := stream()

	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "add a new stream",
			stream: s,
			err:    nil},
		{
			desc:   "add an existing stream",
			stream: s,
			err:    streams.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := svc.AddStream(tc.stream)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestAddBulkStreams(t *testing.T) {
	svc := newService()
	all := []streams.Stream{}
	for i := 0; i < 100; i++ {
		all = append(all, stream())
	}

	cases := []struct {
		desc    string
		streams []streams.Stream
		key     string
		err     error
	}{
		{
			desc:    "add 100 streams",
			streams: all,
			err:     nil,
		},
		{
			desc:    "add 0 streams",
			streams: []streams.Stream{},
			err:     nil,
		},
	}

	for _, tc := range cases {
		err := svc.AddBulkStreams(tc.streams)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSearchStreams(t *testing.T) {
	svc := newService()
	size := 0
	streamsOwner := bson.NewObjectId().Hex()
	for i := 0; i < 50; i++ {
		size++
		s := stream()
		s.Owner = streamsOwner
		id, err := svc.AddStream(s)
		s.ID = id
		require.Nil(t, err, "Repo should save streams.")
	}
	// Specify two special Streams to match different
	// types of query and different result sets.
	s1Price := uint64(40)
	s2Price := uint64(50)

	s1 := stream()
	s1.Price = s1Price
	s1.Name = "different name"
	s1.Type = "different type"
	s1.Owner = streamsOwner
	id, _ := svc.AddStream(s1)
	s1.ID = id
	size++

	s2 := stream()
	s2.Price = s2Price
	s2.Owner = bson.NewObjectId().Hex()
	id, _ = svc.AddStream(s2)
	s2.ID = id
	size++

	total := uint64(size)
	lmt := int(limit)

	cases := []struct {
		desc  string
		owner string
		size  int
		query streams.Query
		page  streams.Page
	}{
		{
			desc:  "search with query with only the limit specified",
			owner: streamsOwner,
			size:  lmt,
			query: streams.Query{
				Limit: limit,
			},
			page: streams.Page{
				Limit: limit,
				Total: total,
			},
		},
		{
			desc:  "search reset too big offest to default value silently",
			owner: streamsOwner,
			size:  0,
			query: streams.Query{
				Limit: limit,
				Page:  uint64(total + maxPage),
			},
			page: streams.Page{
				Page:  maxPage,
				Limit: limit,
				Total: total,
			},
		},
		{
			desc:  "search with min price specified",
			owner: streamsOwner,
			size:  lmt,
			// Get all except the one with the price1.
			// Content is caluclated this way because MongoDB
			// pages results from the last inserted entry.
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(s1Price + 1),
			},
			page: streams.Page{
				Limit: limit,
				Total: total - 1,
			},
		},
		{
			desc:  "search with max price specified",
			owner: s2.Owner,
			size:  1,
			query: streams.Query{
				Limit:    limit,
				MaxPrice: pointer(s2Price),
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
		{
			desc:  "search with price range specified",
			owner: streamsOwner,
			size:  2,
			// GTE price1 and LT price2 + 1 (to include price2)
			query: streams.Query{
				Limit:    limit,
				MinPrice: pointer(s1Price),
				MaxPrice: pointer(s2Price + 1),
			},
			page: streams.Page{
				Limit: limit,
				Total: 2,
			},
		},
		{
			desc:  "search by owner",
			owner: s2.Owner,
			size:  1,
			query: streams.Query{
				Limit: limit,
				Owner: s2.Owner,
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
		{
			desc:  "search by name",
			owner: streamsOwner,
			size:  1,
			query: streams.Query{
				Limit: limit,
				Name:  s1.Name,
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
		{
			desc:  "search by type",
			owner: streamsOwner,
			size:  1,
			query: streams.Query{
				Limit:      limit,
				StreamType: s1.Type,
			},
			page: streams.Page{
				Limit: limit,
				Total: 1,
			},
		},
		{
			desc: "search by owner other than provided",
			size: lmt,
			query: streams.Query{
				Limit: limit,
				Owner: fmt.Sprintf("-%s", s2.Owner),
			},
			page: streams.Page{
				Limit: limit,
				Total: total - 1,
			},
		},
		{
			desc:  "search by name other than provided",
			owner: streamsOwner,
			size:  lmt,
			query: streams.Query{
				Limit: limit,
				Name:  fmt.Sprintf("-%s", s1.Name),
			},
			page: streams.Page{
				Limit: limit,
				Total: total - 1,
			},
		},
		{
			desc:  "search by type other than provided",
			owner: streamsOwner,
			size:  lmt,
			query: streams.Query{
				Limit:      limit,
				StreamType: fmt.Sprintf("-%s", s1.Type),
			},
			page: streams.Page{
				Limit: limit,
				Total: total - 1,
			},
		},
	}

	for _, tc := range cases {
		res, err := svc.SearchStreams(tc.owner, tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.page.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.page.Limit, res.Limit))
		assert.Equal(t, tc.page.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.page.Total, res.Total))
		assert.Equal(t, tc.size, len(res.Content), fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.size, len(res.Content)))
		for _, s := range res.Content {
			if tc.owner == s.Owner {
				assert.NotEmpty(t, s.URL, fmt.Sprintf("%s: expected Streams of the owner to have an URL, but got an empty value instead", tc.desc))
				continue
			}
			assert.Empty(t, s.URL, fmt.Sprintf("%s: expected Streams of the other owners not to have an URL, but got a value instead", tc.desc))
		}
	}
}

func TestSearchStreamsShared(t *testing.T) {
	svc := newService()
	stream0 := streams.Stream{
		Owner:      "user0",
		ID:         "stream0",
		Visibility: "protected",
		URL:        "url0",
	}
	_, _ = svc.AddStream(stream0)
	stream1 := streams.Stream{
		Owner:      "user1",
		ID:         "stream1",
		Visibility: "protected",
		URL:        "url1",
	}
	_, _ = svc.AddStream(stream1)
	cases := []struct {
		desc        string
		userId      string
		query       streams.Query
		resultsPage streams.Page
	}{
		{
			desc:   "Search a stream shared to a user",
			userId: "sharingReceiverUser",
			query: streams.Query{
				Limit: 1_000_000,
			},
			resultsPage: streams.Page{
				Page:    0,
				Limit:   1_000_000,
				Total:   1,
				Content: []streams.Stream{stream0},
			},
		},
		{
			desc:   "Search a stream shared to a group including the requested user",
			userId: "userInSomeGroups",
			query: streams.Query{
				Limit: 1_000_000,
			},
			resultsPage: streams.Page{
				Page:    0,
				Limit:   1_000_000,
				Total:   1,
				Content: []streams.Stream{stream1},
			},
		},
		{
			desc:   "Search for streams shared to both a requested user and a group including the requested user",
			userId: "sharingReceiverUserInSomeGroups",
			query: streams.Query{
				Limit: 1_000_000,
			},
			resultsPage: streams.Page{
				Page:  0,
				Limit: 1_000_000,
				Total: 2,
				Content: []streams.Stream{
					stream0,
					stream1,
				},
			},
		},
		{
			desc:   "Search a stream shared to a user but filtered by owner filter",
			userId: "sharingReceiverUser",
			query: streams.Query{
				Limit: 1_000_000,
				Owner: "user1",
			},
			resultsPage: streams.Page{
				Page:    0,
				Limit:   1_000_000,
				Total:   0,
				Content: []streams.Stream{},
			},
		},
	}
	for _, tc := range cases {
		res, err := svc.SearchStreams(tc.userId, tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.resultsPage.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.resultsPage.Limit, res.Limit))
		assert.Equal(t, tc.resultsPage.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.resultsPage.Total, res.Total))
		for _, expectedStream := range tc.resultsPage.Content {
			found := false
			for _, s := range res.Content {
				if expectedStream.ID != s.ID {
					continue
				}
				found = true
			}
			assert.True(t, found)
		}
	}
}

func TestSearchStreamsByMetadata(t *testing.T) {
	svc := newService()
	stream0 := streams.Stream{
		Owner:      "user0",
		ID:         "stream0",
		Visibility: "protected",
		URL:        "url0",
		Metadata: map[string]interface{}{
			"Ends": "2022-05-31T09:54:15Z",
		},
	}
	_, _ = svc.AddStream(stream0)
	stream1 := streams.Stream{
		Owner:      "user0",
		ID:         "stream1",
		Visibility: "protected",
		URL:        "url1",
	}
	_, _ = svc.AddStream(stream1)
	stream2 := streams.Stream{
		Owner:      "user0",
		ID:         "stream2",
		Visibility: "protected",
		URL:        "url2",
		Metadata: map[string]interface{}{
			"Starts": "2022-05-30T09:54:15Z",
			"Ends":   "2022-05-31T09:54:15Z",
		},
	}
	_, _ = svc.AddStream(stream2)
	cases := []struct {
		desc        string
		userId      string
		query       streams.Query
		resultsPage streams.Page
	}{
		{
			desc:   "Search streams w/o metadata constraint",
			userId: "user0",
			query: streams.Query{
				Limit: 1_000_000,
				Owner: "user0",
			},
			resultsPage: streams.Page{
				Page:  0,
				Limit: 1_000_000,
				Total: 3,
				Content: []streams.Stream{
					stream0,
					stream1,
					stream2,
				},
			},
		},
		{
			desc:   "Search streams w/ metadata constraint",
			userId: "user0",
			query: streams.Query{
				Limit: 1_000_000,
				Owner: "user0",
				Metadata: map[string]interface{}{
					"Ends": "2022-05-31T09:54:15Z",
				},
			},
			resultsPage: streams.Page{
				Page:  0,
				Limit: 1_000_000,
				Total: 2,
				Content: []streams.Stream{
					stream0,
					stream2,
				},
			},
		},
		{
			desc:   "Search streams w/ non equal metadata constraint",
			userId: "user0",
			query: streams.Query{
				Limit: 1_000_000,
				Owner: "user0",
				Metadata: map[string]interface{}{
					"foo": "bar",
				},
			},
			resultsPage: streams.Page{
				Page:    0,
				Limit:   1_000_000,
				Total:   0,
				Content: []streams.Stream{},
			},
		},
		{
			desc:   "Search streams w/ multiple metadata constraints",
			userId: "user0",
			query: streams.Query{
				Limit: 1_000_000,
				Owner: "user0",
				Metadata: map[string]interface{}{
					"Starts": "2022-05-30T09:54:15Z",
					"Ends":   "2022-05-31T09:54:15Z",
				},
			},
			resultsPage: streams.Page{
				Page:  0,
				Limit: 1_000_000,
				Total: 1,
				Content: []streams.Stream{
					stream2,
				},
			},
		},
	}
	for _, tc := range cases {
		res, err := svc.SearchStreams(tc.userId, tc.query)
		assert.Nil(t, err, "There should be no error searching streams")
		assert.Equal(t, tc.resultsPage.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.resultsPage.Limit, res.Limit))
		assert.Equal(t, tc.resultsPage.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.resultsPage.Total, res.Total))
		for _, expectedStream := range tc.resultsPage.Content {
			found := false
			for _, s := range res.Content {
				if expectedStream.ID != s.ID {
					continue
				}
				found = true
			}
			assert.True(t, found)
		}
	}
}

func TestUpdateStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	wrongOwner := s
	wrongOwner.Owner = bson.NewObjectId().Hex()

	cases := []struct {
		desc   string
		stream streams.Stream
		err    error
	}{
		{
			desc:   "update an existing stream",
			stream: s,
			err:    nil,
		},
		{
			desc:   "update a non-existing stream",
			stream: stream(),
			err:    streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateStream(tc.stream)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestFullViewStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "view an existing stream",
			id:   s.ID,
			err:  nil,
		},
		{
			desc: "view a non-existing stream",
			id:   bson.NewObjectId().String(),
			err:  streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewFullStream(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStream(t *testing.T) {
	s := stream()
	svc := newService(s.Owner)
	svc.AddStream(s)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "view an existing stream",
			id:   s.ID,
			err:  nil,
		},
		{
			desc: "view a non-existing stream",
			id:   bson.NewObjectId().String(),
			err:  streams.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewStream(tc.id, "")
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewStreamShared(t *testing.T) {
	svc := newService()
	stream0 := streams.Stream{
		Owner:      "user0",
		ID:         "stream0",
		Visibility: "protected",
		URL:        "url0",
	}
	_, _ = svc.AddStream(stream0)
	stream1 := streams.Stream{
		Owner:      "user1",
		ID:         "stream1",
		Visibility: "protected",
		URL:        "url1",
	}
	_, _ = svc.AddStream(stream1)
	cases := []struct {
		desc           string
		streamId       string
		userId         string
		expectedErr    error
		expectedStream *streams.Stream
	}{
		{
			desc:           "View a stream shared to a user",
			streamId:       "stream0",
			userId:         "sharingReceiverUser",
			expectedErr:    nil,
			expectedStream: &stream0,
		},
		{
			desc:           "View a stream not shared to a user",
			streamId:       "stream1",
			userId:         "sharingReceiverUser",
			expectedErr:    streams.ErrNotFound,
			expectedStream: nil,
		},
		{
			desc:           "View a stream not shared to any group including the requested user",
			streamId:       "stream0",
			userId:         "userInSomeGroups",
			expectedErr:    streams.ErrNotFound,
			expectedStream: nil,
		},
		{
			desc:           "View a stream shared to a group including the requested user",
			streamId:       "stream1",
			userId:         "userInSomeGroups",
			expectedErr:    nil,
			expectedStream: &stream1,
		},
	}
	for _, tc := range cases {
		s, err := svc.ViewStream(tc.streamId, tc.userId)
		assert.Equal(t, tc.expectedErr, err)
		if tc.expectedStream != nil {
			assert.Equal(t, tc.expectedStream.ID, s.ID)
		}
	}
}

func TestRemoveStream(t *testing.T) {
	s := stream()
	svc := newService()
	svc.AddStream(s)

	cases := []struct {
		desc     string
		streamId string
		owner    string
		err      error
	}{
		{
			desc:     "remove an existing stream",
			streamId: s.ID,
			owner:    s.Owner,
			err:      nil,
		},
		{
			desc:     "remove a non-existing stream",
			streamId: bson.NewObjectId().Hex(),
			owner:    s.Owner,
			err:      nil,
		},
	}

	for _, tc := range cases {
		err := svc.RemoveStream(tc.owner, tc.streamId)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestExportStreamsShared(t *testing.T) {
	svc := newService()
	stream0 := streams.Stream{
		Owner:      "user0",
		ID:         "stream0",
		Visibility: "protected",
		URL:        "url0",
	}
	_, _ = svc.AddStream(stream0)
	stream1 := streams.Stream{
		Owner:      "user1",
		ID:         "stream1",
		Visibility: "protected",
		URL:        "url1",
	}
	_, _ = svc.AddStream(stream1)
	cases := []struct {
		desc            string
		userId          string
		expectedStreams []streams.Stream
	}{
		{
			desc:            "Export a stream shared to a user",
			userId:          "sharingReceiverUser",
			expectedStreams: []streams.Stream{stream0},
		},
		{
			desc:            "Export a stream shared to a group including the requested user",
			userId:          "userInSomeGroups",
			expectedStreams: []streams.Stream{stream1},
		},
		{
			desc:   "Export multiple streams shared to both a requested user and a group including the requested user",
			userId: "sharingReceiverUserInSomeGroups",
			expectedStreams: []streams.Stream{
				stream0,
				stream1,
			},
		},
	}
	for _, tc := range cases {
		ss, err := svc.ExportStreams(tc.userId)
		assert.Nil(t, err, "There should be no error exporting streams")
		for _, expectedStream := range tc.expectedStreams {
			found := false
			for _, s := range ss {
				if expectedStream.ID != s.ID {
					continue
				}
				found = true
			}
			assert.True(t, found)
		}
	}
}

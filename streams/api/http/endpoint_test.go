package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"monetasa/streams"
	httpapi "monetasa/streams/api/http"
	"monetasa/streams/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const (
	validFilePath    = "../../../assets/test/validBulk.csv"
	invalidFilePath  = "../../../assets/test/invalidBulk.csv"
	conflictFilePath = "../../../assets/test/conflictBulk.csv"
	urlLen           = 20
)

var (
	validKey    = bson.NewObjectId().Hex()
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	counter     = rand.Intn(100)
)

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func genStream() streams.Stream {
	counter++
	return streams.Stream{
		ID:          bson.NewObjectId(),
		Name:        "name",
		Type:        "type",
		Description: "description",
		Snippet: `{
			"sensor_id": "8746",
			"sensor_type": "DHT22",
			"location": "4409",
			"lat": "50.873",
			"lon": "4.698",
			"timestamp": "2018-03-09T00:02:09",
			"temperature": "5.20"
		}`,
		Price: 123,
		URL:   fmt.Sprintf("https://myStream%d.com", counter),
		Owner: validKey,
		Location: streams.Location{
			Type:        "Point",
			Coordinates: [2]float64{50, 50},
		},
	}
}

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	token       string
	contentType string
	body        io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}

	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}

	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}

	return tr.client.Do(req)
}

func newService() streams.Service {
	repo := mocks.NewStreamRepository()
	return streams.NewService(repo)
}

func newServer(svc streams.Service) *httptest.Server {
	auth := mocks.NewAuth([]string{validKey})
	mux := httpapi.MakeHandler(svc, auth)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func sendFile(fileName, targetURL string) (io.Reader, string) {
	bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)
	fileWriter, _ := bodyWriter.CreateFormFile("csv", fileName)
	f, _ := os.Open(fileName)
	defer f.Close()

	io.Copy(fileWriter, f)
	ct := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return bytes.NewReader(bodyBuff.Bytes()), ct
}

func makeQuery(page, limit uint, name, streamType, owner string, minPrice, maxPrice *uint64) string {
	ret := fmt.Sprintf("?page=%d&limit=%d", page, limit)
	if owner != "" {
		ret = fmt.Sprintf("%s&owner=%s", ret, owner)
	}
	if name != "" {
		ret = fmt.Sprintf("%s&name=%s", ret, name)
	}
	if streamType != "" {
		ret = fmt.Sprintf("%s&type=%s", ret, streamType)
	}
	if minPrice != nil {
		ret = fmt.Sprintf("%s&minPrice=%d", ret, *minPrice)
	}
	if maxPrice != nil {
		ret = fmt.Sprintf("%s&maxPrice=%d", ret, *maxPrice)
	}

	return ret
}

func TestAddStream(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	stream := genStream()
	valid := toJSON(stream)

	s := genStream()
	s.Name = ""
	invalid := toJSON(s)

	cases := []struct {
		desc   string
		req    string
		auth   string
		status int
	}{
		{
			desc:   "add a valid stream",
			req:    valid,
			auth:   validKey,
			status: http.StatusCreated,
		},
		{
			desc:   "add an ivalid stream",
			req:    invalid,
			auth:   validKey,
			status: http.StatusBadRequest,
		},
		{
			desc:   "add an empty stream",
			req:    "{}",
			auth:   validKey,
			status: http.StatusBadRequest,
		},
		{
			desc:   "add a stream with invalid token",
			req:    valid,
			auth:   "invalid key",
			status: http.StatusForbidden,
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodPost,
			url:    fmt.Sprintf("%s/streams", ts.URL),
			token:  tc.auth,
			body:   strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestAddBulkStreams(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()
	valid, ct := sendFile(validFilePath, ts.URL)
	invalid, ct1 := sendFile(invalidFilePath, ts.URL)
	conflict, ct2 := sendFile(conflictFilePath, ts.URL)

	cases := []struct {
		desc        string
		req         io.Reader
		contentType string
		auth        string
		status      int
	}{
		{
			desc:        "add a valid bulk of streams",
			req:         valid,
			contentType: ct,
			auth:        validKey,
			status:      http.StatusCreated,
		},
		{
			desc:        "add a bulk of streams unauthorized",
			req:         valid,
			contentType: ct,
			auth:        "unauthorized",
			status:      http.StatusForbidden,
		},
		{
			desc:        "add a valid bulk of streams with wrong contentType",
			req:         valid,
			contentType: "json",
			auth:        validKey,
			status:      http.StatusUnsupportedMediaType,
		},
		{
			desc:        "add an invalid bulk of streams",
			req:         invalid,
			contentType: ct1,
			auth:        validKey,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "add a bulk of streams with conflicts",
			req:         conflict,
			contentType: ct2,
			auth:        validKey,
			status:      http.StatusConflict,
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/streams/bulk", ts.URL),
			token:       tc.auth,
			contentType: tc.contentType,
			body:        tc.req,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestSearchStreams(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	total := uint64(200)
	for i := uint64(0); i < total; i++ {
		svc.AddStream(genStream())
	}

	// Specify two special Streams to match different
	// types of query and different result sets.
	price1 := uint64(40)
	price2 := uint64(50)

	s := genStream()
	s.ID = bson.NewObjectId()
	s.Price = price1
	svc.AddStream(s)

	s = genStream()
	s.ID = bson.NewObjectId()
	s.Price = price2
	s.Owner = bson.NewObjectId().Hex()
	s.Name = "special_name"
	s.Type = "special_type"
	svc.AddStream(s)
	// Add special streams to count.
	total += 2

	cases := []struct {
		desc   string
		auth   string
		status int
		query  string
		size   int
		res    streams.Page
	}{
		{
			desc:   "search streams with no query provided",
			auth:   validKey,
			status: http.StatusOK,
			query:  "",
			size:   20,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: total,
			},
		},
		{
			desc:   "search streams unauthorized",
			auth:   "invalid key",
			status: http.StatusForbidden,
			query:  "",
		},
		{
			desc:   "search streams by page and limit only",
			auth:   validKey,
			status: http.StatusOK,
			query:  makeQuery(3, 30, "", "", "", nil, nil),
			size:   30,
			res: streams.Page{
				Page:  3,
				Limit: 30,
				Total: total,
			},
		},
		{
			desc:   "search streams by the owner",
			auth:   validKey,
			status: http.StatusOK,
			query:  makeQuery(0, 20, "", "", s.Owner, nil, nil),
			size:   1,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: 1,
			},
		},
		{
			desc:   "search streams by page, limit and price",
			auth:   validKey,
			status: http.StatusOK,
			query:  makeQuery(3, 30, "", "", "", &price2, nil),
			size:   30,
			res: streams.Page{
				Page:  3,
				Limit: 30,
				Total: total - 1,
			},
		},
		{
			desc:   "search streams by page, limit and price range",
			auth:   validKey,
			status: http.StatusOK,
			query:  makeQuery(0, 30, "", "", "", &price1, &price2),
			size:   1,
			res: streams.Page{
				Page:  0,
				Limit: 30,
				Total: 1,
			},
		},
		{
			desc:   "search streams by page, limit and price with too big page and limit",
			auth:   validKey,
			status: http.StatusOK,
			query:  makeQuery(3, 30, "", "", "", &price1, &price2),
			size:   0,
			res: streams.Page{
				Page:  3,
				Limit: 30,
				Total: 1,
			},
		},
		{
			desc:   "search streams by owner",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?owner=%s", s.Owner),
			size:   1,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: 1,
			},
		},
		{
			desc:   "search streams by name",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?name=%s", s.Name[0:5]),
			size:   1,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: 1,
			},
		},
		{
			desc:   "search streams by type",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?type=%s", s.Type[0:5]),
			size:   1,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: 1,
			},
		},
		{
			desc:   "search streams by owner other than provided",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?owner=-%s", s.Owner),
			size:   20,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: total - 1,
			},
		},
		{
			desc:   "search streams by name other than provided",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?name=-%s", s.Name[0:5]),
			size:   20,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: total - 1,
			},
		},
		{
			desc:   "search streams by type other than provided",
			auth:   validKey,
			status: http.StatusOK,
			query:  fmt.Sprintf("?type=-%s", s.Type[0:5]),
			size:   20,
			res: streams.Page{
				Page:  0,
				Limit: 20,
				Total: total - 1,
			},
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/streams%s", ts.URL, tc.query),
			token:  tc.auth,
		}
		r, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		defer r.Body.Close()
		// Unauthorized requests should not be processed.
		if tc.auth != validKey {
			assert.Equal(t, tc.status, r.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, r.StatusCode))
			continue
		}
		res := streams.Page{}
		err = json.NewDecoder(r.Body).Decode(&res)

		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.res.Limit, res.Limit, fmt.Sprintf("%s: expected limit %d got %d\n", tc.desc, tc.res.Limit, res.Limit))
		assert.Equal(t, tc.res.Total, res.Total, fmt.Sprintf("%s: expected total %d got %d\n", tc.desc, tc.res.Total, res.Total))
		// Don't use actual content, only compare expected size.
		assert.Equal(t, tc.size, len(res.Content), fmt.Sprintf("%s: expected size of batch %d got %d\n", tc.desc, tc.size, len(res.Content)))
		assert.Equal(t, tc.status, r.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, r.StatusCode))
	}
}

func TestUpdateStream(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	stream := genStream()
	svc.AddStream(stream)
	valid := toJSON(stream)
	// Create an invalid stream.
	s := genStream()
	s.Name = ""
	invalid := toJSON(s)
	// Create stream that does not exist in database.
	s = stream
	nonExistingId := bson.NewObjectId()
	s.ID = nonExistingId
	nonExisting := toJSON(s)

	cases := []struct {
		desc   string
		req    string
		auth   string
		status int
		id     string
	}{
		{
			desc:   "update an existing stream",
			req:    valid,
			auth:   validKey,
			status: http.StatusOK,
			id:     stream.ID.Hex(),
		},
		{
			desc:   "update a stream with non-matching stream ID and URL id",
			req:    valid,
			auth:   validKey,
			status: http.StatusBadRequest,
			id:     bson.NewObjectId().Hex(),
		},
		{
			desc:   "update a non-existing stream",
			req:    nonExisting,
			auth:   validKey,
			status: http.StatusNotFound,
			id:     nonExistingId.Hex(),
		},
		{
			desc:   "update a stream with invalid data",
			req:    invalid,
			auth:   validKey,
			status: http.StatusBadRequest,
			id:     stream.ID.Hex(),
		},
		{
			desc:   "update a stream without an auth key",
			req:    valid,
			auth:   "",
			status: http.StatusForbidden,
			id:     stream.ID.Hex(),
		},
		{
			desc:   "update a stream with invalid auth key",
			req:    valid,
			auth:   "invalid",
			status: http.StatusForbidden,
			id:     stream.ID.Hex(),
		},
		{
			desc:   "update stream with an empty request",
			req:    "",
			auth:   validKey,
			status: http.StatusBadRequest,
			id:     stream.ID.Hex(),
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodPut,
			url:    fmt.Sprintf("%s/streams/%s", ts.URL, tc.id),
			token:  tc.auth,
			body:   strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestViewStream(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	stream := genStream()
	svc.AddStream(stream)

	cases := []struct {
		desc   string
		id     string
		auth   string
		status int
	}{
		{
			desc:   "get a stream with no errors",
			auth:   validKey,
			id:     stream.ID.Hex(),
			status: http.StatusOK,
		},
		{
			desc:   "gat a stream with no auth key",
			auth:   "",
			id:     stream.ID.Hex(),
			status: http.StatusForbidden,
		},
		{
			desc:   "get a nonexisting stream",
			auth:   validKey,
			id:     bson.NewObjectId().Hex(),
			status: http.StatusNotFound,
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/streams/%s", ts.URL, tc.id),
			token:  tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestRemoveStream(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	stream := genStream()
	svc.AddStream(stream)

	cases := []struct {
		desc   string
		id     string
		auth   string
		status int
	}{
		{
			desc:   "remove an existing stream",
			auth:   validKey,
			id:     stream.ID.Hex(),
			status: http.StatusNoContent,
		},
		{
			desc:   "remove a stream with no auth key",
			auth:   "",
			id:     stream.ID.Hex(),
			status: http.StatusForbidden,
		},
		{
			desc:   "remove a nonexisting stream",
			auth:   validKey,
			id:     bson.NewObjectId().Hex(),
			status: http.StatusNoContent,
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodDelete,
			url:    fmt.Sprintf("%s/streams/%s", ts.URL, tc.id),
			token:  tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

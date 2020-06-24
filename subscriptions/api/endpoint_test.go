package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/subscriptions"
	httpapi "github.com/datapace/datapace/subscriptions/api"
	"github.com/datapace/datapace/subscriptions/mocks"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

const (
	wrong       = "wrong"
	contentType = "application/json"
	streamURL   = "myUrl"
	hours       = 24
	balance     = 100000
)

var (
	token    = bson.NewObjectId().Hex()
	userID   = bson.NewObjectId().Hex()
	streamID = bson.NewObjectId().Hex()

	sub = subscriptions.Subscription{
		ID:        bson.NewObjectId(),
		UserID:    userID,
		StreamID:  streamID,
		StreamURL: streamURL,
		Hours:     hours,
	}
	subs = []subscriptions.Subscription{sub}
)

func newService() subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	streams := mocks.NewStreamsService(map[string]subscriptions.Stream{
		streamID: subscriptions.Stream{},
	})
	proxy := mocks.NewProxy()
	transactions := mocks.NewTransactionsService(balance)
	auth := mocks.NewAuthClient(nil, nil)
	return subscriptions.New(auth, subs, streams, proxy, transactions)
}

func newServer(svc subscriptions.Service, ac datapace.AuthServiceClient) *httptest.Server {
	mux := httpapi.MakeHandler(svc, ac)
	return httptest.NewServer(mux)
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
	jm, _ := json.Marshal(subs)
	msg := bytes.NewBuffer([]byte(jm))

	req, err := http.NewRequest(tr.method, tr.url, msg)
	if err != nil {
		return nil, err
	}

	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}

	return tr.client.Do(req)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateSubscription(t *testing.T) {
	ac := mocks.NewAuthClient(map[string]string{
		token: userID,
	}, nil)
	svc := newService()
	ss := newServer(svc, ac)
	defer ss.Close()

	body := toJSON(subs)

	cases := []struct {
		desc        string
		auth        string
		contentType string
		body        string
		status      int
	}{
		{
			desc:        "create subscriptions with valid credentials",
			auth:        token,
			contentType: contentType,
			body:        body,
			status:      http.StatusMultiStatus,
		},
		{
			desc:        "create subscriptions with invalid credentials",
			auth:        wrong,
			contentType: contentType,
			body:        body,
			status:      http.StatusForbidden,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: ss.Client(),
			method: http.MethodPost,
			url:    fmt.Sprintf("%s/subscriptions", ss.URL),
			token:  tc.auth,
			body:   strings.NewReader(tc.body),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		contentType := res.Header.Get("Content-Type")
		assert.Equal(t, tc.contentType, contentType, fmt.Sprintf("%s: expected content type %s got %s", tc.desc, tc.contentType, contentType))
	}
}

func TestViewSubscription(t *testing.T) {
	ac := mocks.NewAuthClient(map[string]string{
		token: userID,
	}, nil)

	svc := newService()
	_, err := svc.AddSubscription(userID, "", sub)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ss := newServer(svc, ac)
	defer ss.Close()

	cases := []struct {
		desc        string
		auth        string
		contentType string
		id          string
		status      int
	}{
		{
			desc:        "get subscriptions with valid credentials",
			auth:        token,
			contentType: contentType,
			id:          sub.ID.Hex(),
			status:      http.StatusOK,
		},
		{
			desc:        "get subscriptions with invalid credentials",
			auth:        wrong,
			contentType: contentType,
			id:          sub.ID.Hex(),
			status:      http.StatusForbidden,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: ss.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/subscriptions/%s", ss.URL, tc.id),
			token:  tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		contentType := res.Header.Get("Content-Type")
		assert.Equal(t, tc.contentType, contentType, fmt.Sprintf("%s: expected content type %s got %s", tc.desc, tc.contentType, contentType))
	}
}

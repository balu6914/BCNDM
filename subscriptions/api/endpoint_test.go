package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"monetasa"
	"monetasa/subscriptions"
	httpapi "monetasa/subscriptions/api"
	"monetasa/subscriptions/mocks"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	token       = "token"
	wrong       = "wrong"
	userID      = "userID"
	contentType = "application/json"
	streamID    = "streamID"
	streamURL   = "myUrl"
	hours       = 24
)

var sub = subscriptions.Subscription{
	UserID:    userID,
	StreamID:  streamID,
	StreamURL: streamURL,
	Hours:     hours,
}

func newService() subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	return subscriptions.New(subs)
}

func newServer(svc subscriptions.Service, ac monetasa.AuthServiceClient) *httptest.Server {
	mux := httpapi.MakeHandler(svc, ac)
	return httptest.NewServer(mux)
}

type testRequest struct {
	client *http.Client
	method string
	url    string
	token  string
}

func (tr testRequest) make() (*http.Response, error) {
	jm, _ := json.Marshal(sub)
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

func TestCreateSubscriptions(t *testing.T) {
	ac := mocks.NewAuthClient(map[string]string{
		token: userID,
	})
	svc := newService()
	ss := newServer(svc, ac)
	defer ss.Close()

	cases := map[string]struct {
		auth        string
		contentType string
		method      string
		status      int
	}{
		"create subscriptions with valid credentials": {
			auth:        token,
			contentType: contentType,
			method:      http.MethodPost,
			status:      http.StatusCreated,
		},
		"create subscriptions with invalid credentials": {
			auth:        wrong,
			contentType: contentType,
			method:      http.MethodPost,
			status:      http.StatusForbidden,
		},
		"get subscriptions with valid credentials": {
			auth:        token,
			contentType: contentType,
			method:      http.MethodGet,
			status:      http.StatusOK,
		},
		"get subscriptions with invalid credentials": {
			auth:        wrong,
			contentType: contentType,
			method:      http.MethodGet,
			status:      http.StatusForbidden,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client: ss.Client(),
			method: tc.method,
			url:    fmt.Sprintf("%s/subscriptions", ss.URL),
			token:  tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		contentType := res.Header.Get("Content-Type")
		assert.Equal(t, tc.contentType, contentType, fmt.Sprintf("%s: expected content type %s got %s", desc, tc.contentType, contentType))
	}
}

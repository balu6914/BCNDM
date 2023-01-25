package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/datapace/datapace/subscriptions/accessv2"
	"github.com/datapace/datapace/subscriptions/sharing"
	"io"
	"io/ioutil"
	"strings"

	authproto "github.com/datapace/datapace/proto/auth"
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
	token          = bson.NewObjectId().Hex()
	userID         = bson.NewObjectId().Hex()
	streamID       = bson.NewObjectId().Hex()
	nonOwnStreamId = bson.NewObjectId().Hex()

	sub = subscriptions.Subscription{
		ID:        bson.NewObjectId(),
		UserID:    userID,
		StreamID:  streamID,
		StreamURL: streamURL,
		Hours:     hours,
	}
	subs        = []subscriptions.Subscription{sub}
	nonOwnerSub = subscriptions.Subscription{
		ID:        bson.NewObjectId(),
		StreamID:  nonOwnStreamId,
		StreamURL: streamURL,
		Hours:     hours,
	}
	nonOwnerSubs = []subscriptions.Subscription{nonOwnerSub}
)

func newService() subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	streams := mocks.NewStreamsService(map[string]subscriptions.Stream{
		streamID: {
			Owner:      userID,
			Visibility: "private",
		},
		nonOwnStreamId: {
			Owner:      bson.NewObjectId().Hex(),
			Visibility: "private",
		},
	})
	proxy := mocks.NewProxy()
	transactions := mocks.NewTransactionsService(balance)
	auth := mocks.NewAuthClient(nil, nil)
	sharingSvc := sharing.NewServiceMock()
	accessV2Svc := accessv2.NewServiceMock()
	return subscriptions.New(auth, subs, streams, proxy, transactions, sharingSvc, accessV2Svc)
}

func newServer(svc subscriptions.Service, ac authproto.AuthServiceClient) *httptest.Server {
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
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
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
	nonOwnerBody := toJSON(nonOwnerSubs)

	cases := []struct {
		desc             string
		auth             string
		contentType      string
		body             string
		status           int
		respBodyContains string
	}{
		{
			desc:        "create subscriptions with valid credentials",
			auth:        token,
			contentType: contentType,
			body:        body,
			status:      http.StatusCreated,
		},
		{
			desc:        "create subscriptions with invalid credentials",
			auth:        wrong,
			contentType: contentType,
			body:        body,
			status:      http.StatusForbidden,
		},
		{
			desc:             "create non-own private stream subscription fails",
			auth:             token,
			contentType:      contentType,
			body:             nonOwnerBody,
			status:           http.StatusCreated,
			respBodyContains: subscriptions.ErrFailedCreateSub.Error(),
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
		if tc.respBodyContains != "" {
			respBody, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Contains(t, string(respBody), tc.respBodyContains)
		}
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

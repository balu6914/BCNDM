package http_test

import (
	"encoding/json"
	"fmt"
	"io"
	"monetasa"
	"monetasa/transactions"
	httpapi "monetasa/transactions/api/http"
	"monetasa/transactions/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	token       = "token"
	userID      = "user"
	chanID      = "chan"
	balance     = 42
	contentType = "application/json"
)

func newService() transactions.Service {
	bn := mocks.NewBlockchainNetwork(map[string]uint64{
		userID: balance,
	})

	return transactions.New(bn)
}

func newServer(svc transactions.Service, ac monetasa.AuthServiceClient) *httptest.Server {
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
	req, err := http.NewRequest(tr.method, tr.url, nil)
	if err != nil {
		return nil, err
	}

	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}

	return tr.client.Do(req)
}

type balanceRes struct {
	Balance uint64 `json:"balance"`
}

func getBalance(r io.Reader) uint64 {
	var res balanceRes
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return 0
	}

	return res.Balance
}

func TestBalance(t *testing.T) {
	ac := mocks.NewAuthClient(map[string]string{
		token: userID,
	})
	svc := newService()
	ts := newServer(svc, ac)
	defer ts.Close()

	cases := map[string]struct {
		auth        string
		chanID      string
		contentType string
		status      int
		balance     uint64
	}{
		"get balance of existing user": {
			auth:        token,
			chanID:      chanID,
			contentType: contentType,
			status:      http.StatusOK,
			balance:     balance,
		},
		"get balance of nonexistent user": {
			auth:        "invalid-token",
			chanID:      chanID,
			contentType: contentType,
			status:      http.StatusForbidden,
			balance:     0,
		},
		"get balance with empty token": {
			auth:        "",
			chanID:      chanID,
			contentType: contentType,
			status:      http.StatusForbidden,
			balance:     0,
		},
		"get balance with invalid channel id": {
			auth:        token,
			chanID:      "",
			contentType: contentType,
			status:      http.StatusBadRequest,
			balance:     0,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/channels/%s/balance", ts.URL, tc.chanID),
			token:  tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		contentType := res.Header.Get("Content-Type")
		assert.Equal(t, tc.contentType, contentType, fmt.Sprintf("%s: expected content type %s got %s", desc, tc.contentType, contentType))
		balance := getBalance(res.Body)
		assert.Equal(t, tc.balance, balance, fmt.Sprintf("%s: expected balance %d got %d", desc, tc.balance, balance))
	}
}

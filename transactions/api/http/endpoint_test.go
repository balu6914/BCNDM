package http_test

import (
	"encoding/json"
	"fmt"
	"io"
	"monetasa/transactions"
	httpapi "monetasa/transactions/api/http"
	"monetasa/transactions/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	token           = "token"
	userID          = "user"
	balance         = 42
	remainingTokens = 100
	contentType     = "application/json"
	invalid         = "invalid"
)

func newService() transactions.Service {
	repo := mocks.NewUserRepository(map[string]string{
		userID: token,
	})
	bn := mocks.NewBlockchainNetwork(map[string]uint64{
		userID: balance,
	}, remainingTokens)

	return transactions.New(repo, bn)
}

func newServer(svc transactions.Service) *httptest.Server {
	ac := mocks.NewAuthClient(map[string]string{
		token: userID,
	})
	mux := httpapi.MakeHandler(svc, ac)
	return httptest.NewServer(mux)
}

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
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

type buyReq struct {
	Amount uint64 `json:"amount"`
}

func TestBalance(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	cases := map[string]struct {
		auth        string
		contentType string
		status      int
		balance     uint64
	}{
		"get balance of existing user": {
			auth:        token,
			contentType: contentType,
			status:      http.StatusOK,
			balance:     balance,
		},
		"get balance of nonexistent user": {
			auth:        invalid,
			contentType: contentType,
			status:      http.StatusForbidden,
			balance:     0,
		},
		"get balance with empty token": {
			auth:        "",
			contentType: contentType,
			status:      http.StatusForbidden,
			balance:     0,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/tokens", ts.URL),
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

func TestBuyTokens(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	req := buyReq{Amount: remainingTokens}
	byteReq, err := json.Marshal(req)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	data := string(byteReq)

	cases := []struct {
		desc        string
		auth        string
		contentType string
		body        string
		status      int
	}{
		{
			desc:        "buy tokens for existing user",
			auth:        token,
			contentType: contentType,
			body:        data,
			status:      http.StatusOK,
		},
		{
			desc:        "buy tokens for nonexistent user",
			auth:        invalid,
			contentType: contentType,
			body:        data,
			status:      http.StatusForbidden,
		},
		{
			desc:        "buy zero tokens for user",
			auth:        token,
			contentType: contentType,
			body:        `{"amount":0}`,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "buy tokens with empty auth token",
			auth:        "",
			contentType: contentType,
			body:        data,
			status:      http.StatusForbidden,
		},
		{
			desc:        "buy tokens with invalid request",
			auth:        token,
			contentType: contentType,
			body:        "}",
			status:      http.StatusBadRequest,
		},
		{
			desc:        "buy tokens with empty request",
			auth:        token,
			contentType: contentType,
			body:        "",
			status:      http.StatusBadRequest,
		},
		{
			desc:        "buy tokens with invalid content type",
			auth:        token,
			contentType: invalid,
			body:        data,
			status:      http.StatusUnsupportedMediaType,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/tokens/buy", ts.URL),
			contentType: tc.contentType,
			token:       tc.auth,
			body:        strings.NewReader(tc.body),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

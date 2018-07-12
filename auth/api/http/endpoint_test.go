package http_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"monetasa/auth"
	httpapi "monetasa/auth/api/http"
	"monetasa/auth/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	contentType = "application/json"
	token       = "token"
	email       = "john.doe@email.com"
	invalid     = "invalid"
)

var user = auth.User{
	ID:       email,
	Email:    email,
	Password: "pass",
}

func newService() auth.Service {
	repo := mocks.NewUserRepository()
	hasher := mocks.NewHasher()
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()

	return auth.New(repo, hasher, idp, ts)
}

func newServer(svc auth.Service) *httptest.Server {
	mux := httpapi.MakeHandler(svc)
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

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRegister(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	data := toJSON(user)
	invalidData := toJSON(auth.User{})
	invalidEmailData := toJSON(auth.User{
		Email:    invalid,
		Password: "pass",
	})

	cases := []struct {
		desc        string
		contentType string
		req         string
		status      int
	}{
		{
			desc:        "register new user",
			contentType: contentType,
			req:         data,
			status:      http.StatusCreated,
		},
		{
			desc:        "register existing user",
			contentType: contentType,
			req:         data,
			status:      http.StatusConflict,
		},
		{
			desc:        "register invalid user",
			contentType: contentType,
			req:         invalidData,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "register user with invalid email",
			contentType: contentType,
			req:         invalidEmailData,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "register user with missing content type",
			contentType: "",
			req:         data,
			status:      http.StatusUnsupportedMediaType,
		},
		{
			desc:        "register user with invalid request format",
			contentType: contentType,
			req:         "}",
			status:      http.StatusBadRequest,
		},
		{
			desc:        "register user with empty JSON request",
			contentType: contentType,
			req:         "{}",
			status:      http.StatusBadRequest,
		},
		{
			desc:        "register user with empty request",
			contentType: contentType,
			req:         "",
			status:      http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestLogin(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	svc.Register(user)

	credentials := user
	credentials.ID = ""
	data := toJSON(credentials)
	tokenData := toJSON(map[string]string{"token": user.Email})

	nonexistentUser := auth.User{
		Email:    "nonexistent.user@email.com",
		Password: "pass",
	}
	nonexistentData := toJSON(nonexistentUser)

	invalidData := toJSON(auth.User{})
	invalidEmailData := toJSON(auth.User{
		Email:    invalid,
		Password: "pass",
	})

	cases := map[string]struct {
		contentType string
		req         string
		status      int
		res         string
	}{
		"login existing user": {
			contentType: contentType,
			req:         data,
			status:      http.StatusCreated,
			res:         tokenData,
		},
		"login non-existent user": {
			contentType: contentType,
			req:         nonexistentData,
			status:      http.StatusForbidden,
			res:         "",
		},
		"login user with invalid data": {
			contentType: contentType,
			req:         invalidData,
			status:      http.StatusBadRequest,
			res:         "",
		},
		"login user with invalid email": {
			contentType: contentType,
			req:         invalidEmailData,
			status:      http.StatusBadRequest,
			res:         "",
		},
		"login user with empty content type": {
			contentType: "",
			req:         data,
			status:      http.StatusUnsupportedMediaType,
			res:         "",
		},
		"login user with invalid request format": {
			contentType: contentType,
			req:         "}",
			status:      http.StatusBadRequest,
			res:         "",
		},
		"login user with empty JSON request": {
			contentType: contentType,
			req:         "{}",
			status:      http.StatusBadRequest,
			res:         "",
		},
		"login user with empty request": {
			contentType: contentType,
			req:         "",
			status:      http.StatusBadRequest,
			res:         "",
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/tokens", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpteds error %s", desc, err))
		token := strings.Trim(string(body), "\n")

		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, token, fmt.Sprintf("%s: expected body %s got %s", desc, tc.res, token))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	svc.Register(user)
	key, _ := svc.Login(user)

	updatedUser := user
	updatedUser.Password = "new_pass"
	data := toJSON(updatedUser)

	invalidData := toJSON(auth.User{})
	invalidEmailData := toJSON(auth.User{
		Email:    invalid,
		Password: "pass",
	})

	cases := []struct {
		desc        string
		contentType string
		req         string
		token       string
		status      int
	}{
		{
			desc:        "update existing user",
			contentType: contentType,
			req:         data,
			token:       key,
			status:      http.StatusOK,
		},
		{
			desc:        "update non-existent user",
			contentType: contentType,
			req:         data,
			token:       "non-existent",
			status:      http.StatusNotFound,
		},
		{
			desc:        "update user with invalid data",
			contentType: contentType,
			req:         invalidData,
			token:       key,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "update user with invalid email",
			contentType: contentType,
			req:         invalidEmailData,
			token:       key,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "update user with empty token",
			contentType: contentType,
			req:         data,
			token:       "",
			status:      http.StatusForbidden,
		},
		{
			desc:        "update user with empty content type",
			contentType: "",
			req:         data,
			token:       key,
			status:      http.StatusUnsupportedMediaType,
		},
		{
			desc:        "update user with invalid request format",
			contentType: contentType,
			req:         "}",
			token:       key,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "update user with empty JSON request",
			contentType: contentType,
			req:         "{}",
			token:       key,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "update user with empty request",
			contentType: contentType,
			req:         "",
			token:       key,
			status:      http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPut,
			url:         fmt.Sprintf("%s/users", ts.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}

}

func TestView(t *testing.T) {
	svc := newService()
	ts := newServer(svc)
	defer ts.Close()

	svc.Register(user)
	key, _ := svc.Login(user)

	cases := map[string]struct {
		token  string
		status int
	}{
		"view existing user": {
			token:  key,
			status: http.StatusOK,
		},
		"view user with invalid token": {
			token:  invalid,
			status: http.StatusForbidden,
		},
		"view user with empty token": {
			token:  "",
			status: http.StatusForbidden,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/users", ts.URL),
			token:  tc.token,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
	}
}
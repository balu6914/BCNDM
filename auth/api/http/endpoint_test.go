package http_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/datapace/datapace/auth"
	httpapi "github.com/datapace/datapace/auth/api/http"
	"github.com/datapace/datapace/auth/mocks"

	"github.com/stretchr/testify/assert"
)

const (
	contentType  = "application/json"
	token        = "token"
	email        = "john.doe@email.com"
	contactEmail = "john.doe@email.com"
	invalid      = "invalid"
)

var policies = map[string]auth.Policy{
	"admin": {
		Name:    "admin",
		Owner:   "admin",
		Version: "1.0.0",
		Rules: []auth.Rule{
			{
				Action: auth.Any,
				Type:   "user",
			},
			{
				Action: auth.Any,
				Type:   "stream",
			},
			{
				Action: auth.Any,
				Type:   "subscription",
			},
			{
				Action: auth.Any,
				Type:   "policy",
			},
			{
				Action: auth.Any,
				Type:   "contract",
			},
		},
	},
	"user": {
		Name:    "user",
		Owner:   "admin",
		Version: "1.0.0",
		Rules: []auth.Rule{
			{
				Action: auth.CreateBulk,
				Type:   "stream",
			},
			{
				Action: auth.List,
				Type:   "stream",
			},
			{
				Action: auth.List,
				Type:   "user",
			},
			{
				Action: auth.Any,
				Type:   "stream",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.Any,
				Type:   "contract",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.List,
				Type:   "subscription",
			},
			{
				Action: auth.Any,
				Type:   "subscription",
				Condition: auth.SimpleCondition{
					Key: "ownerID",
				},
			},
			{
				Action: auth.Any,
				Type:   "user",
				Condition: auth.SimpleCondition{
					Key: "id",
				},
			},
			{
				Action: auth.Any,
				Type:   "token",
			},
		},
	},
}

var policiesMu sync.Mutex

var userForDisable = auth.User{
	ID:       "246disable",
	Email:    "disable@example.com",
	Password: "Pass1234!",
	Role:     auth.UserRole,
}

var user = auth.User{
	ID:        email,
	Email:     email,
	Password:  "Pass1234!",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Role:      auth.UserRole,
	Policies:  []auth.Policy{policies["user"]},
}

var userForUpdate = auth.User{
	ID:        "123update",
	Email:     "update@example.com",
	Password:  "Pass5678!",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Role:      auth.UserRole,
	Policies:  []auth.Policy{policies["user"]},
}

var nonAdminForUpdate = auth.User{
	ID:        "nonadmin333update",
	Email:     "nonadminforupdate@example.com",
	Password:  "Pass1234!",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Role:      auth.UserRole,
	Policies:  []auth.Policy{policies["user"]},
}

var admin = auth.User{
	ID:        "admin",
	Email:     "admin@example.com",
	Password:  "Pass1234!",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Role:      auth.AdminRole,
	Policies:  []auth.Policy{policies["admin"]},
}

var nonadmin = auth.User{
	ID:        "nonadmin@example.com",
	Email:     "nonadmin@example.com",
	Password:  "Pass1234!",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Role:      auth.UserRole,
	Policies:  []auth.Policy{policies["user"]},
}

func newServiceWithAdmin() (auth.Service, string, auth.User) {
	hasher := mocks.NewHasher()
	urepo := mocks.NewUserRepository(hasher, admin, policies, &policiesMu)
	prepo := mocks.NewPolicyRepository(policies, &policiesMu)
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()
	ac := mocks.NewAccessControl()

	svc := auth.New(urepo, prepo, hasher, idp, ts, ac)
	key, _ := svc.Login(auth.User{
		Email:    admin.Email,
		Password: admin.Password,
	})
	return svc, key, admin
}

func newService() (auth.Service, string) {
	svc, key, _ := newServiceWithAdmin()
	return svc, key
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
	svc, key := newService()
	ts := newServer(svc)
	defer ts.Close()
	svc.Register(key, nonadmin)
	data := toJSON(testRegisterReq{
		Email:        user.Email,
		Password:     user.Password,
		ContactEmail: user.ContactEmail,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Company:      user.Company,
		Address:      user.Address,
		Phone:        user.Phone,
	})

	invalidData := toJSON(auth.User{})
	invalidEmailData := toJSON(auth.User{
		Email:     invalid,
		Password:  "Pass1234!",
		FirstName: "John",
		LastName:  "Doe",
	})

	nonadminkey, err := svc.Login(auth.User{
		Email:    nonadmin.Email,
		Password: nonadmin.Password,
	})
	require.Nil(t, err, "unexpected error logging nonadmin user: %s", err)

	cases := []struct {
		desc        string
		contentType string
		req         string
		status      int
		key         string
	}{
		{
			desc:        "fail to register new user without admin role",
			contentType: contentType,
			req:         data,
			status:      http.StatusForbidden,
			key:         nonadminkey,
		},
		{
			desc:        "fail to register new user without authorization credentials",
			contentType: contentType,
			req:         data,
			status:      http.StatusForbidden,
			key:         "",
		},
		{
			desc:        "register new user",
			contentType: contentType,
			req:         data,
			status:      http.StatusCreated,
			key:         key,
		},
		{
			desc:        "register existing user",
			contentType: contentType,
			req:         data,
			status:      http.StatusConflict,
			key:         key,
		},
		{
			desc:        "register invalid user",
			contentType: contentType,
			req:         invalidData,
			status:      http.StatusBadRequest,
			key:         key,
		},
		{
			desc:        "register user with invalid email",
			contentType: contentType,
			req:         invalidEmailData,
			status:      http.StatusBadRequest,
			key:         key,
		},
		{
			desc:        "register user with missing content type",
			contentType: "",
			req:         data,
			status:      http.StatusUnsupportedMediaType,
			key:         key,
		},
		{
			desc:        "register user with invalid request format",
			contentType: contentType,
			req:         "}",
			status:      http.StatusBadRequest,
			key:         key,
		},
		{
			desc:        "register user with empty JSON request",
			contentType: contentType,
			req:         "{}",
			status:      http.StatusBadRequest,
			key:         key,
		},
		{
			desc:        "register user with empty request",
			contentType: contentType,
			req:         "",
			status:      http.StatusBadRequest,
			key:         key,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
			token:       tc.key,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestLogin(t *testing.T) {
	svc, key := newService()
	ts := newServer(svc)
	defer ts.Close()

	_, err := svc.Register(key, user)
	require.Nil(t, err, "unexpected error registering user: %s", err)
	_, err = svc.Register(key, userForDisable)
	require.Nil(t, err, "unexpected error registering user: %s", err)
	userForDisable2 := userForDisable
	userForDisable2.Disabled = true
	userForDisable2.Password = ""
	err = svc.UpdateUser(key, userForDisable2)
	require.Nil(t, err, "unexpected error disabling user: %s", err)

	credentials := user
	credentials.ID = ""
	data := toJSON(credentials)
	tokenData := toJSON(map[string]string{"token": user.Email + "|" + user.Role})

	nonexistentUser := auth.User{
		Email:    "nonexistent.user@email.com",
		Password: "Pass1234!",
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
		"login disabled user": {
			contentType: contentType,
			req:         toJSON(userForDisable),
			status:      http.StatusForbidden,
			res:         "",
		},
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
	svc, k := newService()
	ts := newServer(svc)
	defer ts.Close()

	_, err := svc.Register(k, userForUpdate)
	require.Nil(t, err, "unexpected error registering user: %s", err)
	key, err := svc.Login(userForUpdate)
	require.Nil(t, err, "unexpected error logging in user: %s", err)
	_, err = svc.Register(k, nonAdminForUpdate)
	require.Nil(t, err, "unexpected error registering user: %s", err)
	nk, err := svc.Login(nonAdminForUpdate)
	require.Nil(t, err, "unexpected error logging in user: %s", err)

	data := toJSON(testUpdateReq{
		ContactEmail: "john@email.com",
		FirstName:    "john",
		LastName:     "doe",
	})

	cases := []struct {
		desc        string
		contentType string
		req         string
		userID      string
		token       string
		status      int
	}{
		{
			desc:        "update existing user by admin",
			contentType: contentType,
			req:         data,
			userID:      userForUpdate.ID,
			token:       k,
			status:      http.StatusOK,
		},
		{
			desc:        "update self by user",
			contentType: contentType,
			req:         data,
			userID:      userForUpdate.ID,
			token:       key,
			status:      http.StatusOK,
		},
		{
			desc:        "update user by nonadmin",
			contentType: contentType,
			req:         data,
			userID:      userForUpdate.ID,
			token:       nk,
			status:      http.StatusForbidden,
		},
		{
			desc:        "update with bad token",
			contentType: contentType,
			req:         data,
			userID:      userForUpdate.ID,
			token:       "non-existent",
			status:      http.StatusForbidden,
		},
		{
			desc:        "update user with empty token",
			contentType: contentType,
			req:         data,
			userID:      userForUpdate.ID,
			token:       "",
			status:      http.StatusForbidden,
		},
		{
			desc:        "update user with empty content type",
			contentType: "",
			req:         data,
			userID:      userForUpdate.ID,
			token:       k,
			status:      http.StatusUnsupportedMediaType,
		},
	}
	for _, tc := range cases {
		req := testRequest{
			client:      ts.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/%s", ts.URL, tc.userID),
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
	svc, k := newService()
	ts := newServer(svc)
	defer ts.Close()

	id, err := svc.Register(k, user)
	require.Nil(t, err, "unexpected error registering user: %s", err)
	key, err := svc.Login(user)
	require.Nil(t, err, "unexpected error logging in user: %s", err)

	cases := map[string]struct {
		token  string
		id     string
		status int
	}{
		"view existing user": {
			token:  key,
			id:     id,
			status: http.StatusOK,
		},
		"view user with invalid token": {
			token:  invalid,
			id:     id,
			status: http.StatusForbidden,
		},
		"view user with empty token": {
			token:  "",
			id:     id,
			status: http.StatusForbidden,
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client: ts.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/users/%s", ts.URL, tc.id),
			token:  tc.token,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
	}
}

type testUpdateReq struct {
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type testRegisterReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Company      string `json:"company,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
}

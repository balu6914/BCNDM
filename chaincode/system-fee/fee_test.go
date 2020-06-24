package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/datapace/datapace/chaincode/mocks"

	fee "github.com/datapace/datapace/chaincode/system-fee"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	user1CN = "Admin@org1.datapace.com"
	user2CN = "testUser2"
	user3CN = "testUser3"

	user1Cert = `
-----BEGIN CERTIFICATE-----
MIIENzCCAx+gAwIBAgIUOVuKU1UBcvkOgBXWjGD7COxU88gwDQYJKoZIhvcNAQEL
BQAwgaoxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQH
DA1TYW4gRnJhbmNpc2NvMREwDwYDVQQKDAhEYXRhcGFjZTERMA8GA1UECwwIZGF0
YXBhY2UxIDAeBgNVBAMMF0FkbWluQG9yZzEuZGF0YXBhY2UuY29tMSYwJAYJKoZI
hvcNAQkBFhdBZG1pbkBvcmcxLmRhdGFwYWNlLmNvbTAeFw0xOTAxMTUxMTQwNDJa
Fw0yMDAxMTUxMTQwNDJaMIGqMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZv
cm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzERMA8GA1UECgwIRGF0YXBhY2Ux
ETAPBgNVBAsMCGRhdGFwYWNlMSAwHgYDVQQDDBdBZG1pbkBvcmcxLmRhdGFwYWNl
LmNvbTEmMCQGCSqGSIb3DQEJARYXQWRtaW5Ab3JnMS5kYXRhcGFjZS5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCvlZURKWXT6HRxsYBQx8TBudyV
O97Jw/zpSLwpgBee+SniJv/d2VgPT/xunkA06vwxGSVLThz325nU+SCO5qY5lUfF
xdM540qjmqJDpX35Cpm1oKAP+ySaZniz/HUtkVQ5kzrlFQktZJycoOTNpcNNk1QT
UBWQ71Z+vYcvGmXiLCK1z0xvjlfBMg0A4KbHXkRk8NZ7/T61r2vbxwB2RRi5eUts
h2dB9UiUmaCbEmF9IWvGqGxylOHgFEAgHYWvOsUFGk4WmnrZWcwUJLxg0VeA6KYJ
diMGzfjf+a+T/uJM23gXIBb9O6NbMu/n27gHRFh8FsA7WvYK6ugoWF4G20WDAgMB
AAGjUzBRMB0GA1UdDgQWBBRpLVW8op8VY9gzqWKz24/0aZBbhzAfBgNVHSMEGDAW
gBRpLVW8op8VY9gzqWKz24/0aZBbhzAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQB9Q0rI61Eref6L6biwr0LbkRbUWqIo2JISQjmydKoR1Nfs8jtu
lJcSAgrBkUGTCRr/Ze5NNNQB1JtmG0nEcYrc7TghFwioU0mMGtv7vxO/acXFWk6w
cwyv6AWLUwdeN/XjRq09rC3oRBYg9rdqgpj/w5tHtaaRl9d+SFes9999iBjFp+p9
jYZuP4FQ69f6+MfF/+qmZGmz+WyXkpwQWaq5PW2D3n4iC9s3NWwGrbKpxKldq8uw
4ZDaRCP8B9z1myBQ5lQ+FYYRsfGfc4kiMOXVP/njteCsieCQ7/R04M26pJtFmyZV
5tikbH0H4y7Awd8ijx0OzpHMnn5xbHhQ3xJW
-----END CERTIFICATE-----`

	user2Cert = `
-----BEGIN CERTIFICATE-----
MIIB9TCCAZugAwIBAgIUSkK6FlbMHMj8lUytz1/l0IJPw7swCgYIKoZIzj0EAwIw
fzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xHzAdBgNVBAoTFkludGVybmV0IFdpZGdldHMsIEluYy4xDDAK
BgNVBAsTA1dXVzEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTcwMzEzMTAxODAw
WhcNMTgwMjA5MTgxODAwWjAUMRIwEAYDVQQDEwl0ZXN0VXNlcjIwWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAAQYEbEXfqVfArb9u2p8JHiqSwEiE9cQ5mn9CKr76prT
yjZYVmYnImQparjDhtYfiab2cEJaOqJ2J7Au16C6jJ/so2AwXjAOBgNVHQ8BAf8E
BAMCAgQwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUN0EDf71qzYVtIxW5PQQzFvje
7L8wHwYDVR0jBBgwFoAUF2dCPaqegj/ExR2fW8OZ0bWcSBAwCgYIKoZIzj0EAwID
SAAwRQIhAPmJKZYTYiJwvtHbG41XAeIBytyEYA0usiLEvevhN1oFAiB+sLNsJ5Y+
BtcMVPta45X0/aZ5oyI/IJYFWBGSpvgyRQ==
-----END CERTIFICATE-----`

	user3Cert = `
-----BEGIN CERTIFICATE-----
MIIB9DCCAZugAwIBAgIUSC46fLwShh0o0HEzRpvqBe0LEZ0wCgYIKoZIzj0EAwIw
fzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xHzAdBgNVBAoTFkludGVybmV0IFdpZGdldHMsIEluYy4xDDAK
BgNVBAsTA1dXVzEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTcwMzEzMTAxOTAw
WhcNMTgwMjA5MTgxOTAwWjAUMRIwEAYDVQQDEwl0ZXN0VXNlcjMwWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAAQaCAMezFMF7K1xwBy6pR9LP/zVKo/Nh45LMqAuM2IL
mE1ZTFCqc1HJ3ijiSyid+uMOQyo9Jdu2ylj2qECEwYoRo2AwXjAOBgNVHQ8BAf8E
BAMCAgQwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUcXgf0aeO1taJyaDs0c2B274w
h1gwHwYDVR0jBBgwFoAUF2dCPaqegj/ExR2fW8OZ0bWcSBAwCgYIKoZIzj0EAwID
RwAwRAIgJP9ARAqRHl6f2KxB+YJ6ICA9YyYAEkqRnBY4UcTMSIUCID7LFYDewEj3
LmQ6Yvctwv0WEeTCLAuRSmPZL9+hNzX+
-----END CERTIFICATE-----`

	invalid  = "invalid"
	mspID    = "default"
	ccID     = "1"
	ccName   = "fee"
	feeValue = 10000
)

func TestInit(t *testing.T) {
	ts := NewMockTransferService()
	svc := fee.NewService(ts)
	cc := fee.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	payload, err := json.Marshal(feeReq{
		Owner: user2CN,
		Value: feeValue,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		cert   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "create fee with valid user",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", string(payload)),
			status: int32(shim.OK),
		},
		{
			desc:   "create fee with invalid user",
			cert:   invalid,
			args:   util.ToChaincodeArgs("init", string(payload)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create fee with invalid data",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", "{"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create fee with invalid num of arguments",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", string(payload), "other_arg"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create fee with invalid function name",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs(invalid, string(payload)),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		stub.MockCreator(mspID, tc.cert)
		res := stub.MockInit(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", res.GetMessage(), tc.status, res.Status))
	}
}

func TestFee(t *testing.T) {
	ts := NewMockTransferService()
	svc := fee.NewService(ts)
	cc := fee.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	payload, err := json.Marshal(feeReq{
		Owner: user2CN,
		Value: 100,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	stub.MockCreator(mspID, user1Cert)
	res := stub.MockInvoke(ccID, util.ToChaincodeArgs("setFee", string(payload)))
	require.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), res.Status))

	cases := map[string]struct {
		args   [][]byte
		status int32
	}{
		"get fee": {
			args:   util.ToChaincodeArgs("fee"),
			status: int32(shim.OK),
		},
		"get fee with invalid function name": {
			args:   util.ToChaincodeArgs(invalid),
			status: int32(shim.ERROR),
		},
	}

	for desc, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.Status))
	}
}

func TestSetFee(t *testing.T) {
	ts := NewMockTransferService()
	svc := fee.NewService(ts)
	cc := fee.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	payload, err := json.Marshal(feeReq{
		Owner: user2CN,
		Value: 100,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	tooBigFeePayload, err := json.Marshal(feeReq{
		Owner: user2CN,
		Value: 100001,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		cert   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "set fee",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("setFee", string(payload)),
			status: int32(shim.OK),
		},
		{
			desc:   "set fee that is over 100%",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("setFee", string(tooBigFeePayload)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "set fee with user with no permission",
			cert:   user2Cert,
			args:   util.ToChaincodeArgs("setFee", string(payload)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "set fee with invalid data",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("setFee", "{"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "set fee with invalid num of arguments",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("setFee", string(payload), "other_arg"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "set fee with invalid function name",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs(invalid, string(payload)),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		stub.MockCreator(mspID, tc.cert)
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

func TestTransfer(t *testing.T) {
	ts := NewMockTransferService()
	svc := fee.NewService(ts)
	cc := fee.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	initPayload, err := json.Marshal(feeReq{
		Owner: user2CN,
		Value: 10000,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	stub.MockCreator(mspID, user1Cert)
	res := stub.MockInvoke(ccID, util.ToChaincodeArgs("setFee", string(initPayload)))
	require.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), res.Status))

	payload, err := json.Marshal(transferReq{
		To:    user3CN,
		Value: 100,
	})
	require.Nil(t, err, fmt.Sprintf("unexppected error: %s", err))

	cases := []struct {
		desc   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "transfer tokens with fee",
			args:   util.ToChaincodeArgs("transfer", string(payload)),
			status: int32(shim.OK),
		},
		{
			desc:   "transfer tokens with invalid data",
			args:   util.ToChaincodeArgs("transfer", "{"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with invalid num of arguments",
			args:   util.ToChaincodeArgs("transfer", string(payload), "other_arg"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with invalid function name",
			args:   util.ToChaincodeArgs(invalid, string(payload)),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

type feeReq struct {
	Owner string `json:"owner"`
	Value uint64 `json:"value"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

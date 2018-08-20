package main_test

import (
	"encoding/json"
	"fmt"
	"monetasa/chaincode/mocks"
	"testing"

	fee "monetasa/chaincode/system-fee"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	user1CN = "Admin@org1.monetasa.com"
	user2CN = "testUser2"
	user3CN = "testUser3"

	user1Cert = `
-----BEGIN CERTIFICATE-----
MIICHDCCAcOgAwIBAgIRAP51ybjVIXxBiihZb6q7aqIwCgYIKoZIzj0EAwIwdTEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGjAYBgNVBAoTEW9yZzEubW9uZXRhc2EuY29tMR0wGwYDVQQDExRj
YS5vcmcxLm1vbmV0YXNhLmNvbTAeFw0xODA4MTYxNDE3MDJaFw0yODA4MTMxNDE3
MDJaMFwxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH
Ew1TYW4gRnJhbmNpc2NvMSAwHgYDVQQDDBdBZG1pbkBvcmcxLm1vbmV0YXNhLmNv
bTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEyrrQ78tw+xqDhYnpp+v4cfoxLU
W0VHbW7evgPrKQuutib3lMjCooMRq4+fs3VSfjl0Ho+JPq9e6/qR0S1k7gijTTBL
MA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAILb7EO8q
Tiz3Hlw+UtOjSLej4Foxbl3UfJ/jZXDCobodMAoGCCqGSM49BAMCA0cAMEQCIC5M
PXLK69sIYS4rZ/qf7Jz1anGppP/98gP94OntQnbGAiB5UsofXK5Zw+yGcs6zks8b
qTm7Nxz8TYQ9GhXIxPg5ag==
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

	payload, err := json.Marshal(feeReq{Value: feeValue})
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
		Value: 100,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	tooBigFeePayload, err := json.Marshal(feeReq{
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
	Value uint64 `json:"value"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

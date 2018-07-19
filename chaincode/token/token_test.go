package main_test

import (
	"encoding/json"
	"fmt"
	token "monetasa/chaincode/token"
	"monetasa/chaincode/token/mocks"
	"testing"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	user1CN = "testUser"
	user2CN = "testUser2"
	user3CN = "testUser3"

	user1Cert = `
-----BEGIN CERTIFICATE-----
MIIB9DCCAZqgAwIBAgIUDda1JZnuPZ5dlcwSlOmU/KWSn7MwCgYIKoZIzj0EAwIw
fzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xHzAdBgNVBAoTFkludGVybmV0IFdpZGdldHMsIEluYy4xDDAK
BgNVBAsTA1dXVzEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTcwMjEzMTQyOTAw
WhcNMTgwMTEyMjIyOTAwWjATMREwDwYDVQQDEwh0ZXN0VXNlcjBZMBMGByqGSM49
AgEGCCqGSM49AwEHA0IABKqm8JxN53RW1/muhqPxO7F7dnEMhguy23MVj4CXybqP
rY70z4AJdXKZTxPeU06kIwb1c0NMii+NMUAjp624z0qjYDBeMA4GA1UdDwEB/wQE
AwICBDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBT6YW1Vq07nRK502xj3Y76/lqsu
3zAfBgNVHSMEGDAWgBQXZ0I9qp6CP8TFHZ9bw5nRtZxIEDAKBggqhkjOPQQDAgNI
ADBFAiEA5tzFnCPvASFWQku49vrGNGhmJeASlbo2W1ipWarkTlQCIHpI4eWFj6na
4Xtb5djZAMGlfC2jJl/FTKzFj/xd4s3E
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

	invalid = "invalid"
	mspID   = "default"
	ccID    = "1"
	ccName  = "token"
)

var td = tokenData{
	Name:        "DatapaceToken",
	Symbol:      "TAS",
	Decimals:    8,
	TotalSupply: 1000,
}

func TestInit(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	reqPayload, err := json.Marshal(td)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		cert   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "create token with valid user",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", string(reqPayload)),
			status: int32(shim.OK),
		},
		{
			desc:   "create token with invalid user",
			cert:   invalid,
			args:   util.ToChaincodeArgs("init", string(reqPayload)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create token with invalid data",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", "{"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create token with invalid num of arguments",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs("init", string(reqPayload), "other_arg"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "create token with invalid function name",
			cert:   user1Cert,
			args:   util.ToChaincodeArgs(invalid, string(reqPayload), "other_arg"),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		stub.MockCreator(mspID, tc.cert)
		res := stub.MockInit(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

func TestTotalSupply(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	createToken(t, stub)

	res := stub.MockInvoke(ccID, util.ToChaincodeArgs("totalSupply"))
	assert.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("get token total supply: expected %d got %d", int32(shim.OK), res.Status))
}

func TestBalanceOf(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	// Prepare data
	createToken(t, stub)

	tr := transferReq{
		To:    user2CN,
		Value: 100,
	}
	transferData, err := json.Marshal(tr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	// Transfer some tokens to user 2.
	res := stub.MockInvoke(ccID, util.ToChaincodeArgs("transfer", string(transferData)))
	require.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("initial transfer failed: expected %d got %d", int32(shim.OK), res.Status))

	br := balanceReq{
		Owner: user2CN,
	}
	balanceData, err := json.Marshal(br)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	nonexistentBr := balanceReq{
		Owner: user3CN,
	}
	nonexistentBalanceData, err := json.Marshal(nonexistentBr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		args   [][]byte
		status int32
	}{
		"get balance of existing user": {
			args:   util.ToChaincodeArgs("balanceOf", string(balanceData)),
			status: int32(shim.OK),
		},
		"get balance of nonexistent user": {
			args:   util.ToChaincodeArgs("balanceOf", string(nonexistentBalanceData)),
			status: int32(shim.OK),
		},
		"get balance with invalid request": {
			args:   util.ToChaincodeArgs("balanceOf", "}"),
			status: int32(shim.ERROR),
		},
		"get balance with invalid number of requests": {
			args:   util.ToChaincodeArgs("balanceOf", string(balanceData), ""),
			status: int32(shim.ERROR),
		},
	}

	for desc, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.Status))
	}
}

func TestTransfer(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	// Prepare data
	createToken(t, stub)

	tr := transferReq{
		To:    user2CN,
		Value: 100,
	}
	transferData, err := json.Marshal(tr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	bigTr := transferReq{
		To:    user2CN,
		Value: 10000,
	}
	bigTransferData, err := json.Marshal(bigTr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "transfer tokens to valid account",
			args:   util.ToChaincodeArgs("transfer", string(transferData)),
			status: int32(shim.OK),
		},
		{
			desc:   "transfer too many tokens",
			args:   util.ToChaincodeArgs("transfer", string(bigTransferData)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with invalid request",
			args:   util.ToChaincodeArgs("transfer", "}"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with too many arguments",
			args:   util.ToChaincodeArgs("transfer", string(transferData), ""),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

func TestTransferFrom(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	// Prepare data
	createToken(t, stub)

	setupTransfers(t, stub)

	stub.MockCreator(mspID, user1Cert)

	tr := transferFromReq{
		From:  user2CN,
		To:    user3CN,
		Value: 100,
	}
	transferData, err := json.Marshal(tr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	noAllowanceTr := transferFromReq{
		From:  user3CN,
		To:    user2CN,
		Value: 100,
	}
	noAllowanceTd, err := json.Marshal(noAllowanceTr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "transfer tokens from one account to other",
			args:   util.ToChaincodeArgs("transferFrom", string(transferData)),
			status: int32(shim.OK),
		},
		{
			desc:   "transfer tokens from one account to other with no allowance",
			args:   util.ToChaincodeArgs("transferFrom", string(noAllowanceTd)),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with invalid request",
			args:   util.ToChaincodeArgs("transferFrom", "}"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "transfer tokens with invalid number of arguments",
			args:   util.ToChaincodeArgs("transferFrom", string(transferData), ""),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

func TestApprove(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	createToken(t, stub)

	ar := approveReq{
		Spender: user2CN,
		Value:   100,
	}
	ad, err := json.Marshal(ar)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		args   [][]byte
		status int32
	}{
		{
			desc:   "approve allowance for spender",
			args:   util.ToChaincodeArgs("approve", string(ad)),
			status: int32(shim.OK),
		},
		{
			desc:   "approve allowance with invalid request",
			args:   util.ToChaincodeArgs("approve", "}"),
			status: int32(shim.ERROR),
		},
		{
			desc:   "approve allowance with invalid number of arguments",
			args:   util.ToChaincodeArgs("approve", string(ad), ""),
			status: int32(shim.ERROR),
		},
	}

	for _, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", tc.desc, tc.status, res.Status))
	}
}

func TestAllowance(t *testing.T) {
	svc := token.NewService()
	cc := token.NewChaincode(svc)
	stub := mocks.NewFullMockStub(ccName, cc)

	createToken(t, stub)

	stub.MockCreator(mspID, user2Cert)

	approve := approveReq{
		Spender: user1CN,
		Value:   100,
	}
	approveData, err := json.Marshal(approve)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	approveRes := stub.MockInvoke(ccID, util.ToChaincodeArgs("approve", string(approveData)))
	require.Equal(t, int32(shim.OK), approveRes.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), approveRes.Status))

	stub.MockCreator(mspID, user1Cert)

	ar := allowanceReq{
		Owner:   user2CN,
		Spender: user1CN,
	}
	ad, err := json.Marshal(ar)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	invalidAr := allowanceReq{}
	invalidAd, err := json.Marshal(invalidAr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		args   [][]byte
		status int32
	}{
		"read allowance for given spender owner pair": {
			args:   util.ToChaincodeArgs("allowance", string(ad)),
			status: int32(shim.OK),
		},
		"read allowance for empty spender and owner": {
			args:   util.ToChaincodeArgs("allowance", string(invalidAd)),
			status: int32(shim.ERROR),
		},
	}

	for desc, tc := range cases {
		res := stub.MockInvoke(ccID, tc.args)
		assert.Equal(t, tc.status, res.Status, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.Status))
	}
}

type tokenData struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
}

type balanceReq struct {
	Owner string `json:"owner"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

type transferFromReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

type approveReq struct {
	Spender string `json:"spender"`
	Value   uint64 `json:"value"`
}

type allowanceReq struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
}

func createToken(t *testing.T, stub *mocks.FullMockStub) {
	reqPayload, err := json.Marshal(td)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	stub.MockCreator(mspID, user1Cert)
	res := stub.MockInit(ccID, util.ToChaincodeArgs("init", string(reqPayload)))
	require.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("creating token failed"))
}

func setupTransfers(t *testing.T, stub *mocks.FullMockStub) {
	firstTr := transferReq{
		To:    user2CN,
		Value: 100,
	}
	firstTransferData, err := json.Marshal(firstTr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	res := stub.MockInvoke(ccID, util.ToChaincodeArgs("transfer", string(firstTransferData)))
	require.Equal(t, int32(shim.OK), res.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), res.Status))

	otherTr := transferReq{
		To:    user3CN,
		Value: 100,
	}
	otherTransferData, err := json.Marshal(otherTr)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	otherRes := stub.MockInvoke(ccID, util.ToChaincodeArgs("transfer", string(otherTransferData)))
	require.Equal(t, int32(shim.OK), otherRes.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), otherRes.Status))

	stub.MockCreator(mspID, user2Cert)

	ar := approveReq{
		Spender: user1CN,
		Value:   100,
	}
	approveData, err := json.Marshal(ar)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	approveRes := stub.MockInvoke(ccID, util.ToChaincodeArgs("approve", string(approveData)))
	require.Equal(t, int32(shim.OK), approveRes.Status, fmt.Sprintf("expected %d got %d", int32(shim.OK), approveRes.Status))
}

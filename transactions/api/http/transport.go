package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/transactions"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	contentType = "application/json"
	page        = "page"
	limit       = "limit"
	owner       = "owner"
	partner     = "partner"
	defPage     = 0
	defLimit    = 10
	defOwner    = false
	defPartner  = false
	shareScale  = 1000
)

var (
	errMalformedEntity        = errors.New("malformed entity")
	errUnauthorizedAccess     = errors.New("missing or invalid credentials provided")
	errUnsupportedContentType = errors.New("unsupported content type")
	auth                      datapace.AuthServiceClient
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc transactions.Service, ac datapace.AuthServiceClient) http.Handler {
	auth = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()
	r.Get("/tokens", kithttp.NewServer(
		balanceEndpoint(svc),
		decodeBalanceReq,
		encodeResponse,
		opts...,
	))
	r.Post("/tokens/buy", kithttp.NewServer(
		buyEndpoint(svc),
		decodeBuyReq,
		encodeResponse,
		opts...,
	))
	r.Post("/tokens/withdraw", kithttp.NewServer(
		withdrawEndpoint(svc),
		decodeWithdrawReq,
		encodeResponse,
		opts...,
	))
	r.Post("/contracts", kithttp.NewServer(
		createContractsEndpoint(svc),
		decodeCreateContractsReq,
		encodeResponse,
		opts...,
	))
	r.Patch("/contracts/sign", kithttp.NewServer(
		signContractEndpoint(svc),
		decodeSignContractReq,
		encodeResponse,
		opts...,
	))
	r.Get("/contracts", kithttp.NewServer(
		listContractsEndpoint(svc),
		decodeListContractsReq,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeBalanceReq(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := balanceReq{userID: id}

	return req, nil
}

func decodeBuyReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req buyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.userID = id

	return req, nil
}

func decodeWithdrawReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req withdrawReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.userID = id

	return req, nil
}

func decodeCreateContractsReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req createContractsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ownerID = id

	for i := range req.Items {
		req.Items[i].Share *= shareScale
	}

	return req, nil
}

func decodeSignContractReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req signContractReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.partnerID = id

	return req, nil
}

func decodeListContractsReq(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := authorize(r)
	if err != nil {
		return nil, err
	}
	pageNo := uintQueryParam(r, page, defPage)
	limit := uintQueryParam(r, limit, defLimit)

	isOwner := boolQueryParam(r, owner, defOwner)
	isPartner := boolQueryParam(r, partner, defPartner)
	var role transactions.Role
	if isOwner {
		role = transactions.Owner
	}
	if isPartner {
		role = transactions.Partner
	}
	if isOwner && isPartner {
		role = transactions.AllRoles
	}

	req := listContractsReq{
		userID: id,
		page:   pageNo,
		limit:  limit,
		role:   role,
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(apiRes); ok {
		for k, v := range ar.headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.code())

		if ar.empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func authorize(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errUnauthorizedAccess
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := auth.Identify(ctx, &datapace.Token{Value: token})
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.Unauthenticated {
			return "", errUnauthorizedAccess
		}
		return "", err
	}

	return id.GetValue(), nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case errMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case transactions.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func boolQueryParam(req *http.Request, name string, fallback bool) bool {
	vals := bone.GetQuery(req, name)
	if len(vals) == 0 {
		return fallback
	}

	val, err := strconv.ParseBool(vals[0])
	if err != nil {
		return fallback
	}

	return val
}

func uintQueryParam(req *http.Request, name string, fallback uint64) uint64 {
	vals := bone.GetQuery(req, name)
	if len(vals) == 0 {
		return fallback
	}

	val, err := strconv.ParseUint(vals[0], 10, 64)
	if err != nil {
		return fallback
	}

	return val
}

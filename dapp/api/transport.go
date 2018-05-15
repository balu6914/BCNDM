package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"gopkg.in/mgo.v2/bson"

	"monetasa"
	"monetasa/auth"
	"monetasa/auth/client"

	"monetasa/dapp"
)

var (
	dappService dapp.Service
	authClient  client.AuthClient
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(sr dapp.Service, ac client.AuthClient) http.Handler {
	dappService = sr
	authClient = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/streams", kithttp.NewServer(
		addStreamEndpoint(sr),
		decodeCreateStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/streams/bulk", kithttp.NewServer(
		addBulkStreamEndpoint(sr),
		decodeCreateBulkStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/search", kithttp.NewServer(
		searchStreamEndpoint(sr),
		decodeSearchStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Put("/streams/:id", kithttp.NewServer(
		updateStreamEndpoint(sr),
		decodeUpdateStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:id", kithttp.NewServer(
		viewStreamEndpoint(sr),
		decodeReadStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Delete("/streams/:id", kithttp.NewServer(
		removeStreamEndpoint(sr),
		decodeDeleteStreamRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())

	return r
}

func authenticate(r *http.Request) (string, error) {
	// apiKey is a JWT token
	apiKey := r.Header.Get("Authorization")

	if apiKey == "" {
		return "", auth.ErrUnauthorizedAccess
	}

	userId, err := authClient.VerifyToken(apiKey)
	if err != nil {
		return "", err
	}

	// id is an email of the user
	return userId, nil
}

func decodeCreateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user, err := authenticate(r)
	if err != nil {
		return nil, err
	}

	var stream dapp.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	stream.ID = bson.NewObjectId()

	req := createStreamReq{
		User:   user,
		Stream: stream,
	}
	return req, nil
}

func decodeCreateBulkStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user, err := authenticate(r)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("csv")
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	// remove the first csv row
	records = records[1:]

	var streams []dapp.Stream
	for _, record := range records {
		price, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		longitude, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		latitude, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}

		stream := dapp.Stream{
			Owner:       user,
			ID:          bson.NewObjectId(),
			Name:        record[0],
			Type:        record[1],
			Description: record[2],
			Price:       price,
			Location: dapp.Location{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
			URL: record[6],
		}

		streams = append(streams, stream)
	}

	req := createBulkStreamRequest{
		Streams: streams,
	}

	return req, nil
}

func decodeUpdateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user, err := authenticate(r)
	if err != nil {
		return nil, err
	}

	var stream dapp.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	req := updateStreamReq{
		User:     user,
		StreamId: bone.GetValue(r, "id"),
		Stream:   stream,
	}
	return req, nil
}

func decodeReadStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	_, err := authenticate(r)
	if err != nil {
		return nil, err
	}

	req := readStreamReq{
		StreamId: bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeDeleteStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user, err := authenticate(r)
	if err != nil {
		return nil, err
	}

	req := deleteStreamReq{
		User:     user,
		StreamId: bone.GetValue(r, "id"),
	}
	return req, nil

}

func decodeSearchStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if _, err := authenticate(r); err != nil {
		return nil, err
	}

	q := r.URL.Query()

	var req searchStreamReq
	req.Type = q["type"][0]
	req.x0, _ = strconv.ParseFloat(q["x0"][0], 64)
	req.y0, _ = strconv.ParseFloat(q["y0"][0], 64)
	req.x1, _ = strconv.ParseFloat(q["x1"][0], 64)
	req.y1, _ = strconv.ParseFloat(q["y1"][0], 64)
	req.x2, _ = strconv.ParseFloat(q["x2"][0], 64)
	req.y2, _ = strconv.ParseFloat(q["y2"][0], 64)
	req.x3, _ = strconv.ParseFloat(q["x3"][0], 64)
	req.y3, _ = strconv.ParseFloat(q["y3"][0], 64)

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

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case dapp.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case dapp.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case dapp.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case dapp.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case dapp.ErrUnsupportedContentType:
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

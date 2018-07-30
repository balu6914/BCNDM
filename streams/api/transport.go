package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

	"monetasa"
	"monetasa/streams"
)

var (
	streamsService streams.Service
	authService    streams.Authorization
	searchFields   = [4][2]string{{"x0", "y0"}, {"x1", "y1"}, {"x2", "y2"}, {"x3", "y3"}}
	fileFields     = []string{"name", "type", "description", "price", "longitude", "latitude", "url"}
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc streams.Service, auth streams.Authorization) http.Handler {
	streamsService = svc
	authService = auth

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/streams", kithttp.NewServer(
		addStreamEndpoint(svc),
		decodeAddStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/streams/bulk", kithttp.NewServer(
		addBulkStreamEndpoint(svc),
		decodeAddBulkStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/search", kithttp.NewServer(
		searchStreamEndpoint(svc),
		decodeSearchStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Put("/streams/:id", kithttp.NewServer(
		updateStreamEndpoint(svc),
		decodeUpdateStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:id", kithttp.NewServer(
		viewStreamEndpoint(svc),
		decodeViewStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Delete("/streams/:id", kithttp.NewServer(
		removeStreamEndpoint(svc),
		decodeRemoveStreamRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())

	return r
}

func decodeAddStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	var stream streams.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	req := createStreamReq{
		owner:  owner,
		stream: stream,
	}
	return req, nil
}

func stringInSlice(a string, list []string) bool {
	for _, s := range list {
		if s == a {
			return true
		}
	}

	return false
}

type csvFile struct {
	columns []string
	records [][]string
}

func readFile(r *http.Request) (*csvFile, error) {
	file, _, err := r.FormFile("csv")
	if err != nil {
		return nil, streams.ErrMalformedData
	}

	reader := csv.NewReader(file)

	content, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	columns, records := content[0], content[1:]
	for idx, field := range columns {
		columns[idx] = strings.ToLower(field)
	}
	for _, field := range fileFields {
		if !stringInSlice(field, columns) {
			return nil, streams.ErrMalformedData
		}
	}

	ret := &csvFile{
		columns: columns,
		records: records,
	}
	return ret, nil
}

func parseStream(record []string, keys map[string]int) (*streams.Stream, error) {
	price, err := strconv.ParseUint(record[keys["price"]], 10, 64)
	if err != nil {
		return nil, err
	}

	longitude, err := strconv.ParseFloat(record[keys["longitude"]], 64)
	if err != nil {
		return nil, err
	}

	latitude, err := strconv.ParseFloat(record[keys["latitude"]], 64)
	if err != nil {
		return nil, err
	}

	ret := &streams.Stream{
		Name:        record[keys["name"]],
		Type:        record[keys["type"]],
		Description: record[keys["description"]],
		Price:       price,
		Location: streams.Location{
			Type:        "Point",
			Coordinates: [2]float64{longitude, latitude},
		},
		URL: record[6],
	}
	return ret, nil
}

func decodeAddBulkStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	file, err := readFile(r)
	if err != nil {
		return nil, err
	}

	// Keys represent a map with csv field names as keys and field col
	// numbers as values.
	keys := make(map[string]int)
	for idx, attr := range file.columns {
		keys[attr] = idx
	}

	s := []streams.Stream{}
	for _, record := range file.records {
		stream, err := parseStream(record, keys)
		if err != nil {
			return nil, err
		}
		s = append(s, *stream)
	}

	req := createBulkStreamReq{
		owner:   owner,
		Streams: s,
	}

	return req, nil
}

func decodeUpdateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	var stream streams.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	req := updateStreamReq{
		owner:  owner,
		id:     bone.GetValue(r, "id"),
		stream: stream,
	}
	return req, nil
}

func decodeViewStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if _, err := authService.Authorize(r); err != nil {
		return nil, err
	}

	req := readStreamReq{
		id: bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeRemoveStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	req := deleteStreamReq{
		owner: owner,
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeSearchStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if _, err := authService.Authorize(r); err != nil {
		return nil, err
	}

	q := r.URL.Query()
	req := searchStreamReq{}
	if len(q["type"]) != 1 {
		return nil, streams.ErrMalformedData
	}
	req.locationType = q["type"][0]
	var err error
	for i, v := range searchFields {
		// X and Y coordinates are q[v[0]] and q[v[1]] respectively.
		if len(q[v[0]]) != 1 || len(q[v[1]]) != 1 {
			return nil, streams.ErrMalformedData
		}
		req.points = append(req.points, []float64{0, 0})
		for j := range v {
			req.points[i][j], err = strconv.ParseFloat(q[v[j]][0], 64)
			if err != nil {
				return nil, streams.ErrMalformedData
			}
		}
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

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case streams.ErrMalformedData:
		w.WriteHeader(http.StatusBadRequest)
	case streams.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case streams.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case streams.ErrConflict:
		w.WriteHeader(http.StatusConflict)
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

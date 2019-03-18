package http

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
	"gopkg.in/mgo.v2/bson"

	"datapace"
	"datapace/streams"
)

const (
	defLocType  = "Point"
	gmailSuffix = "@gmail.com"
)

var (
	authService    streams.Authorization
	locationPoints = [4][2]string{{"x0", "y0"}, {"x1", "y1"}, {"x2", "y2"}, {"x3", "y3"}}
	fields         = []string{"name", "type", "description", "snippet", "price", "longitude", "latitude", "url", "terms"}
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc streams.Service, auth streams.Authorization) http.Handler {
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
		addBulkStreamsEndpoint(svc),
		decodeAddBulkStreamsRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams", kithttp.NewServer(
		searchStreamsEndpoint(svc),
		decodeSearchStreamsRequest,
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

	r.GetFunc("/version", datapace.Version())

	return r
}

func checkEmail(userEmail datapace.UserEmail) (string, error) {
	email := strings.ToLower(userEmail.Email)
	contactEmail := strings.ToLower(userEmail.ContactEmail)
	if strings.HasSuffix(email, gmailSuffix) {
		return userEmail.Email, nil
	}
	if strings.HasSuffix(contactEmail, gmailSuffix) {
		return userEmail.ContactEmail, nil
	}
	return "", streams.ErrMalformedData
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

	if stream.Location.Type == "" {
		stream.Location.Type = defLocType
	}

	if stream.External {
		token := r.Header.Get("Authorization")
		ownerEmail, err := authService.Email(token)
		if err != nil {
			return nil, err
		}
		email, err := checkEmail(ownerEmail)
		if err != nil {
			return nil, err
		}
		stream.BQ.Email = email
	}

	stream.Owner = owner
	req := addStreamReq{
		owner:  owner,
		stream: stream,
	}
	return req, nil
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
	if err != nil || len(content) < 2 {
		return nil, streams.ErrMalformedData
	}

	columns, records := content[0], content[1:]
	for idx, field := range columns {
		columns[idx] = strings.ToLower(field)
	}

	for _, field := range fields {
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

	// Convert Metadata from string to bson.M if present
	data := []byte(record[keys["metadata"]])
	metadata := bson.M{}
	if len(data) != 0 {
		json.Unmarshal(data, &metadata)
	}

	ret := &streams.Stream{
		Visibility:  streams.Visibility(record[keys["visibility"]]),
		Name:        record[keys["name"]],
		Type:        record[keys["type"]],
		Description: record[keys["description"]],
		Snippet:     record[keys["snippet"]],
		Price:       price,
		Location: streams.Location{
			Type:        "Point",
			Coordinates: [2]float64{longitude, latitude},
		},
		URL:      record[keys["url"]],
		Terms:    record[keys["terms"]],
		Metadata: metadata,
	}
	return ret, nil
}

func decodeAddBulkStreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), fileContentType) {
		return nil, streams.ErrWrongType
	}

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

		stream.Owner = owner

		if stream.Location.Type == "" {
			stream.Location.Type = defLocType
		}
		s = append(s, *stream)
	}

	if len(s) < 1 {
		return nil, streams.ErrMalformedData
	}

	req := addBulkStreamsReq{
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

	if stream.Location.Type == "" {
		stream.Location.Type = defLocType
	}

	req := updateStreamReq{
		owner:  owner,
		id:     bone.GetValue(r, "id"),
		stream: stream,
	}
	return req, nil
}

func decodeViewStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	req := viewStreamReq{
		id:    bone.GetValue(r, "id"),
		owner: owner,
	}
	return req, nil
}

func decodeRemoveStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	req := removeStreamReq{
		owner: owner,
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeSearchStreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authService.Authorize(r)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	req := newSearchStreamReq()
	req.user = owner

	if err := searchFields(&req, q); err != nil {
		return nil, err
	}

	if err := locationFields(&req, q); err != nil {
		return nil, err
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
	case streams.ErrMalformedData, streams.ErrBigQuery:
		w.WriteHeader(http.StatusBadRequest)
	case streams.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case streams.ErrInvalidBQAccess:
		w.WriteHeader(http.StatusUnauthorized)
	case streams.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case streams.ErrWrongType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
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

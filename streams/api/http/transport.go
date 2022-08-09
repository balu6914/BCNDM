package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/datapace/datapace"
	"github.com/datapace/datapace/auth"
	authproto "github.com/datapace/datapace/proto/auth"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

	"github.com/datapace/datapace/errors"
	"github.com/datapace/datapace/streams"
	"github.com/datapace/datapace/streams/executions"
)

const (
	defLocType  = "Point"
	streamType  = "stream"
	gmailSuffix = "@gmail.com"
)

var (
	authService    streams.Authorization
	locationPoints = [4][2]string{{"x0", "y0"}, {"x1", "y1"}, {"x2", "y2"}, {"x3", "y3"}}
	fields         = []string{"name", "type", "description", "snippet", "price", "longitude", "latitude", "url", "terms"}
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc streams.Service, auth streams.Authorization, accessSvc streams.AccessControl) http.Handler {
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
		updateStreamEndpoint(svc, accessSvc),
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

	r.Get("/export", kithttp.NewServer(
		exportStreamsEndpoint(svc),
		decodeExportStreamsRequest,
		encodeExportStreamsResponse,
		opts...,
	))

	r.Post("/search", kithttp.NewServer(
		searchStreamsEndpoint(svc),
		decodeSearchStreamsJsonRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())

	return r
}

func checkEmail(userEmail *authproto.UserEmail) (string, error) {
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
	var stream streams.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	ar := &authproto.AuthRequest{
		Action:     int64(auth.Create),
		Token:      r.Header.Get("Authorization"),
		Type:       streamType,
		Attributes: stream.Attributes(),
	}

	owner, err := authService.Authorize(ar)
	if err != nil {
		return nil, err
	}

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

func decodeAddBulkStreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentTypeFormData) {
		return nil, streams.ErrWrongType
	}

	ar := &authproto.AuthRequest{
		Action: int64(auth.CreateBulk),
		Token:  r.Header.Get("Authorization"),
		Type:   streamType,
	}
	owner, err := authService.Authorize(ar)
	if err != nil {
		return nil, err
	}

	file, fileHeader, err := r.FormFile("data")
	if err != nil {
		return nil, streams.ErrMalformedData
	}
	fileName := fileHeader.Filename

	var s []streams.Stream
	switch {
	case strings.HasSuffix(fileName, ".csv"):
		s, err = decodeCsvStreams(file, owner)
	case strings.HasSuffix(fileName, ".json"):
		s, err = decodeJsonStreams(file, owner)
	default:
		err = streams.ErrMalformedData
	}
	if err != nil {
		return nil, err
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

func decodeCsvStreams(file multipart.File, owner string) ([]streams.Stream, error) {
	csvFile, err := readCsvFile(file)
	if err != nil {
		return nil, err
	}

	// Keys represent a map with csv field names as keys and field col
	// numbers as values.
	keys := make(map[string]int)
	for idx, attr := range csvFile.columns {
		keys[attr] = idx
	}

	s := []streams.Stream{}
	for _, record := range csvFile.records {
		stream, err := streams.NewFromCsv(record, keys)
		if err != nil {
			return nil, err
		}
		stream.Owner = owner
		if stream.Location.Type == "" {
			stream.Location.Type = defLocType
		}
		s = append(s, *stream)
	}
	return s, nil
}

type csvFile struct {
	columns []string
	records [][]string
}

func readCsvFile(file multipart.File) (*csvFile, error) {
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

func decodeJsonStreams(file multipart.File, owner string) ([]streams.Stream, error) {
	ss := []streams.Stream{}
	err := json.NewDecoder(file).Decode(&ss)
	// take care on the stream.Owner and stream.Location.Type
	results := []streams.Stream{}
	for _, s := range ss {
		s.Owner = owner
		if s.Location.Type == "" {
			s.Location.Type = defLocType
		}
		results = append(results, s)
	}
	return results, err
}

func decodeUpdateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var stream streams.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if stream.Location.Type == "" {
		stream.Location.Type = defLocType
	}

	req := updateStreamReq{
		owner:  r.Header.Get("Authorization"),
		id:     bone.GetValue(r, "id"),
		stream: stream,
	}
	return req, nil
}

func decodeViewStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.Read),
		Token:  r.Header.Get("Authorization"),
		Type:   streamType,
	}
	owner, err := authService.Authorize(ar)
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
	req := removeStreamReq{
		owner: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeSearchStreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.List),
		Token:  r.Header.Get("Authorization"),
		Type:   streamType,
	}
	owner, err := authService.Authorize(ar)
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

func decodeExportStreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.List),
		Token:  r.Header.Get("Authorization"),
		Type:   streamType,
	}
	owner, err := authService.Authorize(ar)
	if err != nil {
		return nil, err
	}
	req := exportStreamsReq{owner: owner}
	return req, nil
}

func decodeSearchStreamsJsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.List),
		Token:  r.Header.Get("Authorization"),
		Type:   streamType,
	}
	owner, err := authService.Authorize(ar)
	if err != nil {
		return nil, err
	}

	req := newSearchStreamReq()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	req.user = owner

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentTypeJson)

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

func encodeExportStreamsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentTypeCsv)
	resp := response.(exportStreamsResp)
	for k, v := range resp.headers() {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.code())
	if resp.empty() {
		return nil
	}
	csvWriter := csv.NewWriter(w)
	csvRecs := [][]string{
		streams.CsvHeader,
	}
	for _, stream := range resp.streams {
		csvRec, err := stream.Csv()
		if err != nil {
			return err
		}
		csvRecs = append(csvRecs, csvRec)
	}
	return csvWriter.WriteAll(csvRecs)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentTypeJson)
	switch errVal := err.(type) {
	case errors.Error:
		switch {
		case errors.Contains(errVal, executions.ErrCrateDataset), errors.Contains(errVal, executions.ErrCrateAlgorithm):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errVal, streams.ErrMalformedData), errors.Contains(errVal, streams.ErrBigQuery):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errVal, streams.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(errVal, streams.ErrInvalidBQAccess):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Contains(errVal, streams.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Contains(errVal, streams.ErrWrongType):
			w.WriteHeader(http.StatusUnsupportedMediaType)
		case errors.Contains(errVal, streams.ErrConflict):
			w.WriteHeader(http.StatusConflict)
		}
		if errVal.Msg() != "" {
			json.NewEncoder(w).Encode(errorRes{Err: errVal.Msg()})
		}
	default:
		switch err {
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
}

package streams

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"cloud.google.com/go/bigquery"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/api/googleapi"
)

const bqURL = "http://bigquery.cloud.google.com/table"

var (
	// ErrConflict indicates usage of the existing stream id
	// for the new stream.
	ErrConflict = errors.New("stream id or url already taken")

	// ErrUnauthorizedAccess indicates missing or invalid
	// credentials provided when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrWrongType indicates wrong contant type error.
	ErrWrongType = errors.New("wrong type")

	// ErrMalformedData indicates a malformed request.
	ErrMalformedData = errors.New("malformed data")

	// ErrBigQuery indicates a problem with Google Big Query API.
	ErrBigQuery = errors.New("Google Big Query error")

	// ErrInvalidBQAccess indicates unauthorized user for accessing
	// the Big Query API.
	ErrInvalidBQAccess = errors.New("wrong user role")
)

// ErrBulkConflict represents an error when saving bulk
// of Streams. It's different from other errors because
// it contains data about conflicted fields values.
type ErrBulkConflict struct {
	Message   string   `json:"message,omitempty"`
	Conflicts []string `json:"conflicts,omitempty"`
}

func (e ErrBulkConflict) Error() string {
	return e.Message
}

var _ Service = (*streamService)(nil)

// Service specifies an API that must be fullfiled by the
// domain service implementation, and all of its decorators
// (e.g. logging & metrics).
type Service interface {
	// Adds new stream to the user identified by the provided id.
	AddStream(Stream) (string, error)

	// Adds new streams via parsed csv file.
	AddBulkStreams([]Stream) error

	// Retrieves data about subset of streams given geolocation
	// coordinates, name, type, owner or price range. Data is returned
	// in the Page form. Provides check if the user is actual owner of
	// the Stream to prevent access to the real Stream URL.
	SearchStreams(string, Query) (Page, error)

	// Updates the Stream identified by the provided id.
	UpdateStream(Stream) error

	// ViewFullStream retrieves Stream data including URL by id.
	ViewFullStream(string) (Stream, error)

	// Retrieves data about the Stream identified by the id.
	// Provides check if the user is actual owner of the
	// Stream to prevent access to the real Stream URL.
	ViewStream(string, string) (Stream, error)

	// Removes the Stream identified with the provided id, that
	// belongs to the user identified by the provided id.
	RemoveStream(string, string) error
}

type streamService struct {
	streams       StreamRepository
	accessControl AccessControl
	ai            AIService
}

// NewService instantiates the domain service implementation.
func NewService(streams StreamRepository, accessControl AccessControl, ai AIService) Service {
	return streamService{
		streams:       streams,
		accessControl: accessControl,
		ai:            ai,
	}
}

func (ss streamService) AddStream(stream Stream) (string, error) {
	if stream.External {
		return ss.addBqStream(stream)
	}

	id, err := ss.streams.Save(stream)
	if err != nil {
		return "", err
	}

	stream.ID = bson.ObjectIdHex(id)

	if stream.Type == "Dataset" {
		if err := ss.ai.CreateDataset(stream); err != nil {
			return "", err
		}
	}

	if stream.Type == "Algorithm" {
		if err := ss.ai.CreateAlgorithm(stream); err != nil {
			return "", err
		}
	}

	return id, nil
}

func (ss streamService) checkAccess(owner string, access []*bigquery.AccessEntry) error {
	for _, ae := range access {
		if ae.EntityType == bigquery.UserEmailEntity && ae.Entity == owner && ae.Role == bigquery.OwnerRole {
			return nil
		}
	}
	return ErrInvalidBQAccess
}

func (ss streamService) addBqStream(stream Stream) (string, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, stream.BQ.Project)
	if err != nil {
		return "", ErrBigQuery
	}
	defer client.Close()
	ds := client.Dataset(stream.BQ.Dataset)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			switch e.Code {
			case http.StatusBadRequest, http.StatusNotFound:
				return "", ErrMalformedData
			}
		}
		return "", ErrBigQuery
	}
	if err := ss.checkAccess(stream.BQ.Email, meta.Access); err != nil {
		return "", err
	}

	bq := stream.BQ
	q := fmt.Sprintf("SELECT %s FROM `%s.%s.%s`", bq.Fields, bq.Project, bq.Dataset, bq.Table)
	id := strings.Replace(uuid.NewV4().String(), "-", "_", -1)

	// Try to create table to check if query is valid. In the
	// case of an invalid data, a Stream won't be saved.
	t := ds.Table(fmt.Sprintf("%s_%s", bq.Table, id))
	md := bigquery.TableMetadata{
		Name:           bq.Table,
		ViewQuery:      q,
		ExpirationTime: time.Now().Add(10 * time.Second),
	}

	if err := t.Create(ctx, &md); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			switch e.Code {
			case http.StatusConflict:
				return "", ErrConflict
			case http.StatusBadRequest, http.StatusNotFound:
				return "", ErrMalformedData
			}
		}
		return "", ErrBigQuery
	}

	stream.URL = fmt.Sprintf("%s/%s:%s.%s", bqURL, bq.Project, bq.Dataset, t.TableID)

	return ss.streams.Save(stream)
}

func (ss streamService) AddBulkStreams(streams []Stream) error {
	return ss.streams.SaveAll(streams)
}

func (ss streamService) SearchStreams(owner string, query Query) (Page, error) {
	partners, err := ss.accessControl.Partners(owner)
	if err != nil {
		return Page{}, err
	}
	partners = append(partners, owner)

	query.Partners = partners

	p, err := ss.streams.Search(query)
	if err != nil {
		return p, err
	}

	// Prevent sending real URL to the end user.
	for i := range p.Content {
		if p.Content[i].Owner != owner {
			p.Content[i].URL = ""
		}
	}

	return p, nil
}

func (ss streamService) UpdateStream(stream Stream) error {
	return ss.streams.Update(stream)
}

func (ss streamService) ViewFullStream(id string) (Stream, error) {
	return ss.streams.One(id)
}

func (ss streamService) ViewStream(id, owner string) (Stream, error) {
	s, err := ss.streams.One(id)
	if err != nil {
		return s, err
	}

	partners, err := ss.accessControl.Partners(owner)
	if err != nil {
		return Stream{}, err
	}
	partners = append(partners, owner)

	in := false
	for _, partner := range partners {
		if s.Owner == partner {
			in = true
			break
		}
	}
	if !in {
		return Stream{}, ErrNotFound
	}

	if s.Owner != owner {
		s.URL = ""
	}
	return s, nil
}

func (ss streamService) RemoveStream(owner string, id string) error {
	return ss.streams.Remove(owner, id)
}

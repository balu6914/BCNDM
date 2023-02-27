package streams

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/datapace/datapace/streams/groups"
	"github.com/datapace/datapace/streams/sharing"

	"github.com/datapace/datapace/errors"

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

	// ErrCreateTerms indicates terms not created
	ErrCreateTerms = errors.New("terms not created")
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

	// Adds new streams.
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

	// ExportStreams returns all streams available to the specified owner
	ExportStreams(string) ([]Stream, error)

	//Add category/subcategory endpoint
	AddCategory(string, []string) (string, error)

	//List all the categories and subcategories
	ListCatgeories(string) ([]Category, error)
}

type streamService struct {
	streams       StreamRepository
	categories    CategoryRepository
	accessControl AccessControl
	ai            AIService
	terms         TermsService
	groupsSvc     groups.Service
	sharingSvc    sharing.Service
}

// NewService instantiates the domain service implementation.
func NewService(
	streams StreamRepository,
	categories CategoryRepository,
	accessControl AccessControl,
	ai AIService,
	terms TermsService,
	groupsSvc groups.Service,
	sharingSvc sharing.Service,
) Service {
	return streamService{
		streams:       streams,
		categories:    categories,
		accessControl: accessControl,
		ai:            ai,
		terms:         terms,
		groupsSvc:     groupsSvc,
		sharingSvc:    sharingSvc,
	}
}

func (ss streamService) AddStream(stream Stream) (string, error) {
	if stream.External {
		return ss.addBqStream(stream)
	}
	now := time.Now()
	stream.StartDate = &now

	id, err := ss.streams.Save(stream)
	if err != nil {
		return "", err
	}

	stream.ID = id

	if stream.Type == "Dataset" {
		if err := ss.ai.CreateDataset(stream); err != nil {
			return "", errors.Wrap(err, ss.RemoveStream(stream.Owner, id))
		}
	}

	if stream.Type == "Algorithm" {
		if err := ss.ai.CreateAlgorithm(stream); err != nil {
			return "", errors.Wrap(err, ss.RemoveStream(stream.Owner, id))
		}
	}

	if err := ss.terms.CreateTerms(stream); err != nil {
		return "", errors.Wrap(ErrCreateTerms, errors.Wrap(err, ss.RemoveStream(stream.Owner, id)))
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

	ids := ss.resolveSharedStreams(owner)
	query.Shared = ids

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

	if s.Visibility == Protected && s.Owner != owner {
		if !ss.allowedToAccess(owner, s) {
			return Stream{}, ErrNotFound
		}
	}

	if s.Owner != owner {
		s.URL = ""
	}
	return s, nil
}

func (ss streamService) RemoveStream(owner string, id string) error {
	if err := ss.streams.Remove(owner, id); err != nil {
		return err
	}
	_ = ss.sharingSvc.DeleteSharing(owner, id)
	return nil
}

func (ss streamService) ExportStreams(owner string) ([]Stream, error) {
	partners, err := ss.accessControl.Partners(owner)
	if err != nil {
		return []Stream{}, err
	}
	partners = append(partners, owner)
	sharedIds := ss.resolveSharedStreams(owner)
	q := Query{
		Name:       "",
		Partners:   partners,
		Shared:     sharedIds,
		StreamType: "",
		Coords:     nil,
		Page:       0,
		Limit:      math.MaxInt32,
		MinPrice:   nil,
		MaxPrice:   nil,
	}
	p, err := ss.streams.Search(q)
	if err != nil {
		return p.Content, err
	}
	return p.Content, nil
}

func (ss streamService) resolveSharedStreams(userId string) map[string]bool {
	groupIds, _ := ss.groupsSvc.GetUserGroups(userId)
	sharings, _ := ss.sharingSvc.GetSharings(userId, groupIds)
	streamIds := make(map[string]bool)
	for _, s := range sharings {
		streamIds[string(s.StreamId)] = true
	}
	return streamIds
}

func (ss streamService) allowedToAccess(userId string, stream Stream) bool {
	if ss.arePartners(stream.Owner, userId) {
		return true
	}
	if ss.isShared(stream.ID, userId) {
		return true
	}
	return false
}

func (ss streamService) arePartners(owner, userId string) bool {

	partners, err := ss.accessControl.Partners(userId)
	if err != nil {
		return false
	}

	in := false
	for _, partner := range partners {
		if owner == partner {
			in = true
			break
		}
	}
	return in
}

func (ss streamService) isShared(id string, userId string) bool {
	ids := ss.resolveSharedStreams(userId)
	return ids[id]
}

func (ss streamService) AddCategory(categoryName string, subCategoryNames []string) (string, error) {
	id, err := ss.categories.Save(categoryName, subCategoryNames)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (ss streamService) ListCatgeories(key string) ([]Category, error) {
	categoriesList, err := ss.categories.List(key)
	if err != nil {
		return nil, err
	}

	return categoriesList, nil
}

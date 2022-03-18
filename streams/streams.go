package streams

import (
	"net/url"
	"strconv"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

var _ Service = (*streamService)(nil)

const (
	// Public streams are visible for all users
	Public Visibility = "public"
	// Protected streams are visibible for users with access
	Protected Visibility = "protected"
	// Private streams are only visible to owner
	Private Visibility = "private"
)

// Visibility of streams
type Visibility string

// Location represents Stream location to enable geo
// search streams. Official MongoDB docs could be found here
// http://docs.mongoengine.org/guide/querying.html#geo-queries
type Location struct {
	Type string `json:"type,omitempty"`
	// Coordinates represent longitude and latitude. It's represented this
	// way to match the way MongoDB represents geo data.
	Coordinates [2]float64 `json:"coordinates,omitempty"`
}

// BigQuery holds Big Query related data.
type BigQuery struct {
	// Email represents Gmail address of the owner. It can be either Email
	// or ContactEmail of the owner.
	Email   string `json:"email,omitempty"`
	Project string `json:"project,omitempty"`
	Dataset string `json:"dataset,omitempty"`
	Table   string `json:"table,omitempty"`
	Fields  string `json:"fields,omitempty"`
}

// Validate provides basic checks of parameters related to the Big Query.
func (bq BigQuery) Validate() bool {
	return bq.Email != "" &&
		bq.Project != "" &&
		bq.Dataset != "" &&
		bq.Fields != "" &&
		bq.Table != ""
}

// Stream represents data stream to be exchanged through platform.
type Stream struct {
	Owner       string                 `json:"owner,omitempty"`
	ID          string                 `json:"id,omitempty"`
	Visibility  Visibility             `json:"visibility,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	Snippet     string                 `json:"snippet,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Price       uint64                 `json:"price,omitempty"`
	Location    Location               `json:"location,omitempty"`
	Terms       string                 `json:"terms,omitempty"`
	External    bool                   `json:"external,omitempty"`
	BQ          BigQuery               `json:"bq,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Attributes returns auth.Resource attributes.
func (s Stream) Attributes() map[string]string {
	return map[string]string{
		"id":          s.ID,
		"ownerID":     s.Owner,
		"visibility":  string(s.Visibility),
		"name":        s.Name,
		"type":        s.Type,
		"description": s.Description,
		"snippet":     s.Snippet,
		"url":         s.URL,
		"price":       strconv.FormatUint(s.Price, 10),
		"terms":       s.Terms,
		"external":    strconv.FormatBool(s.External),
	}
}

// ResourceType returns auth.Resource type.
func (s Stream) ResourceType() string {
	return "stream"
}

// Page represents paged result for list response.
type Page struct {
	Page    uint64   `json:"page"`
	Limit   uint64   `json:"limit"`
	Total   uint64   `json:"total"`
	Content []Stream `json:"content"`
}

const (
	maxNameLength        = 256
	maxTypeLength        = 32
	maxDescriptionLength = 2048
	maxSnippetLength     = 2048
	maxURLLength         = 2048
	minPrice             = 0
	minLongitude         = -180
	maxLongitude         = 180
	minLatitude          = -90
	maxLatitude          = 90
	maxMetadataLength    = 2048
)

// Validate returns an error if stream representation is invalid.
func (s *Stream) Validate() error {
	if s.Name == "" || (len(s.Name) > maxNameLength) ||
		s.Type == "" || (len(s.Type) > maxTypeLength) ||
		s.Description == "" || (len(s.Description) > maxDescriptionLength) ||
		(len(s.Snippet) > maxSnippetLength) ||
		s.Price <= minPrice ||
		s.Location.Coordinates[0] < minLongitude || s.Location.Coordinates[0] > maxLongitude ||
		s.Location.Coordinates[1] < minLatitude || s.Location.Coordinates[1] > maxLatitude ||
		// // TODO: Add Metadata length validation
		s.Visibility != Public && s.Visibility != Protected && s.Visibility != Private {
		return ErrMalformedData
	}

	var err error
	if s.Terms, err = url.PathUnescape(s.Terms); err != nil {
		return ErrMalformedData
	}
	if s.URL, err = url.PathUnescape(s.URL); err != nil {
		return ErrMalformedData
	}
	if !s.External && len(s.URL) > maxURLLength {
		return ErrMalformedData
	}

	if !s.External && (!govalidator.IsURL(s.Terms) || len(s.Terms) > maxURLLength) {
		return ErrMalformedData
	}

	if s.ID != "" && !bson.IsObjectIdHex(s.Owner) {
		return ErrMalformedData
	}

	if s.Owner != "" && !bson.IsObjectIdHex(s.Owner) {
		return ErrMalformedData
	}

	if s.External && !s.BQ.Validate() {
		return ErrMalformedData
	}

	return nil
}

// StreamRepository specifies a stream persistence API.
type StreamRepository interface {
	// Save persists the stream. A non-nil error is returned to indicate
	// operation failure.
	Save(Stream) (string, error)

	// Save persists an array of streams. A non-nil error is returned to
	// indicate operation failure.
	SaveAll([]Stream) error

	// Search for streams by specified query parameters.
	Search(Query) (Page, error)

	// Update performs an update of an existing stream. A non-nil error is
	// returned to indicate operation failure.
	Update(Stream) error

	// One retrieves a stream by its unique identifier (i.e. id).
	One(string) (Stream, error)

	// Removes the stream having the provided identifier, that is owned
	// by the specified user.
	Remove(string, string) error
}

package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"github.com/datapace/datapace/streams"
	"github.com/datapace/datapace/subscriptions/accessv2"
	"github.com/datapace/datapace/subscriptions/sharing"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"

	authproto "github.com/datapace/datapace/proto/auth"

	"cloud.google.com/go/bigquery"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/api/googleapi"
)

const (
	bqURL       = "http://bigquery.cloud.google.com/table"
	gmailSuffix = "@gmail.com"
)

var (
	// ErrConflict indicates usage of the existing stream id for the new stream.
	ErrConflict = errors.New("Subscription ID already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrFailedTransfer indicates that token transfer failed.
	ErrFailedTransfer = errors.New("failed to transfer tokens")

	// ErrFailedCreateSub indicates that creation of subscription failed.
	ErrFailedCreateSub = errors.New("failed to create subscription")

	// ErrStreamAccess indicates that the stream requires the access to be approved to create a subscription.
	ErrStreamAccess = errors.New("the stream needs access approval by the stream owner")

	// ErrNotEnoughTokens indicates that spender doesn't have enough tokens.
	ErrNotEnoughTokens = errors.New("not enough tokens")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// AddSubscription subscribes to a stream making a new Subscription.
	AddSubscription(string, string, Subscription) (Subscription, error)

	// SearchSubscriptions searches subscriptions by the query.
	SearchSubscriptions(Query) (Page, error)

	// ViewSubscription retrieves subscription by the given ID.
	ViewSubscription(string, string) (Subscription, error)

	// ViewSubByUserAndStream retrieves subscription by the given user ID and
	// stream ID.
	ViewSubByUserAndStream(string, string) (Subscription, error)

	// Report generates PDF report.
	Report(Query, string) ([]byte, error)
}

var _ Service = (*subscriptionsService)(nil)

type subscriptionsService struct {
	auth          authproto.AuthServiceClient
	subscriptions SubscriptionRepository
	streams       StreamsService
	proxy         Proxy
	transactions  TransactionsService
	sharingSvc    sharing.Service
	accessV2Svc   accessv2.Service
}

// New instantiates the domain service implementation.
func New(auth authproto.AuthServiceClient, subs SubscriptionRepository, streams StreamsService, proxy Proxy, transactions TransactionsService, sharingSvc sharing.Service, accessV2Svc accessv2.Service) Service {
	return &subscriptionsService{
		auth:          auth,
		subscriptions: subs,
		streams:       streams,
		proxy:         proxy,
		transactions:  transactions,
		sharingSvc:    sharingSvc,
		accessV2Svc:   accessV2Svc,
	}
}

func (ss subscriptionsService) checkEmail(userEmail authproto.UserEmail) (string, error) {
	email := strings.ToLower(userEmail.Email)
	contactEmail := strings.ToLower(userEmail.ContactEmail)
	if strings.HasSuffix(email, gmailSuffix) {
		return userEmail.Email, nil
	}
	if strings.HasSuffix(contactEmail, gmailSuffix) {
		return userEmail.ContactEmail, nil
	}
	return "", ErrMalformedEntity
}

func (ss subscriptionsService) createDataset(client *bigquery.Client, userToken, datasetID string, stream Stream) (*bigquery.Dataset, error) {
	if client == nil {
		return nil, ErrFailedCreateSub
	}

	ds := client.Dataset(datasetID)
	email, err := ss.auth.Email(context.Background(), &authproto.Token{Value: userToken})
	if err != nil {
		return nil, err
	}

	gmail, err := ss.checkEmail(*email)
	if err != nil {
		return nil, err
	}

	ae := &bigquery.AccessEntry{
		Role:       bigquery.ReaderRole,
		EntityType: bigquery.UserEmailEntity,
		Entity:     gmail,
	}

	meta := bigquery.DatasetMetadata{
		Access: []*bigquery.AccessEntry{ae},
	}

	if err = ds.Create(context.Background(), &meta); err != nil {
		return nil, err
	}

	return ds, nil
}

func (ss subscriptionsService) createTable(ds *bigquery.Dataset, stream Stream, expire time.Time, tableID string) error {
	q := fmt.Sprintf("SELECT %s FROM `%s.%s.%s`", stream.Fields, stream.Project, stream.Dataset, stream.Table)
	t := ds.Table(tableID)
	tableMeta := bigquery.TableMetadata{
		Name:           stream.Table,
		ViewQuery:      q,
		ExpirationTime: expire,
	}

	if err := t.Create(context.Background(), &tableMeta); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			switch e.Code {
			case http.StatusConflict:
				return ErrConflict
			case http.StatusBadRequest, http.StatusNotFound:
				return ErrMalformedEntity
			}
		}
		return err
	}

	return nil
}

func (ss subscriptionsService) createViewDataset(email, datasetID string, client *bigquery.Client) error {
	ds := client.Dataset(datasetID)
	ae := &bigquery.AccessEntry{
		Role:       bigquery.ReaderRole,
		EntityType: bigquery.UserEmailEntity,
		Entity:     email,
	}

	meta := bigquery.DatasetMetadata{
		Access: []*bigquery.AccessEntry{ae},
	}

	return ds.Create(context.Background(), &meta)
}

func (ss subscriptionsService) createView(datasetID, viewID string, stream Stream, expire time.Time, client *bigquery.Client) error {
	q := fmt.Sprintf("SELECT %s FROM `%s.%s.%s`", stream.Fields, stream.Project, stream.Dataset, stream.Table)
	t := client.Dataset(datasetID).Table(viewID)
	tableMeta := bigquery.TableMetadata{
		Name:           stream.Table,
		ViewQuery:      q,
		ExpirationTime: expire,
	}

	ctx := context.Background()
	if err := t.Create(ctx, &tableMeta); err != nil {
		client.Dataset(datasetID).Delete(ctx)
		if e, ok := err.(*googleapi.Error); ok {
			switch e.Code {
			case http.StatusConflict:
				return ErrConflict
			case http.StatusBadRequest, http.StatusNotFound:
				return ErrMalformedEntity
			}
		}
		return err
	}

	return nil
}

func (ss subscriptionsService) setViewAccess(srcID, dstID, tableID string, client *bigquery.Client) error {
	ctx := context.Background()
	ds := client.Dataset(srcID)
	md, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}

	mdUpdate := []*bigquery.AccessEntry{}
	for _, v := range md.Access {
		if v.View != nil {
			_, err := client.Dataset(v.View.DatasetID).Table(v.View.TableID).Metadata(ctx)
			// TODO Add error check and possibly much more complex handling.
			if err == nil {
				mdUpdate = append(mdUpdate, v)
			}
			continue
		}
		mdUpdate = append(mdUpdate, v)
	}

	t := client.Dataset(dstID).Table(tableID)
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(mdUpdate, &bigquery.AccessEntry{
			EntityType: bigquery.ViewEntity,
			View:       t},
		),
	}

	if _, err := ds.Update(ctx, update, md.ETag); err != nil {
		return err
	}

	return nil
}

func (ss subscriptionsService) createBQ(url *string, token string, stream Stream, client *bigquery.Client, expire time.Time) (*bigquery.Dataset, error) {
	email, err := ss.auth.Email(context.Background(), &authproto.Token{Value: token})
	if err != nil {
		return nil, err
	}

	gmail, err := ss.checkEmail(*email)
	if err != nil {
		return nil, err
	}

	id := strings.Replace(uuid.NewV4().String(), "-", "_", -1)
	datasetID := fmt.Sprintf("%s_%s", stream.Dataset, id)
	viewID := fmt.Sprintf("%s_%s", stream.Table, id)
	// Create dataset for authorized view.
	if err := ss.createViewDataset(gmail, datasetID, client); err != nil {
		return nil, err
	}
	// Create authorized view.
	if err := ss.createView(datasetID, viewID, stream, expire, client); err != nil {
		return nil, err
	}
	// Authorize view.
	if err := ss.setViewAccess(stream.Dataset, datasetID, viewID, client); err != nil {
		return nil, err
	}

	*url = fmt.Sprintf("%s/%s:%s.%s", bqURL, stream.Project, datasetID, viewID)
	return client.Dataset(datasetID), nil
}

func (ss subscriptionsService) AddSubscription(userID, token string, sub Subscription) (Subscription, error) {
	sub.UserID = userID
	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))

	stream, err := ss.streams.One(sub.StreamID)
	if err != nil {
		return Subscription{}, ErrNotFound
	}

	if userID != stream.Owner && stream.Visibility == string(streams.Private) {
		return Subscription{}, fmt.Errorf("%w: non-own private stream subscription is forbidden", ErrFailedCreateSub)
	}

	err = ss.checkStreamAccess(userID, stream)
	if err != nil {
		return Subscription{}, err
	}

	sub.StreamOwner = stream.Owner
	sub.StreamName = stream.Name
	sub.StreamPrice = stream.Price
	url := stream.URL

	var ds *bigquery.Dataset
	if stream.External {
		client, err := bigquery.NewClient(context.Background(), stream.Project)
		if err != nil {
			return Subscription{}, err
		}
		defer client.Close()
		ds, err = ss.createBQ(&url, token, stream, client, sub.EndDate)
		if err != nil {
			return Subscription{}, err
		}
	}

	hash, err := ss.proxy.Register(sub.Hours, url)
	if err != nil {
		if ds != nil {
			// Ignore if deletion fails, table will be removed after expiration anyway.
			ds.Delete(context.Background())
		}
		return Subscription{}, err
	}

	sub.StreamURL = hash

	id, err := ss.subscriptions.Save(sub)
	if err != nil {
		if ds != nil {
			ds.Delete(context.Background())
		}

		if err == ErrConflict {
			return Subscription{}, ErrConflict
		}

		return Subscription{}, ErrFailedCreateSub
	}

	// do not invoke transactions if the stream is shared to the current user
	if !ss.isStreamSharedTo(stream.ID, userID) {
		if err := ss.transactions.Transfer(stream.ID, userID, stream.Owner, stream.Price*sub.Hours); err != nil {
			if ds != nil {
				ds.Delete(context.Background())
			}
			if err == ErrNotEnoughTokens {
				return Subscription{}, err
			}
			return Subscription{}, ErrFailedTransfer
		}
	}

	ss.subscriptions.Activate(id)
	sub.ID = bson.ObjectIdHex(id)

	return sub, nil
}

func (ss subscriptionsService) SearchSubscriptions(q Query) (Page, error) {
	return ss.subscriptions.Search(q)
}

func (ss subscriptionsService) ViewSubscription(userID, subID string) (Subscription, error) {
	sub, err := ss.subscriptions.One(subID)
	if err != nil {
		return Subscription{}, err
	}

	if userID != sub.StreamOwner && userID != sub.UserID {
		return Subscription{}, ErrNotFound
	}

	return sub, nil
}

func (ss subscriptionsService) ViewSubByUserAndStream(userID, streamID string) (Subscription, error) {
	return ss.subscriptions.OneByUserAndStream(userID, streamID)
}

func (ss subscriptionsService) isStreamSharedTo(streamId, rcvUserId string) bool {
	if ss.sharingSvc != nil {
		sharings, err := ss.sharingSvc.GetSharings(rcvUserId, []string{})
		if err != nil {
			return false
		}
		for _, s := range sharings {
			if string(s.StreamId) == streamId {
				return true
			}
		}
	}
	return false
}

func (ss subscriptionsService) checkStreamAccess(userId string, stream Stream) (err error) {
	if ss.accessV2Svc != nil && userId != stream.Owner && stream.AccessType == streams.AccessTypeProtected {
		k := accessv2.Key{
			ConsumerId: userId,
			ProviderId: stream.Owner,
			ProductId:  stream.ID,
		}
		var a accessv2.Access
		a, err = ss.accessV2Svc.Get(context.TODO(), k)
		switch {
		case errors.Is(err, accessv2.ErrNotFound):
			err = fmt.Errorf("%w: access is not requested", ErrStreamAccess)
		case errors.Is(err, accessv2.ErrNotAvailable):
			err = nil // access service is not available, do nothing
		case err != nil:
			err = fmt.Errorf("%w: failed to check whether access is granted: %s", ErrStreamAccess, err.Error())
		case a.State != accessv2.StateApproved:
			err = fmt.Errorf("%w: access request state: %s", ErrStreamAccess, a.State.String())
		}
	}
	return
}

package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"datapace"
	"net/http"
	"strings"
	"time"

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

	// ErrNotEnoughTokens indicates that spender doesn't have enough tokens.
	ErrNotEnoughTokens = errors.New("not enough tokens")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// AddSubscription subscribes to a stream making a new Subscription.
	AddSubscription(string, string, Subscription) (string, error)

	// SearchSubscriptions searches subscriptions by the query.
	SearchSubscriptions(Query) (Page, error)

	// ViewSubscription retrieves subscription by the given ID.
	ViewSubscription(string, string) (Subscription, error)
}

var _ Service = (*subscriptionsService)(nil)

type subscriptionsService struct {
	auth          datapace.AuthServiceClient
	subscriptions SubscriptionRepository
	streams       StreamsService
	proxy         Proxy
	transactions  TransactionsService
}

// New instantiates the domain service implementation.
func New(auth datapace.AuthServiceClient, subs SubscriptionRepository, streams StreamsService, proxy Proxy, transactions TransactionsService) Service {
	return &subscriptionsService{
		auth:          auth,
		subscriptions: subs,
		streams:       streams,
		proxy:         proxy,
		transactions:  transactions,
	}
}

func (ss subscriptionsService) checkEmail(userEmail datapace.UserEmail) (string, error) {
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
	email, err := ss.auth.Email(context.Background(), &datapace.Token{Value: userToken})
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

func (ss subscriptionsService) createVewDataset(email, datasetID string, client *bigquery.Client) error {
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
			_, err := v.View.Metadata(ctx)
			// TODO Add error check and possibly much more complex handling.
			if err != nil {
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
	email, err := ss.auth.Email(context.Background(), &monetasa.Token{Value: token})
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
	if err := ss.createVewDataset(gmail, datasetID, client); err != nil {
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

	viewURL := fmt.Sprintf("%s/%s:%s.%s", bqURL, stream.Project, datasetID, viewID)
	url = &viewURL
	return client.Dataset(datasetID), nil
}

func (ss subscriptionsService) AddSubscription(userID, token string, sub Subscription) (string, error) {
	sub.UserID = userID
	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))

	stream, err := ss.streams.One(sub.StreamID)
	if err != nil {
		return "", ErrNotFound
	}

	sub.StreamOwner = stream.Owner
	sub.StreamName = stream.Name
	sub.StreamPrice = stream.Price
	url := stream.URL

	var ds *bigquery.Dataset
	if stream.External {
		client, err := bigquery.NewClient(context.Background(), stream.Project)
		if err != nil {
			return "", err
		}
		defer client.Close()
		ds, err = ss.createBQ(&url, token, stream, client, sub.EndDate)
		if err != nil {
			return "", err
		}
	}

	hash, err := ss.proxy.Register(sub.Hours, url)
	if err != nil {
		if ds != nil {
			// Ignore if deletion fails, table will be removed after expiration anyway.
			ds.Delete(context.Background())
		}
		return "", err
	}

	sub.StreamURL = hash

	id, err := ss.subscriptions.Save(sub)
	if err != nil {
		if ds != nil {
			ds.Delete(context.Background())
		}

		if err == ErrConflict {
			return "", ErrConflict
		}

		return "", ErrFailedCreateSub
	}

	if err := ss.transactions.Transfer(stream.ID, userID, stream.Owner, stream.Price*sub.Hours); err != nil {
		if ds != nil {
			ds.Delete(context.Background())
		}

		if err == ErrNotEnoughTokens {
			return "", err
		}

		return "", ErrFailedTransfer
	}

	ss.subscriptions.Activate(id)

	return id, nil
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

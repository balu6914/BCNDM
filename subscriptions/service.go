package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/api/googleapi"
)

const bqURL = "http://bigquery.cloud.google.com/table"

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
	AddSubscription(string, Subscription) (string, error)

	// SearchSubscriptions searches subscriptions by the query.
	SearchSubscriptions(Query) (Page, error)

	// ViewSubscription retrieves subscription by the given ID.
	ViewSubscription(string, string) (Subscription, error)
}

var _ Service = (*subscriptionsService)(nil)

type subscriptionsService struct {
	subscriptions SubscriptionRepository
	streams       StreamsService
	proxy         Proxy
	transactions  TransactionsService
}

// New instantiates the domain service implementation.
func New(subs SubscriptionRepository, streams StreamsService, proxy Proxy, transactions TransactionsService) Service {
	return &subscriptionsService{
		subscriptions: subs,
		streams:       streams,
		proxy:         proxy,
		transactions:  transactions,
	}
}

func (ss subscriptionsService) createBqView(stream Stream, end time.Time) (*bigquery.Table, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, stream.Project)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	q := fmt.Sprintf("SELECT %s FROM `%s.%s.%s`", stream.Fields, stream.Project, stream.Dataset, stream.Table)
	id := strings.Replace(uuid.NewV4().String(), "-", "_", -1)
	id = fmt.Sprintf("%s_%s", stream.Table, id)

	ds := client.Dataset(stream.Dataset)
	t := ds.Table(id)
	md := bigquery.TableMetadata{
		Name:           stream.Table,
		ViewQuery:      q,
		ExpirationTime: end,
	}

	if err := t.Create(ctx, &md); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			switch e.Code {
			case http.StatusConflict:
				return nil, ErrConflict
			case http.StatusBadRequest:
				return nil, ErrMalformedEntity
			}
		}
		return nil, err
	}

	return t, nil
}

func (ss subscriptionsService) AddSubscription(userID string, sub Subscription) (string, error) {
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

	var table *bigquery.Table
	if stream.External {
		table, err := ss.createBqView(stream, sub.EndDate)
		if err != nil {
			return "", err
		}

		url = fmt.Sprintf("%s/%s:%s.%s", bqURL, stream.Project, stream.Dataset, table.TableID)
	}

	hash, err := ss.proxy.Register(sub.Hours, url)
	if err != nil {
		if table != nil {
			// Ignore if deletion fails, table will be removed after expiration anyway.
			table.Delete(context.Background())
		}
		return "", err
	}

	sub.StreamURL = hash

	id, err := ss.subscriptions.Save(sub)
	if err != nil {
		if table != nil {
			table.Delete(context.Background())
		}

		if err == ErrConflict {
			return "", ErrConflict
		}

		return "", ErrFailedCreateSub
	}

	if err := ss.transactions.Transfer(stream.ID, userID, stream.Owner, stream.Price*sub.Hours); err != nil {
		if table != nil {
			table.Delete(context.Background())
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

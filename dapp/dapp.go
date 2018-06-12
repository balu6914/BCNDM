package dapp

import (
	"strings"
	"time"
)

var _ Service = (*dappService)(nil)

type dappService struct {
	streams       StreamRepository
	subscriptions SubscriptionsRepository
}

// New instantiates the domain service implementation.
func New(streams StreamRepository, subs SubscriptionsRepository) Service {
	return &dappService{
		streams:       streams,
		subscriptions: subs,
	}
}

func (ds *dappService) authorize(user, id string) (bool, error) {
	s, err := ds.streams.One(id)
	if err != nil {
		return false, err
	}

	if strings.Compare(user, s.Owner) != 0 {
		return false, ErrUnauthorizedAccess
	}

	return true, nil
}

func (ds *dappService) AddStream(stream Stream) (string, error) {
	return ds.streams.Save(stream)
}

func (ds *dappService) AddBulkStream(streams []Stream) error {
	if len(streams) < 1 {
		return ErrMalformedData
	}
	for _, stream := range streams {
		if _, err := ds.streams.Save(stream); err != nil {
			return err
		}
	}
	return nil
}

func (ds *dappService) UpdateStream(owner, id string, stream Stream) error {
	if _, err := ds.authorize(owner, id); err != nil {
		return err
	}
	return ds.streams.Update(id, stream)
}

func (ds *dappService) ViewStream(id string) (Stream, error) {
	return ds.streams.One(id)
}

func (ds *dappService) SearchStreams(coords [][]float64) ([]Stream, error) {
	return ds.streams.Search(coords)
}

func (ds *dappService) RemoveStream(user string, id string) error {
	if _, err := ds.authorize(user, id); err != nil {
		return err
	}

	return ds.streams.Remove(id)
}

func (ds *dappService) CreateSubscription(sub Subscription) error {
	stream, err := ds.streams.One(sub.StreamID)
	if err != nil {
		return err
	}

	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))
	sub.StreamData.Coordinates[0] = stream.Location.Coordinates[0]
	sub.StreamData.Coordinates[1] = stream.Location.Coordinates[1]
	sub.StreamData.Name = stream.Name
	sub.StreamData.Price = stream.Price

	return ds.subscriptions.Create(sub)
}

func (ds *dappService) GetSubscriptions(userID string) ([]Subscription, error) {
	return ds.subscriptions.Read(userID)
}

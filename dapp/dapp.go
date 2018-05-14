package dapp

import (
	"strings"
)

var _ Service = (*dappService)(nil)

type dappService struct {
	streams StreamRepository
}

// New instantiates the domain service implementation.
func New(streams StreamRepository) Service {
	return &dappService{
		streams: streams,
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

func (ds *dappService) AddStream(owner string, stream Stream) (string, error) {
	stream.Owner = owner
	return ds.streams.Save(stream)
}

func (ds *dappService) AddBulkStream(streams []Stream) error {
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

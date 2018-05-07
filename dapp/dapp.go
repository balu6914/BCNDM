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

	if strings.Compare(user, s.User) != 0 {
		return false, ErrUnauthorizedAccess
	}

	return true, nil
}

func (ds *dappService) AddStream(user string, stream Stream) (string, error) {
	stream.User = user

	return ds.streams.Save(stream)
}

func (ds *dappService) UpdateStream(user string, id string, stream Stream) error {
	if _, err := ds.authorize(user, id); err != nil {
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

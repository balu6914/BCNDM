package dapp

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

func (ds *dappService) AddStream(key string, stream Stream) (string, error) {
	return ds.streams.Save(stream)
}

func (ds *dappService) UpdateStream(key string, id string, stream Stream) error {
	return ds.streams.Update(id, stream)
}

func (ds *dappService) ViewStream(key string, id string) (Stream, error) {
	return ds.streams.One(id)
}

func (ds *dappService) SearchStreams(coords [][]float64) ([]Stream, error) {
	return ds.streams.Search(coords)
}

func (ds *dappService) RemoveStream(key string, id string) error {
	return ds.streams.Remove(id)
}

// Package terms contains terms implementation for streams service.
package terms

import (
	"context"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/datapace/datapace/streams"
	"time"
)

var _ streams.TermsService = (*termsService)(nil)

type termsService struct {
	client termsproto.TermsServiceClient
}

// New returns new terms instance.
func New(client termsproto.TermsServiceClient) streams.TermsService {
	return termsService{client: client}
}

func (ts termsService) CreateTerms(s streams.Stream) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := ts.client.CreateTerms(ctx, &termsproto.Terms{
		StreamId: s.ID,
		Url:      s.Terms,
	})
	if err != nil {
		return err
	}
	return nil
}

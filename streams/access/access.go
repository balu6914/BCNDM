// Package access contains access control implementation for streams service.
package access

import (
	"context"
	"time"

	"github.com/datapace"

	"github.com/datapace/streams"
)

var _ streams.AccessControl = (*accessControl)(nil)

type accessControl struct {
	client datapace.AccessServiceClient
}

// New returns new access control instance.
func New(client datapace.AccessServiceClient) streams.AccessControl {
	return accessControl{client: client}
}

func (ac accessControl) Partners(id string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list, err := ac.client.Partners(ctx, &datapace.ID{Value: id})
	if err != nil {
		return nil, err
	}

	return list.Value, nil
}

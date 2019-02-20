// Package access contains access control implementation for streams service.
package access

import (
	"context"
	"datapace"
	"datapace/streams"
	"time"
)

var _ streams.AccessControl = (*accessControl)(nil)

type accessControl struct {
	client datapace.AuthServiceClient
}

// New returns new access control instance.
func New(client datapace.AuthServiceClient) streams.AccessControl {
	return accessControl{client: client}
}

func (ac accessControl) Partners(id string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list, err := ac.client.Partners(ctx, &datapace.UserID{Value: id})
	if err != nil {
		return nil, err
	}

	return list.Value, nil
}

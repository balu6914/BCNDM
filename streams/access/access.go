// Package access contains access control implementation for streams service.
package access

import (
	"context"
	"time"

	accessproto "github.com/datapace/datapace/proto/access"
	commonproto "github.com/datapace/datapace/proto/common"
	"github.com/datapace/datapace/streams"
)

var _ streams.AccessControl = (*accessControl)(nil)

type accessControl struct {
	client accessproto.AccessServiceClient
}

// New returns new access control instance.
func New(client accessproto.AccessServiceClient) streams.AccessControl {
	return accessControl{client: client}
}

func (ac accessControl) Partners(id string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list, err := ac.client.Partners(ctx, &commonproto.ID{Value: id})
	if err != nil {
		return nil, err
	}

	return list.Value, nil
}

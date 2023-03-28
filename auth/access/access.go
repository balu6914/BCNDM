package access

import (
	"context"
	"time"

	"github.com/datapace/datapace/auth"
	accessproto "github.com/datapace/datapace/proto/access"
	commonproto "github.com/datapace/datapace/proto/common"
)

var _ auth.AccessControl = (*accessControl)(nil)

type accessControl struct {
	client accessproto.AccessServiceClient
}

// New returns new access control instance.
func New(client accessproto.AccessServiceClient) auth.AccessControl {
	return accessControl{client: client}
}

func (ac accessControl) PotentialPartners(id string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := ac.client.PotentialPartners(ctx, &commonproto.ID{Value: id})
	if err != nil {
		return []string{}, err
	}

	return res.GetValue(), nil
}

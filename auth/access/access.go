package access

import (
	"context"
	"datapace"
	"datapace/auth"
	"time"
)

var _ auth.AccessControl = (*accessControl)(nil)

type accessControl struct {
	client datapace.AccessServiceClient
}

// New returns new access control instance.
func New(client datapace.AccessServiceClient) auth.AccessControl {
	return accessControl{client: client}
}

func (ac accessControl) PotentialPartners(id string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := ac.client.PotentialPartners(ctx, &datapace.ID{Value: id})
	if err != nil {
		return []string{}, err
	}

	return res.GetValue(), nil
}

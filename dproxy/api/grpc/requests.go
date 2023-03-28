package grpc

import "github.com/datapace/datapace/dproxy/persistence"

type listRequest struct {
	query persistence.Query
}

func (req listRequest) validate() error {
	return nil
}

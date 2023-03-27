package grpc

import (
	"github.com/datapace/datapace/dproxy/persistence"
)

type listResponse struct {
	events []persistence.Event
	err    error
}

package grpc_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	transactionsproto "github.com/datapace/datapace/proto/transactions"
	"github.com/datapace/datapace/transactions"
	grpcapi "github.com/datapace/datapace/transactions/api/grpc"
	"github.com/datapace/datapace/transactions/mocks"

	"google.golang.org/grpc"
)

const (
	port            = 8000
	id1             = "5281b83afbb7f35cb62d0834"
	id2             = "5281b83afbb7f35cb62d0835"
	secret          = "secret"
	balance         = 42
	remainingTokens = 100
)

var svc transactions.Service

func TestMain(m *testing.M) {
	startServer()
	code := m.Run()
	os.Exit(code)
}

func startServer() {
	svc = newService()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	server := grpc.NewServer()
	transactionsproto.RegisterTransactionsServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

func newService() transactions.Service {
	ur := mocks.NewUserRepository(map[string]string{
		id1: secret,
		id2: secret,
	})
	tl := mocks.NewTokenLedger(map[string]uint64{
		id1: balance,
		id2: balance,
	}, remainingTokens)
	cl := mocks.NewContractLedger()
	cr := mocks.NewContractRepository()
	streams := mocks.NewStreamsService(map[string]transactions.Stream{})

	return transactions.New(ur, tl, cl, cr, streams)
}

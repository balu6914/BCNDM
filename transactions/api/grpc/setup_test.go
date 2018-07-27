package grpc_test

import (
	"fmt"
	"monetasa"
	"monetasa/transactions"
	grpcapi "monetasa/transactions/api/grpc"
	"monetasa/transactions/mocks"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
)

const (
	port            = 8080
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
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	monetasa.RegisterTransactionsServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

func newService() transactions.Service {
	repo := mocks.NewUserRepository(map[string]string{
		id1: secret,
		id2: secret,
	})
	bn := mocks.NewBlockchainNetwork(map[string]uint64{
		id1: balance,
		id2: balance,
	}, remainingTokens)

	return transactions.New(repo, bn)
}

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
	port    = 8080
	id      = "id"
	secret  = "secret"
	balance = 42
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
	bn := mocks.NewBlockchainNetwork(map[string]uint64{
		id: balance,
	})
	return transactions.New(bn)
}

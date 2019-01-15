package grpc_test

import (
	"fmt"
	"datapace"
	"datapace/streams"
	grpcapi "datapace/streams/api/grpc"
	"datapace/streams/mocks"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
)

const (
	port            = 8080
	secret          = "secret"
	owner           = "owner"
	balance         = 42
	remainingTokens = 100
)

var svc streams.Service

func TestMain(m *testing.M) {
	startServer()
	code := m.Run()
	os.Exit(code)
}

func startServer() {
	svc = newService()
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	datapace.RegisterStreamsServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

func newService() streams.Service {
	repo := mocks.NewStreamRepository()

	return streams.NewService(repo)
}

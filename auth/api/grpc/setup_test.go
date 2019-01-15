package grpc_test

import (
	"fmt"
	"datapace"
	"datapace/auth"
	grpcapi "datapace/auth/api/grpc"
	"datapace/auth/mocks"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
)

const port = 8081

var svc auth.Service

func TestMain(m *testing.M) {
	svc = newService()
	startGRPCServer(svc, port)
	code := m.Run()
	os.Exit(code)
}

func newService() auth.Service {
	repo := mocks.NewUserRepository()
	hasher := mocks.NewHasher()
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()

	return auth.New(repo, hasher, idp, ts)
}

func startGRPCServer(svc auth.Service, port int) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	datapace.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

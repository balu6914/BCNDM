package grpc_test

import (
	"datapace"
	"datapace/auth"
	grpcapi "datapace/auth/api/grpc"
	"datapace/auth/mocks"
	"fmt"
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
	ac := mocks.NewAccessControl()

	return auth.New(repo, hasher, idp, ts, ac)
}

func startGRPCServer(svc auth.Service, port int) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	datapace.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

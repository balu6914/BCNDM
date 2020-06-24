package grpc_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/auth"
	grpcapi "github.com/datapace/datapace/auth/api/grpc"
	"github.com/datapace/datapace/auth/mocks"

	"google.golang.org/grpc"
)

const port = 8081

var svc auth.Service
var k string

func TestMain(m *testing.M) {
	svc, k = newService()
	startGRPCServer(svc, port)
	code := m.Run()
	os.Exit(code)
}

func newServiceWithAdmin() (auth.Service, string, auth.User) {
	hasher := mocks.NewHasher()
	repo := mocks.NewUserRepository(hasher, admin)
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()
	ac := mocks.NewAccessControl()
	svc := auth.New(repo, hasher, idp, ts, ac)
	key, _ := svc.Login(auth.User{
		Email:    admin.Email,
		Password: admin.Password,
	})
	return svc, key, admin
}

func newService() (auth.Service, string) {
	svc, key, _ := newServiceWithAdmin()
	return svc, key
}

func startGRPCServer(svc auth.Service, port int) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	datapace.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

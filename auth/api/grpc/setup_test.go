package grpc_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/datapace/datapace/auth"
	grpcapi "github.com/datapace/datapace/auth/api/grpc"
	"github.com/datapace/datapace/auth/mocks"
	authproto "github.com/datapace/datapace/proto/auth"

	"google.golang.org/grpc"
)

const port = 8000

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
	urepo := mocks.NewUserRepository(hasher, admin, policies, &policiesMu)
	prepo := mocks.NewPolicyRepository(policies, &policiesMu)
	idp := mocks.NewIdentityProvider()
	ts := mocks.NewTransactionsService()
	ac := mocks.NewAccessControl()
	svc := auth.New(urepo, prepo, hasher, idp, ts, ac)
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	server := grpc.NewServer()
	authproto.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

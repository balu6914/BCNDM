package grpc_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/datapace/datapace/streams/groups"
	"github.com/datapace/datapace/streams/sharing"

	streamsproto "github.com/datapace/datapace/proto/streams"
	"github.com/datapace/datapace/streams"
	grpcapi "github.com/datapace/datapace/streams/api/grpc"
	"github.com/datapace/datapace/streams/mocks"

	"google.golang.org/grpc"
)

const (
	port            = 8000
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	server := grpc.NewServer()
	streamsproto.RegisterStreamsServiceServer(server, grpcapi.NewServer(svc))
	go server.Serve(listener)
}

func newService() streams.Service {
	streamRepo := mocks.NewStreamRepository()
	categoryRepo := mocks.NewCategoryRepository()
	ac := mocks.NewAccessControl([]string{})
	ai := mocks.NewAIService()
	terms := mocks.NewTermsService()
	groupsSvc := groups.NewServiceMock()
	sharingSvc := sharing.NewServiceMock()

	return streams.NewService(streamRepo, categoryRepo, ac, ai, terms, groupsSvc, sharingSvc)
}

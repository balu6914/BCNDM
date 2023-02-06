package grpc

import (
	"context"

	"github.com/datapace/datapace/auth"
	authproto "github.com/datapace/datapace/proto/auth"
	commonproto "github.com/datapace/datapace/proto/common"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ authproto.AuthServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	authproto.UnimplementedAuthServiceServer
	identify  kitgrpc.Handler
	email     kitgrpc.Handler
	userById  kitgrpc.Handler
	exists    kitgrpc.Handler
	authorize kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc auth.Service) authproto.AuthServiceServer {
	identify := kitgrpc.NewServer(
		identifyEndpoint(svc),
		decodeIdentifyRequest,
		encodeIdentifyResponse,
	)

	email := kitgrpc.NewServer(
		emailEndpoint(svc),
		decodeIdentifyRequest,
		encodeEmailResponse,
	)

	userById := kitgrpc.NewServer(
		userByEndpoint(svc),
		decodeByIdRequest,
		encodeUserResponse,
	)

	exists := kitgrpc.NewServer(
		existsEndpoint(svc),
		decodeExistsRequest,
		encodeExistsResponse,
	)

	authorize := kitgrpc.NewServer(
		authorizeEndpoint(svc),
		decodeAuthorizeRequest,
		encodeIdentifyResponse,
	)

	return &grpcServer{
		identify:  identify,
		email:     email,
		userById:  userById,
		exists:    exists,
		authorize: authorize,
	}
}

func (s *grpcServer) Identify(ctx context.Context, token *authproto.Token) (*commonproto.ID, error) {
	_, res, err := s.identify.ServeGRPC(ctx, token)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*commonproto.ID), nil
}

func (s *grpcServer) Email(ctx context.Context, id *authproto.Token) (*authproto.UserEmail, error) {
	_, res, err := s.email.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*authproto.UserEmail), nil
}

func (s *grpcServer) UserById(ctx context.Context, id *commonproto.ID) (*authproto.User, error) {
	_, resp, err := s.userById.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return resp.(*authproto.User), nil
}

func (s *grpcServer) Exists(ctx context.Context, id *commonproto.ID) (*empty.Empty, error) {
	_, res, err := s.exists.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*empty.Empty), nil
}

func (s *grpcServer) Authorize(ctx context.Context, req *authproto.AuthRequest) (*commonproto.ID, error) {
	_, res, err := s.authorize.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*commonproto.ID), nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*authproto.Token)
	return identityReq{req.GetValue()}, nil
}

func encodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return &commonproto.ID{Value: res.id}, encodeError(res.err)
}

func encodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(emailRes)
	return &authproto.UserEmail{Email: res.email, ContactEmail: res.contactEmail}, encodeError(res.err)
}

func decodeByIdRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*commonproto.ID)
	return byIdReq{req.GetValue()}, nil
}

func encodeUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	resp := grpcRes.(userResp)
	e := authproto.UserEmail{
		Email:        resp.email,
		ContactEmail: resp.contactEmail,
	}
	u := authproto.User{
		Id:        resp.id,
		Email:     &e,
		FirstName: resp.firstName,
		LastName:  resp.lastName,
		Company:   resp.company,
		Address:   resp.address,
		Country:   resp.country,
		Mobile:    resp.mobile,
		Phone:     resp.phone,
		Role:      resp.role,
	}
	return &u, encodeError(resp.err)
}

func decodeExistsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*commonproto.ID)
	return existsReq{req.GetValue()}, nil
}

func encodeExistsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(existsRes)
	return &empty.Empty{}, encodeError(res.err)
}

func decodeAuthorizeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*authproto.AuthRequest)
	return authReq{
		token:        req.Token,
		action:       req.Action,
		resourceType: req.Type,
		attributes:   req.Attributes,
	}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case auth.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case auth.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from key")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

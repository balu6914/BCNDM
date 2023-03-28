package http

import (
	"context"

	"github.com/datapace/datapace/auth"
	authproto "github.com/datapace/datapace/proto/auth"
	"github.com/datapace/datapace/streams"
	"github.com/go-kit/kit/endpoint"
)

func addStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addStreamReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddStream(req.stream)
		if err != nil {
			return nil, err
		}

		res := addStreamRes{
			ID: id,
		}
		return res, nil
	}
}

func addBulkStreamsEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addBulkStreamsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AddBulkStreams(req.Streams); err != nil {
			switch err.(type) {
			case streams.ErrBulkConflict:
				return conflictBulkStreamsRes{err.(streams.ErrBulkConflict)}, nil
			default:
				return nil, err
			}
		}

		return addBulkStreamsRes{}, nil
	}
}

func updateStreamEndpoint(svc streams.Service, accessSvc streams.AccessControl) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStreamReq)

		//populating stream.ID with value extracted from request URL as there should be no duality here.
		req.stream.ID = req.id

		//TODO: As we have JWT token with sub (userID), we need just to extract it and use instead of additional auth calls.
		token := req.owner
		ar := &authproto.AuthRequest{
			Action: int64(auth.Read),
			Token:  token,
			Type:   streamType,
		}

		// Getting user id through Auth module.
		ownerID, err := authService.Authorize(ar)
		if err != nil {
			return nil, err
		}

		// Assign received ID to req.owner instead of JWT token
		req.owner = ownerID

		// Fetch stream from DB that has owner field populated so that
		// we can use it to authorize Update request.
		referenceStream, err := svc.ViewStream(req.id, req.owner)
		if err != nil {
			return nil, err
		}

		//Assign owner received from database to request stream.
		req.stream.Owner = referenceStream.Owner

		// Now authorizing user to allow update operation on stream from DB.
		// User id will be provided by Auth service, stream from database (referenceStream) will be checked if user (token bearer) can update it.
		ar = &authproto.AuthRequest{
			Action:     int64(auth.Update),
			Token:      token,
			Type:       streamType,
			Attributes: referenceStream.Attributes(),
		}
		ownerID, err = authService.Authorize(ar)
		if err != nil {
			return nil, err
		}

		// If requester is not a stream owner, then check if they're partners of stream owner.
		if req.owner != referenceStream.Owner {
			partners, err := accessSvc.Partners(req.owner)
			if err != nil {
				return nil, err
			}
			for _, p := range partners {
				if p == referenceStream.Owner {
					req.stream.Owner = "" // don't update the stream owner because it's different
					break
				}
			}
		}

		// Need to set owner before the validation because
		// stream.Validate() won't pass otherwise.
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UpdateStream(req.stream); err != nil {
			return nil, err
		}

		return updateStreamRes{}, nil
	}
}

func viewStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewStream(req.id, req.owner)
		if err != nil {
			return nil, err
		}

		res := viewStreamRes{
			Stream: s,
		}
		return res, nil
	}
}

func removeStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(removeStreamReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		token := req.owner
		ar := &authproto.AuthRequest{
			Action: int64(auth.Read),
			Token:  token,
			Type:   streamType,
		}

		ownerID, err := authService.Authorize(ar)
		if err != nil {
			return nil, err
		}

		// Fetch the stream so that attributes can be used for
		// the authorization.
		stream, err := svc.ViewStream(req.id, ownerID)
		if err != nil {
			return nil, err
		}

		ar = &authproto.AuthRequest{
			Action:     int64(auth.Delete),
			Token:      token,
			Type:       streamType,
			Attributes: stream.Attributes(),
		}

		ownerID, err = authService.Authorize(ar)
		if err != nil {
			return nil, err
		}

		if err := svc.RemoveStream(ownerID, req.id); err != nil {
			return nil, err
		}

		return removeStreamRes{}, nil
	}
}

func searchStreamsEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchStreamsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		q := streams.Query{
			Name:       req.Name,
			Owner:      req.Owner,
			StreamType: req.StreamType,
			Coords:     req.Coords,
			Page:       req.Page,
			Limit:      req.Limit,
			MinPrice:   req.MinPrice,
			MaxPrice:   req.MaxPrice,
			Metadata:   req.Metadata,
		}

		page, err := svc.SearchStreams(req.user, q)
		if err != nil {
			return nil, err
		}

		res := searchStreamsRes{page}
		return res, nil
	}
}

func exportStreamsEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(exportStreamsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		results, err := svc.ExportStreams(req.owner)
		if err != nil {
			return nil, err
		}
		res := exportStreamsResp{results}
		return res, nil
	}
}

func addCategoryEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addCategoryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddCategory(req.CategoryName, req.SubCategoryNames)
		if err != nil {
			return nil, err
		}

		return id, nil
	}
}

func listCategoryEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listCategoryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		categoriesList, err := svc.ListCatgeories(req.key)
		if err != nil {
			return nil, err
		}

		return categoriesList, nil
	}
}

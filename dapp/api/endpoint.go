package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"monetasa/dapp"
)

func versionEndpoint() endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		res := versionRes{
			Version: "0.1.0",
		}
		return res, nil
	}
}

func saveStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Save(req.Stream)

		if err != nil {
			return nil, err
		}

		return createStreamRes{}, nil
	}
}

func updateStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Update(req.Id, req.Stream)

		if err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

func oneStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.One(req.Id)

		if err != nil {
			return nil, err
		}

		res := readStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			Price:       s.Price,
		}
		return res, nil
	}
}

func removeStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Remove(req.Id)

		if err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

// func searchStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
// 	return func(_ context.Context, request interface{}) (interface{}, error) {
// 		req := request.(searchStreamReq)
// 		streams, err := svc.Search(req)
// 		res := searchStreamRes{
// 			Streams: streams,
// 		}

// 		fmt.Println("searchStreamEndpoint response")
// 		fmt.Println(res)

// 		return res, err
// 	}
// }

// func purchaseStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
// 	return func(_ context.Context, request interface{}) (interface{}, error) {
// 		req := request.(purchaseStreamReq)
// 		osReq := oneStreamReq{
// 			Name: req.Name,
// 		}
// 		stream, err := svc.One(osReq)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println("Purchasing the stream", stream.Name)
// 		return stream, nil
// 	}
// }

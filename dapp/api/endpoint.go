package api

import (
	"context"
	// "fmt"
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
		req := request.(modifyStreamReq)
		s := dapp.Stream{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			URL:         req.URL,
			Price:       req.Price,
		}
		s, err := svc.Save(s)

		res := modifyStreamRes{
			Status: "success",
		}
		return res, err
	}
}

func updateStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyStreamReq)
		s := dapp.Stream{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			URL:         req.URL,
			Price:       req.Price,
		}
		err := svc.Update(req.Id, s)

		if err != nil {
			return nil, err
		}
		res := modifyStreamRes{
			Status: "success",
		}
		return res, nil
	}
}

func oneStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readStreamReq)
		s, err := svc.One(req.Id)
		res := readStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			Price:       s.Price,
		}
		return res, err
	}
}

func removeStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyStreamReq)
		err := svc.Remove(req.Id)
		if err != nil {
			return nil, err
		}
		res := modifyStreamRes{
			Status: "success",
		}
		return res, nil
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

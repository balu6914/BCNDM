package api

import (
	"context"
	// "fmt"
	"github.com/go-kit/kit/endpoint"

	"monetasa/dapp"
)

func statusEndpoint() endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		res := statusRes{
			Status: "ok",
		}
		return res, nil
	}
}

func saveStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(saveStreamReq)
		s := Stream{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			URL:         req.URL,
			Price:       req.Price,
		}
		s, err := svc.Save(s)

		res := saveStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			URL:         s.URL,
			Price:       s.Price,
		}
		return res, err
	}
}

func updateStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(saveStreamReq)
		s := Stream{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			URL:         req.URL,
			Price:       req.Price,
		}
		err := svc.Update(req.Name, s)

		if err != nil {
			return nil, err
		}
		res := saveStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			URL:         s.URL,
			Price:       s.Price,
		}
		return res, nil
	}
}

func oneStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(oneStreamReq)
		s, err := svc.One(req.Name)
		res := oneStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			Price:       s.Price,
		}
		return res, err
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

func removeStreamEndpoint(svc StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(removeStreamReq)
		err := svc.Remove(req.Name)
		if err != nil {
			return nil, err
		}
		res := removeStreamRes{
			Status: "ok",
		}
		return res, nil
	}
}

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

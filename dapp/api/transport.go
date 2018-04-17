package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

	"monetasa/dapp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(sr StreamRepository) http.Handler {
	r := bone.New()

	r.Get("/status", kithttp.NewServer(
		statusEndpoint(),
		decodeStatusRequest,
		encodeStatusResponse,
	))

	r.Post("/streams", kithttp.NewServer(
		saveStreamEndpoint(sr),
		decodeSaveStreamRequest,
		encodeSaveStreamResponse,
	))

	r.Put("/streams/:name", kithttp.NewServer(
		updateStreamEndpoint(sr),
		decodeUpdateStreamRequest,
		encodeUpdateStreamResponse,
	))

	// r.Get("/streams/search", kithttp.NewServer(
	// 	searchStreamEndpoint(sr),
	// 	decodeSearchStreamRequest,
	// 	encodeSearchStreamResponse,
	// ))

	r.Get("/streams/:name", kithttp.NewServer(
		oneStreamEndpoint(sr),
		decodeOneStreamRequest,
		encodeOneStreamResponse,
	))

	r.Delete("/streams", kithttp.NewServer(
		removeStreamEndpoint(sr),
		decodeRemoveStreamRequest,
		encodeRemoveStreamResponse,
	))

	// r.Post("/streams/purch", kithttp.NewServer(
	// 	purchaseStreamEndpoint(sr),
	// 	decodePurchaseStreamRequest,
	// 	encodePurchaseStreamResponse,
	// ))

	return r
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeStatusResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeSaveStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req saveStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeSaveStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("encodeSaveStreamResponse response: ", response)
	return json.NewEncoder(w).Encode(response)
}

func decodeUpdateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req saveStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Name = bone.GetValue(r, "name")
	return req, nil
}

func encodeUpdateStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("encodeSaveStreamResponse response: ", response)
	return json.NewEncoder(w).Encode(response)
}

func decodeOneStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req oneStreamReq
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	return nil, err
	// }
	req.Name = bone.GetValue(r, "name")
	return req, nil
}

func encodeOneStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeSearchStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req searchStreamReq
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	return nil, err
	// }
	q := r.URL.Query()

	req.Type = q["type"][0]
	req.x0, _ = strconv.Atoi(q["x0"][0])
	req.y0, _ = strconv.Atoi(q["y0"][0])
	req.x1, _ = strconv.Atoi(q["x1"][0])
	req.y1, _ = strconv.Atoi(q["y1"][0])
	req.x2, _ = strconv.Atoi(q["x2"][0])
	req.y2, _ = strconv.Atoi(q["y2"][0])
	req.x3, _ = strconv.Atoi(q["x3"][0])
	req.y3, _ = strconv.Atoi(q["y3"][0])

	fmt.Println("decodeSearchStreamRequest")
	fmt.Println(req)
	return req, nil
}

func encodeSearchStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeRemoveStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req removeStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRemoveStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodePurchaseStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req purchaseStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodePurchaseStreamResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

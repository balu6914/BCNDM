package grpc

type updateReceiversResponse struct {
	err error
}

type getSharingsResponse struct {
	sharings []sharingPayload
}

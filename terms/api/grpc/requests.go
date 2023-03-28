package grpc

type termsReq struct {
	id       string
	streamId string
	termsUrl string
	metadata map[string]string
}

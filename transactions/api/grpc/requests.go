package grpc

type createUserReq struct {
	id     string
	secret string
}

func (req createUserReq) validate() error {
	if req.id == "" || req.secret == "" {
		return errMalformedEntity
	}

	return nil
}

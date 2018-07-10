package grpc

type createUserReq struct {
	id string
}

func (req createUserReq) validate() error {
	if req.id == "" {
		return errMalformedEntity
	}

	return nil
}

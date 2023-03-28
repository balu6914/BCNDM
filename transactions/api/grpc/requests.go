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

type transferReq struct {
	streamID string
	from     string
	to       string
	value    uint64
	DateTime string
}

func (req transferReq) validate() error {
	if req.streamID == "" || req.from == "" || req.to == "" {
		return errMalformedEntity
	}

	return nil
}

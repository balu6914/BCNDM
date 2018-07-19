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
	from  string
	to    string
	value uint64
}

func (req transferReq) validate() error {
	if req.from == "" || req.to == "" {
		return errMalformedEntity
	}

	return nil
}

package http

type balanceReq struct {
	userID string
	chanID string
}

func (req balanceReq) validate() error {
	if req.userID == "" || req.chanID == "" {
		return errMalformedEntity
	}

	return nil
}

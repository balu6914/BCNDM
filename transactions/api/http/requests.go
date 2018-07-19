package http

type balanceReq struct {
	userID string
}

func (req balanceReq) validate() error {
	if req.userID == "" {
		return errMalformedEntity
	}

	return nil
}

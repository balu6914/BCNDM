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

type buyReq struct {
	userID string
	Amount uint64 `json:"amount"`
}

func (req buyReq) validate() error {
	if req.Amount == 0 {
		return errMalformedEntity
	}

	return nil
}

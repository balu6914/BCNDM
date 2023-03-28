package fabric

type accessReq struct {
	Receiver string `json:"receiver"`
}

type approveReq struct {
	Requester string `json:"requester"`
}

type revokeReq struct {
	Requester string `json:"requester"`
}

type grantReq struct {
	Destination string `json:"destination"`
}

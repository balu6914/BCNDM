package http

const contentType = "application/json"

type createTokenRes struct {
	Token string `json:"token"`
}

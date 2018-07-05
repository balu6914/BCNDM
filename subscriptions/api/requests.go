package api

import "monetasa/subscriptions"

type apiReq interface {
	validate() error
}

type subscriptionReq struct {
	UserID       string
	Subscription subscriptions.Subscription
}

func (req subscriptionReq) validate() error {
	return req.Subscription.Validate()
}

type listSubsReq struct {
	UserID string
}

func (req listSubsReq) validate() error {
	if req.UserID == "" {
		return subscriptions.ErrMalformedEntity
	}

	return nil
}

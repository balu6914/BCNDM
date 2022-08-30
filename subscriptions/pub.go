package subscriptions

import "github.com/datapace/datapace/subscriptions/pub"

type (
	pubMiddleware struct {
		svc    Service
		pubSvc pub.Service
	}
)

func NewPubMiddleware(svc Service, pubSvc pub.Service) Service {
	return pubMiddleware{
		svc:    svc,
		pubSvc: pubSvc,
	}
}

func (pm pubMiddleware) AddSubscription(userId string, authToken string, inSub Subscription) (outSub Subscription, err error) {
	defer func() {
		if err == nil {
			evt := pub.SubscriptionCreateEvent{
				SubscriptionId: outSub.ID.Hex(),
			}
			to := outSub.StreamOwner
			_, _ = pm.pubSvc.Publish(evt, to)
		}
	}()
	return pm.svc.AddSubscription(userId, authToken, inSub)
}

func (pm pubMiddleware) SearchSubscriptions(query Query) (Page, error) {
	return pm.svc.SearchSubscriptions(query)
}

func (pm pubMiddleware) ViewSubscription(userId string, subscriptionId string) (Subscription, error) {
	return pm.svc.ViewSubscription(userId, subscriptionId)
}

func (pm pubMiddleware) ViewSubByUserAndStream(userId string, streamId string) (Subscription, error) {
	return pm.svc.ViewSubByUserAndStream(userId, streamId)
}

func (pm pubMiddleware) Report(query Query, userId string) ([]byte, error) {
	return pm.svc.Report(query, userId)
}

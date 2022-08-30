package pub

import (
	"encoding/json"
	"fmt"
	"github.com/datapace/events/pubsub"
)

type (

	// SubscriptionCreateEvent is the event to be issued on a new subscription
	SubscriptionCreateEvent struct {

		// SubscriptionId is the new subscription id
		SubscriptionId string `json:"subscriptionId"`
	}

	// SubjectFormat contains the subscriptions releated subjects configuration.
	SubjectFormat struct {

		// SubscriptionCreate is the subject format for the SubscriptionCreateEvent.
		SubscriptionCreate string
	}

	// Service is the subscription events publish service.
	Service interface {

		// Publish a SubscriptionCreateEvent to the specified user.
		// Returns sent message id, error otherwise.
		Publish(evt SubscriptionCreateEvent, toUserId string) (uint64, error)
	}

	service struct {
		pubSubSvc pubsub.Service
		subjFmt   SubjectFormat
	}
)

func NewService(pubSubSvc pubsub.Service, subjFmt SubjectFormat) Service {
	return service{
		pubSubSvc: pubSubSvc,
		subjFmt:   subjFmt,
	}
}

func (svc service) Publish(evt SubscriptionCreateEvent, toUserId string) (uint64, error) {
	if svc.pubSubSvc == nil {
		return 0, nil
	}
	subject := fmt.Sprintf(svc.subjFmt.SubscriptionCreate, toUserId)
	return svc.publishJson(subject, evt)
}

func (svc service) publishJson(subject string, evt interface{}) (uint64, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return 0, err
	}
	return svc.pubSubSvc.Publish(subject, map[string][]string{}, data)
}

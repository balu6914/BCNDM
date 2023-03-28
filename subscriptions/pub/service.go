package pub

import (
	"encoding/json"
	"fmt"
	"github.com/datapace/events/pubsub"
)

type (

	// SubjectFormat contains the subscriptions releated subjects configuration.
	SubjectFormat struct {

		// SubscriptionCreate is the subject format for the SubscriptionCreateEvent.
		SubscriptionCreate string
	}

	// Service is the subscription events publish service.
	Service interface {

		// PublishSubscriptionCreated event to the specified user.
		// Returns sent message id, error otherwise.
		PublishSubscriptionCreated(evt interface{}, toUserId string) (uint64, error)
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

func (svc service) PublishSubscriptionCreated(evt interface{}, toUserId string) (uint64, error) {
	if svc.pubSubSvc == nil {
		return 0, nil
	}
	subject := fmt.Sprintf(svc.subjFmt.SubscriptionCreate, toUserId)
	return svc.publishJson(evt, subject)
}

func (svc service) publishJson(evt interface{}, subject string) (uint64, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return 0, err
	}
	return svc.pubSubSvc.Publish(subject, map[string][]string{}, data)
}

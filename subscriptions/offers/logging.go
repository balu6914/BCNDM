package offers

import (
	"context"
	"fmt"

	log "github.com/datapace/datapace/logger"
)

type loggingMiddleware struct {
	svc    Service
	logger log.Logger
}

func NewLoggingMiddleware(svc Service, logger log.Logger) Service {
	return loggingMiddleware{
		svc:    svc,
		logger: logger,
	}
}

func (lm loggingMiddleware) GetPrice(ctx context.Context, k OfferKey) (p OfferPrice, err error) {
	msg := fmt.Sprintf("GetPrice(%v)", k)
	defer func() {
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s: %s", msg, err.Error()))
		} else {
			lm.logger.Info(msg)
		}
	}()
	return lm.svc.GetPrice(ctx, k)
}

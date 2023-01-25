package accessv2

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

func (lm loggingMiddleware) Get(ctx context.Context, k Key) (a Access, err error) {
	msg := fmt.Sprintf("Get(%v)", k)
	defer func() {
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s: %s", msg, err.Error()))
		} else {
			lm.logger.Info(msg)
		}
	}()
	return lm.svc.Get(ctx, k)
}

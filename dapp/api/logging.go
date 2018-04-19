package api

import (
	"fmt"
	"time"

	"monetasa/dapp"
	log "monetasa/logger"
)

var _ dapp.StreamRepository = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    dapp.StreamRepository
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc dapp.StreamRepository, logger log.Logger) dapp.StreamRepository {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Save(stream dapp.Stream) (str dapp.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method save for stream %s took %s to complete", stream.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Save(stream)
}

func (lm *loggingMiddleware) Update(name string, stream dapp.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method Update for stream %s took %s to complete", stream.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Update(name, stream)
}

func (lm *loggingMiddleware) One(name string) (stream dapp.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method One for stream %s, took %s to complete", name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.One(name)
}

func (lm *loggingMiddleware) Remove(name string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method Remove for stream %s, took %s to complete", name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Remove(name)
}

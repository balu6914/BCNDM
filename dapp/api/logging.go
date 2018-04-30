package api

import (
	"fmt"
	"time"

	"monetasa/dapp"
	log "monetasa/logger"
)

var _ dapp.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    dapp.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc dapp.Service, logger log.Logger) dapp.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddStream(key string, stream dapp.Stream) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method save for stream %s took %s to complete", stream.Name, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddStream(key, stream)
}

func (lm *loggingMiddleware) UpdateStream(key string, id string, stream dapp.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method Update for stream %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.UpdateStream(key, id, stream)
}

func (lm *loggingMiddleware) ViewStream(key string, id string) (stream dapp.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method One for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewStream(key, id)
}

func (lm *loggingMiddleware) RemoveStream(key string, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method Remove for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveStream(key, id)
}

func (lm *loggingMiddleware) SearchStreams(coords [][]float64) (streams []dapp.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method Search for took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.SearchStreams(coords)
}

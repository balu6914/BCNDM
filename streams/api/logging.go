package api

import (
	"fmt"
	"time"

	log "monetasa/logger"
	"monetasa/streams"
)

var _ streams.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    streams.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc streams.Service, logger log.Logger) streams.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddStream(key string, stream streams.Stream) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_stream for stream %s took %s to complete", stream.ID.Hex(), time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddStream(key, stream)
}

func (lm *loggingMiddleware) AddBulkStream(key string, streams []streams.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_bulk_stream for streams %v took %s to complete", streams, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddBulkStream(key, streams)
}

func (lm *loggingMiddleware) ViewStream(id string) (stream streams.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_stream for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewStream(id)
}

func (lm *loggingMiddleware) UpdateStream(key string, stream streams.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_stream for stream %s took %s to complete", stream.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.UpdateStream(key, stream)
}

func (lm *loggingMiddleware) SearchStreams(coords [][]float64) (streams []streams.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method search_streams for took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.SearchStreams(coords)
}

func (lm *loggingMiddleware) RemoveStream(key string, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_stream for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveStream(key, id)
}

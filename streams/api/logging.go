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

func (lm *loggingMiddleware) AddStream(stream streams.Stream) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_stream for stream %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddStream(stream)
}

func (lm *loggingMiddleware) AddBulkStreams(streams []streams.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_bulk_stream for streams of size %d took %s to complete", len(streams), time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddBulkStreams(streams)
}

func (lm *loggingMiddleware) ViewFullStream(id string) (stream streams.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_full_stream for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewFullStream(id)
}

func (lm *loggingMiddleware) ViewStream(id, owner string) (stream streams.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_stream for stream %s, took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewStream(id, owner)
}

func (lm *loggingMiddleware) UpdateStream(stream streams.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_stream for stream %s took %s to complete", stream.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.UpdateStream(stream)
}

func (lm *loggingMiddleware) SearchStreams(owner string, query streams.Query) (page streams.Page, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method search_streams for took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.SearchStreams(owner, query)
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

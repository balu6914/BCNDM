package api

import (
	"datapace/executions"
	log "datapace/logger"
	"fmt"
	"time"
)

var _ executions.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    executions.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc executions.Service, logger log.Logger) executions.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm loggingMiddleware) Start(owner, algo, data string, mode executions.JobMode) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method start for owner %s, algo %s, data %s and mode %s took %s to complete",
			owner, algo, data, mode, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Start(owner, algo, data, mode)
}

func (lm loggingMiddleware) Finish(id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method finish for execution %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Finish(id)
}

func (lm loggingMiddleware) Execution(owner, id string) (exec executions.Execution, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method execution for owner %s and execution %s took %s to complete", owner, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Execution(owner, id)
}

func (lm loggingMiddleware) List(owner string) (execs []executions.Execution, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list for owner %s took %s to complete", owner, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.List(owner)
}

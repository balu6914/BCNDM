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

func (lm loggingMiddleware) Start(exec executions.Execution) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method start for owner %s, algo %s, data %s and mode %s with id %s took %s to complete",
			exec.Owner, exec.Algo, exec.Data, exec.Mode, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Start(exec)
}

func (lm loggingMiddleware) Result(owner, id string) (result map[string]interface{}, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method result for owner %s and execution %s took %s to complete", owner, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Result(owner, id)
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

func (lm loggingMiddleware) CreateAlgorithm(algo executions.Algorithm) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_algorithm for algorithm %s took %s to complete", algo.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateAlgorithm(algo)
}

func (lm loggingMiddleware) CreateDataset(data executions.Dataset) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_dataset for dataset %s took %s to complete", data.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateDataset(data)
}

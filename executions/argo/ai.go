package argo

import (
	"encoding/json"
	"fmt"
	"github.com/datapace/datapace/executions"
	log "github.com/datapace/datapace/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var _ executions.AIService = (*aiService)(nil)

type aiService struct {
	url      string
	eventCh  chan executions.Event
	interval time.Duration
	logger   log.Logger
}

// New returns new AI service instance that is implemented using argo.
func New(url string, interval time.Duration, logger log.Logger) executions.AIService {
	return aiService{
		url:      url,
		eventCh:  make(chan executions.Event),
		interval: interval,
		logger:   logger,
	}
}

func (as aiService) CreateAlgorithm(algo executions.Algorithm) (executions.Algorithm, error) {
	config, ok := algo.Metadata["pipeline"]
	if !ok {
		return executions.Algorithm{}, executions.ErrMalformedData
	}
	buf := strings.NewReader(config)
	path := fmt.Sprintf("%s/api/v1/workflow-templates/default", as.url)
	res, err := http.Post(path, "", buf)
	if err != nil {
		return executions.Algorithm{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		emsg, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(emsg))
		return executions.Algorithm{}, executions.ErrCreateAlgoFailed
	}
	var car createAlgoResponse
	if err := json.NewDecoder(res.Body).Decode(&car); err != nil {
		fmt.Println(err.Error())
		return executions.Algorithm{}, executions.ErrCreateAlgoFailed
	}

	algo.ExternalID = car.Metadata.Name
	return algo, nil
}

func (as aiService) CreateDataset(data executions.Dataset) (executions.Dataset, error) {
	// For now do nothing.
	return data, nil
}

func (as aiService) Start(exec executions.Execution, algo executions.Algorithm, data executions.Dataset) (executions.Execution, error) {
	path := fmt.Sprintf("%s/api/v1/workflows/default", as.url)
	dpath, ok := data.Metadata["path"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	post := fmt.Sprintf("{\"workflow\":{\"metadata\":{\"generateName\":\"datapace-\",\"namespace\":\"default\"},\"spec\":{\"entrypoint\":\"datapace\",\"arguments\":{\"parameters\":[{\"name\":\"datapace_url\",\"value\":\"%s\"}]},\"workflowTemplateRef\":{\"name\":\"%s\"}}}}", dpath, algo.ExternalID)
	buf := strings.NewReader(post)
	res, err := http.Post(path, "", buf)
	if err != nil {
		return executions.Execution{}, executions.ErrExecutionFailed
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		emsg, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(emsg))
		return executions.Execution{}, executions.ErrExecutionFailed
	}
	var cer execRunResponse
	if err := json.NewDecoder(res.Body).Decode(&cer); err != nil {
		fmt.Println(err.Error())
		return executions.Execution{}, executions.ErrCreateAlgoFailed
	}

	exec.ExternalID = cer.Metadata.Name
	// Start fetching status of the run
	go as.pullStatus(exec.ExternalID)
	return exec, nil
}

func (as aiService) Result(exec executions.Execution) (map[string]interface{}, error) {
	return nil, nil
}

func (as aiService) Events() (chan executions.Event, error) {
	return as.eventCh, nil
}

func (as aiService) pullStatus(externalID string) {
	ticker := time.NewTicker(as.interval)
	for range ticker.C {
		st, err := as.runStatus(externalID)
		if err != nil {
			// log error
			fmt.Println(err)
			continue
		}
		if st == executions.InProgress {
			continue
		}

		as.eventCh <- executions.Event{
			ExternalID: externalID,
			Status:     st,
		}

		ticker.Stop()
	}
}

func (as aiService) runStatus(externalID string) (executions.State, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/workflows/default/%s", as.url, externalID))
	if err != nil {
		return executions.State(""), err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return executions.State(""), executions.ErrNotFound
	}

	var cer execRunResponse
	if err := json.NewDecoder(res.Body).Decode(&cer); err != nil {
		return executions.State(""), err
	}

	switch cer.Status.Phase {
	case "Failed":
		return executions.Failed, nil
	case "Succeeded":
		return executions.Succeeded, nil
	default:
		return executions.InProgress, nil
	}
}

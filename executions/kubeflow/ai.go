package kubeflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/datapace/datapace/executions"
	log "github.com/datapace/datapace/logger"
	"github.com/datapace/datapace/streams"
)

const (
	defContentType  = "application/json"
	defStorageState = "STORAGESTATE_AVAILABLE"
	datapceURL      = "datapace_url"
)

var _ executions.AIService = (*aiService)(nil)

type aiService struct {
	url      string
	eventCh  chan executions.Event
	interval time.Duration
	logger   log.Logger
}

// New returns new AI service instance that is implemented using kubeflow.
func New(url string, interval time.Duration, logger log.Logger) executions.AIService {
	return aiService{
		url:      url,
		eventCh:  make(chan executions.Event),
		interval: interval,
		logger:   logger,
	}
}

func (as aiService) CreateAlgorithm(algo executions.Algorithm) (executions.Algorithm, error) {
	// Post pipeline yaml file to /apis/v1beta1/pipelines/upload
	buf := &bytes.Buffer{}
	bw := multipart.NewWriter(buf)

	config, ok := algo.Metadata["pipeline"]
	if !ok {
		return executions.Algorithm{}, executions.ErrMalformedData
	}

	filename := fmt.Sprintf("%s.yaml", algo.Name)

	fw, err := bw.CreateFormFile("uploadfile", filename)
	if err != nil {
		return executions.Algorithm{}, err
	}

	if _, err := io.Copy(fw, strings.NewReader(config)); err != nil {
		return executions.Algorithm{}, err
	}

	ct := bw.FormDataContentType()
	bw.Close()
	path := fmt.Sprintf("%s/pipeline/apis/v1beta1/pipelines/upload?name=%s", as.url, algo.Name)
	res, err := http.Post(path, ct, buf)
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

	algo.ExternalID = car.ID
	return algo, nil
}

func (as aiService) CreateDataset(data executions.Dataset) (executions.Dataset, error) {
	// For now do nothing.
	return data, nil
}

func (as aiService) Start(exec executions.Execution, algo executions.Algorithm, data executions.Dataset) (executions.Execution, error) {
	// Start run of pipeline over dataset passed as parameter POST /apis/v1beta1/runs
	ss := ""
	val, ok := exec.Metadata["storage_state"]
	if !ok {
		ss = defStorageState
	}
	ss, ok = val.(string)
	if !ok {
		ss = defStorageState
	}

	dpath, ok := data.Metadata["path"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	params := []paramRequest{
		{
			Name:  datapceURL,
			Value: dpath,
		},
	}
	ps, ok := exec.Metadata["params"]
	if ok {
		mps, ok := ps.(map[string]string)
		if !ok {
			return executions.Execution{}, executions.ErrMalformedData
		}
		for k, v := range mps {
			params = append(params, paramRequest{
				Name:  k,
				Value: v,
			})
		}
	}

	name := exec.Name
	if name == "" {
		name = fmt.Sprintf("%s.%s", algo.ID, data.ID)
	}

	req := createJobRequest{
		Name:         name,
		StorageState: ss,
		PipelineSpec: pipelineSpecRequest{
			ID:         algo.ExternalID,
			Parameters: params,
		},
	}

	path := fmt.Sprintf("%s/pipeline/apis/v1beta1/runs", as.url)
	res, err := sendRequest(http.MethodPost, path, req)
	if err != nil {
		return executions.Execution{}, err
	}

	var cer execRunResponse
	if err := json.Unmarshal(res, &cer); err != nil {
		return executions.Execution{}, executions.ErrExecutionFailed
	}

	exec.ExternalID = cer.Run.ID
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
	res, err := http.Get(fmt.Sprintf("%s/pipeline/apis/v1beta1/runs/%s", as.url, externalID))
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

	switch cer.Run.Status {
	case "Failed":
		return executions.Failed, nil
	case "Succeeded":
		return executions.Succeeded, nil
	default:
		return executions.InProgress, nil
	}
}

func sendRequest(method, url string, r request) ([]byte, error) {
	body, err := r.body()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for key, val := range r.headers() {
		req.Header.Add(key, val)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusUnauthorized:
			return nil, streams.ErrUnauthorizedAccess
		case http.StatusConflict:
			return nil, streams.ErrConflict
		default:
			fmt.Println(res.StatusCode)
			fmt.Println(string(data))
			return nil, errors.New("Unknown error")
		}
	}

	return data, err
}

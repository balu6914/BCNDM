package wwh

import (
	"bytes"
	"datapace/executions"
	"datapace/streams"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	jsonCT       = "application/json"
	relationship = "Relationship"
	attribute    = "Attribute"
)

var _ executions.AIService = (*aiService)(nil)

type aiService struct {
	catalogURL string
	daemonURL  string
	token      string
	username   string
	password   string
}

// NewAIService returns WWH implementation of AI service.
func NewAIService(catalogURL, daemonURL, token, username, password string) executions.AIService {
	return aiService{
		catalogURL: catalogURL,
		daemonURL:  daemonURL,
		token:      token,
		username:   username,
		password:   password,
	}
}

func (as aiService) CreateAlgorithm(algo executions.Algorithm) error {
	lid := fmt.Sprintf("local-%s", algo.ID)
	if err := as.createExecution(lid); err != nil {
		return err
	}

	gid := fmt.Sprintf("global-%s", algo.ID)
	if err := as.createExecution(gid); err != nil {
		return err
	}

	cid := fmt.Sprintf("%s-comp", algo.ID)
	if err := as.createComputation(cid, lid, gid); err != nil {
		return err
	}

	clid := fmt.Sprintf("%s-comp-list", algo.ID)
	return as.createComputationList(clid, cid, "") // leave empty algo parameters for now
}

func (as aiService) CreateDataset(ds executions.Dataset) error {
	hrid := fmt.Sprintf("%s-hr", ds.ID)
	if err := as.createHardResource(hrid, ds.Path); err != nil {
		return err
	}

	mrid := fmt.Sprintf("%s-mr", ds.ID)
	return as.createMetaResource(mrid, hrid)
}

func (as aiService) Start(exec executions.Execution, algo executions.Algorithm, data executions.Dataset) (string, error) {
	url := fmt.Sprintf("%s/daemon/exec/execute", as.daemonURL)
	er := executeReq{
		token:                    as.token,
		MetaResources:            []string{fmt.Sprintf("%s-mr", data.ID)},
		ChainComputation:         algo.Name,
		AdditionalLocalJobArgs:   exec.AdditionalLocalJobArgs,
		LogicJarPath:             algo.Path,
		RequestID:                exec.ID,
		Type:                     exec.Type,
		ResultPath:               fmt.Sprintf("/%s", exec.ID),
		GlobalTimeout:            strconv.FormatUint(exec.GlobalTimeout, 10),
		LocalTimeout:             strconv.FormatUint(exec.LocalTimeout, 10),
		AdditionalPreprocessArgs: exec.AdditionalPreprocessArgs,
		JobMode:                  exec.Mode,
		AdditionalGlobalJobArgs:  exec.AdditionalGlobalJobArgs,
		AdditionalFiles:          exec.AdditionalFiles,
		TrainedModelToken:        algo.ModelToken,
		TrainedModelName:         algo.ModelName,
	}

	res, err := sendRequest(http.MethodPost, url, er)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (as aiService) IsDone(token string) (executions.State, error) {
	url := fmt.Sprintf("%s/daemon/exec/isDone/%s", as.daemonURL, token)
	req := tokenReq{
		token: as.token,
	}
	res, err := sendRequest(http.MethodGet, url, req)
	if err != nil {
		return executions.Executing, err
	}

	done, err := strconv.ParseBool(string(res))
	if err != nil {
		return executions.Executing, err
	}

	if !done {
		return executions.Executing, nil
	}

	return executions.Done, nil
}

func (as aiService) Result(token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/daemon/exec/getResults/%s", as.daemonURL, token)
	req := tokenReq{
		token: as.token,
	}

	res, err := sendRequest(http.MethodGet, url, req)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return map[string]interface{}{
			"error": "execution failed",
		}, nil
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, executions.ErrExecutionFailed
	}

	return result[0], nil
}

func (as aiService) createExecution(id string) error {
	url := fmt.Sprintf("%s/catalog/createInstanceAndProps/PythonExecution/%s", as.catalogURL, id)
	cr := createReq{
		username: as.username,
		password: as.password,
		properties: []propertyReq{
			{
				Name:  "type",
				Type:  attribute,
				Value: "PYTHON-ML",
			},
			{
				Name:  "fullClassPath",
				Type:  attribute,
				Value: "",
			},
		},
	}

	_, err := sendRequest(http.MethodPost, url, cr)
	return err
}

func (as aiService) createComputation(cid, lid, gid string) error {
	url := fmt.Sprintf("%s/catalog/createInstanceAndProps/TwoExecutionsComputation/%s", as.catalogURL, cid)
	cr := createReq{
		username: as.username,
		password: as.password,
		properties: []propertyReq{
			{
				Name:  "execution",
				Type:  relationship,
				Value: []string{lid},
			},
			{
				Name:  "secondExecution",
				Type:  relationship,
				Value: []string{gid},
			},
		},
	}

	_, err := sendRequest(http.MethodPost, url, cr)
	return err
}

func (as aiService) createComputationList(clid, cid, params string) error {
	url := fmt.Sprintf("%s/catalog/createInstanceAndProps/ChainComputationList/%s", as.catalogURL, clid)
	cr := createReq{
		username: as.username,
		password: as.password,
		properties: []propertyReq{
			{
				Name:  "first",
				Type:  relationship,
				Value: []string{cid},
			},
			{
				Name:  "local-parameters-descriptors",
				Type:  attribute,
				Value: params,
			},
			{
				Name:  "member_of",
				Type:  relationship,
				Value: []string{"__self__"},
			},
		},
	}

	_, err := sendRequest(http.MethodPost, url, cr)
	return err
}

func (as aiService) createHardResource(hrid, path string) error {
	url := fmt.Sprintf("%s/catalog/createInstanceAndProps/HardResource/%s", as.catalogURL, hrid)
	cr := createReq{
		username: as.username,
		password: as.password,
		properties: []propertyReq{
			{
				Name:  "path",
				Type:  attribute,
				Value: path,
			},
			{
				Name:  "catalogRef",
				Type:  attribute,
				Value: "__self__",
			},
		},
	}

	_, err := sendRequest(http.MethodPost, url, cr)
	return err
}

func (as aiService) createMetaResource(mrid, hrid string) error {
	url := fmt.Sprintf("%s/catalog/createInstanceAndProps/MetaResource/%s", as.catalogURL, mrid)
	cr := createReq{
		username: as.username,
		password: as.password,
		properties: []propertyReq{
			{
				Name:  "references",
				Type:  relationship,
				Value: []string{hrid},
			},
			{
				Name:  "member_of",
				Type:  relationship,
				Value: []string{"__self__"},
			},
		},
	}

	_, err := sendRequest(http.MethodPost, url, cr)
	return err
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
			return nil, errors.New("Unknown error")
		}
	}

	return data, err
}

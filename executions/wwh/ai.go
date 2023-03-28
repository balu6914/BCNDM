package wwh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/datapace/datapace/executions"
	"github.com/datapace/datapace/streams"

	"github.com/gorilla/websocket"
)

const (
	jsonCT          = "application/json"
	relationship    = "Relationship"
	attribute       = "Attribute"
	stompTerminator = "\u0000"
	stompID         = "0"
	stompTopic      = "/executionDetails"
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

func (as aiService) CreateAlgorithm(algo executions.Algorithm) (executions.Algorithm, error) {
	lid := fmt.Sprintf("local-%s", algo.ID)
	if err := as.createExecution(lid); err != nil {
		return algo, err
	}

	gid := fmt.Sprintf("global-%s", algo.ID)
	if err := as.createExecution(gid); err != nil {
		return algo, err
	}

	cid := fmt.Sprintf("%s-comp", algo.ID)
	if err := as.createComputation(cid, lid, gid); err != nil {
		return algo, err
	}

	clid := fmt.Sprintf("%s-comp-list", algo.ID)
	return algo, as.createComputationList(clid, cid, "") // leave empty algo parameters for now
}

func (as aiService) CreateDataset(ds executions.Dataset) (executions.Dataset, error) {
	hrid := fmt.Sprintf("%s-hr", ds.ID)
	path, ok := ds.Metadata["path"]
	if !ok {
		return ds, executions.ErrMalformedData
	}

	if err := as.createHardResource(hrid, path); err != nil {
		return ds, err
	}

	mrid := fmt.Sprintf("%s-mr", ds.ID)
	return ds, as.createMetaResource(mrid, hrid)
}

func (as aiService) Start(exec executions.Execution, algo executions.Algorithm, data executions.Dataset) (executions.Execution, error) {
	url := fmt.Sprintf("http://%s/daemon/exec/execute", as.daemonURL)

	algoPath, ok := algo.Metadata["path"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	execType, ok := exec.Metadata["type"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	et, ok := execType.(string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	additionalLocalJobArgs, ok := exec.Metadata["additionalLocalJobArguments"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	aljargs, ok := additionalLocalJobArgs.([]string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	globalTime, ok := exec.Metadata["globalTimeout"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	gtime, ok := globalTime.(uint64)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	localTime, ok := exec.Metadata["localTimeout"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	ltime, ok := localTime.(uint64)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	additionalPreprocessArgs, ok := exec.Metadata["additionalPreprocessArguments"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	apa, ok := additionalPreprocessArgs.([]string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	jobMode, ok := exec.Metadata["jobMode"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	jmode, ok := jobMode.(string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	additionalGlobalJobArgs, ok := exec.Metadata["additionalGlobalJobArguments"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	agja, ok := additionalGlobalJobArgs.([]string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	additionalFiles, ok := exec.Metadata["additionalFiles"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}
	af, ok := additionalFiles.([]string)
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	algoModelToken, ok := algo.Metadata["trainedModelToken"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	algoModelName, ok := algo.Metadata["trainedModelName"]
	if !ok {
		return executions.Execution{}, executions.ErrMalformedData
	}

	er := executeReq{
		token:                    as.token,
		MetaResources:            []string{fmt.Sprintf("%s-mr", data.ID)},
		ChainComputation:         algo.Name,
		AdditionalLocalJobArgs:   aljargs,
		LogicJarPath:             algoPath,
		RequestID:                exec.ID,
		Type:                     et,
		ResultPath:               fmt.Sprintf("/%s", exec.ID),
		GlobalTimeout:            strconv.FormatUint(gtime, 10),
		LocalTimeout:             strconv.FormatUint(ltime, 10),
		AdditionalPreprocessArgs: apa,
		JobMode:                  jmode,
		AdditionalGlobalJobArgs:  agja,
		AdditionalFiles:          af,
		TrainedModelToken:        algoModelToken,
		TrainedModelName:         algoModelName,
	}

	res, err := sendRequest(http.MethodPost, url, er)
	if err != nil {
		return executions.Execution{}, err
	}

	exec.ExternalID = string(res)
	return exec, nil
}

func (as aiService) Result(exec executions.Execution) (map[string]interface{}, error) {
	token, ok := exec.Metadata["token"]
	if !ok {
		return nil, executions.ErrMalformedData
	}
	tkn, ok := token.(string)
	if !ok {
		return nil, executions.ErrMalformedData
	}

	url := fmt.Sprintf("http://%s/daemon/exec/getResults/%s", as.daemonURL, tkn)
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

func (as aiService) Events() (chan executions.Event, error) {
	url := fmt.Sprintf("ws://%s/daemon/socket", as.daemonURL)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	connectPayload := fmt.Sprintf("CONNECT\n\n%s", stompTerminator)
	if err := c.WriteMessage(websocket.TextMessage, []byte(connectPayload)); err != nil {
		return nil, err
	}

	subPayload := fmt.Sprintf("SUBSCRIBE\nid:%s\ndestination:%s\n\n%s", stompID, stompTopic, stompTerminator)
	if err := c.WriteMessage(websocket.TextMessage, []byte(subPayload)); err != nil {
		return nil, err
	}

	ch := make(chan executions.Event)
	go readEvents(c, ch)

	return ch, nil
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

func readEvents(c *websocket.Conn, ch chan executions.Event) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			continue
		}

		lines := strings.Split(string(msg), "\n")
		if len(lines) < 1 || lines[0] != "MESSAGE" {
			continue
		}

		content := strings.Split(string(msg), "\n\n")
		if len(content) < 2 {
			continue
		}

		body := strings.Trim(content[1], "\x00")

		var res event
		if err := json.Unmarshal([]byte(body), &res); err != nil {
			continue
		}

		ch <- executions.Event{
			ExternalID: res.Token,
			Status:     res.Status,
		}
	}
}

type event struct {
	Token  string           `json:"token"`
	Status executions.State `json:"status"`
}

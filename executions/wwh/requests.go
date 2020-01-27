package wwh

import "encoding/json"

type request interface {
	headers() map[string]string
	body() ([]byte, error)
}

type createReq struct {
	username   string
	password   string
	properties []propertyReq
}

func (req createReq) headers() map[string]string {
	return map[string]string{
		"Content-Type": jsonCT,
		"username":     req.username,
		"password":     req.password,
	}
}

func (req createReq) body() ([]byte, error) {
	return json.Marshal(req.properties)
}

type propertyReq struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type executeReq struct {
	token                    string
	MetaResources            []string `json:"metaResources"`
	ChainComputation         string   `json:"chainComputation"`
	AdditionalLocalJobArgs   []string `json:"additionalLocalJobArguments"`
	LogicJarPath             string   `json:"logicJarPath"`
	RequestID                string   `json:"requestId"`
	Type                     string   `json:"type"`
	ResultPath               string   `json:"resultPath"`
	GlobalTimeout            string   `json:"globalTimeout"`
	LocalTimeout             string   `json:"localTimeout"`
	AdditionalPreprocessArgs []string `json:"additionalPreprocessArguments"`
	JobMode                  string   `json:"jobMode"`
	AdditionalGlobalJobArgs  []string `json:"additionalGlobalJobArguments"`
	AdditionalFiles          []string `json:"additionalFiles"`
	TrainedModelToken        string   `json:"trainedModelToken,omitempty"`
	TrainedModelName         string   `json:"trainedModelName,omitempty"`
}

func (req executeReq) headers() map[string]string {
	return map[string]string{
		"Content-Type":  jsonCT,
		"Authorization": req.token,
	}
}

func (req executeReq) body() ([]byte, error) {
	return json.Marshal(req)
}

type tokenReq struct {
	token string
}

func (req tokenReq) headers() map[string]string {
	return map[string]string{
		"Authorization": req.token,
	}
}

func (req tokenReq) body() ([]byte, error) {
	return nil, nil
}

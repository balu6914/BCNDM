package kubeflow

import "encoding/json"

type request interface {
	headers() map[string]string
	body() ([]byte, error)
}

type createJobRequest struct {
	Name         string              `json:"name"`
	StorageState string              `json:"storage_state"`
	PipelineSpec pipelineSpecRequest `json:"pipeline_spec"`
}

type pipelineSpecRequest struct {
	ID         string         `json:"pipeline_id"`
	Parameters []paramRequest `json:"parameters"`
}

type paramRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (req createJobRequest) headers() map[string]string {
	return nil
}

func (req createJobRequest) body() ([]byte, error) {
	return json.Marshal(req)
}

package kubeflow

type createAlgoResponse struct {
	ID string `json:"id"`
}

type execRunResponse struct {
	Run runResponse `json:"run"`
}

type runResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

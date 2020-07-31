package argo

type createAlgoResponse struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

type execRunResponse struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}

type runResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

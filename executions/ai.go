package executions

// AIService contains API definition for AI service.
type AIService interface {
	// CreateAlgorithm creates new algorithm on AI service.
	CreateAlgorithm(Algorithm) (Algorithm, error)

	// CreateDataset creates new dataset on AI service.
	CreateDataset(Dataset) (Dataset, error)

	// Start creates new execution and starts it.
	Start(Execution, Algorithm, Dataset) (Execution, error)

	// Result returns execution result if execution has finished.
	Result(Execution) (map[string]interface{}, error)

	// Events tracks events that are comming from AI system and send then
	// to the returned channel.
	Events() (chan Event, error)
}

// Event contains data that is received from event channel.
type Event struct {
	ExternalID string
	Status     State
}

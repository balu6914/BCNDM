package executions

// AIService contains API definition for AI service.
type AIService interface {
	// CreateAlgorithm creates new algorithm on AI service.
	CreateAlgorithm(Algorithm) error

	// CreateDataset creates new dataset on AI service.
	CreateDataset(Dataset) error

	// Start creates new execution and starts it.
	Start(Execution, Algorithm, Dataset) (string, error)

	// Result returns execution result if execution has finished.
	Result(string) (map[string]interface{}, error)

	// Events tracks events that are comming from AI system and send then
	// to the returned channel.
	Events() (chan Event, error)
}

// Event contains data that is received from event channel.
type Event struct {
	Token  string
	Status State
}

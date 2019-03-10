package executions

// AIService contains API definition for AI service.
type AIService interface {
	// CreateAlgorithm creates new algorithm on AI service.
	CreateAlgorithm(Algorithm) error

	// CreateDataset creates new dataset on AI service.
	CreateDataset(Dataset) error

	// Start creates new execution and starts it.
	Start(Execution, Algorithm, Dataset) (string, error)

	// IsDone checks if execution has finished.
	IsDone(string) (State, error)

	// Result returns execution result if execution has finished.
	Result(string) (map[string]interface{}, error)
}

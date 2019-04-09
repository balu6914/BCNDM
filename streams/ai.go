package streams

// AIService contains API definition for AI service.
type AIService interface {
	// CreateAlgorithm creates new algorithm on AI service.
	CreateAlgorithm(Stream) error
	// CreateDataset creates new dataset on AI service.
	CreateDataset(Stream) error
}

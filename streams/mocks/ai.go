package mocks

import "datapace/streams"

var _ streams.AIService = (*aiServiceMock)(nil)

type aiServiceMock struct{}

// NewAIService returns mock instance of AI service.
func NewAIService() streams.AIService {
	return aiServiceMock{}
}

func (ai aiServiceMock) CreateAlgorithm(_ streams.Stream) error {
	return nil
}

func (ai aiServiceMock) CreateDataset(_ streams.Stream) error {
	return nil
}

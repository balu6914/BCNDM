package executions

// AlgorithmRepository specifies algorithm persistence API.
type AlgorithmRepository interface {
	// Creates algorithm and stores it in persistant storage.
	Create(Algorithm) error

	// One finds and returns algorithm by specified id.
	One(string) (Algorithm, error)
}

// Algorithm contains algorithm metadata that is required by AI service.
type Algorithm struct {
	ID         string
	Name       string
	Path       string
	ModelToken string
	ModelName  string
}

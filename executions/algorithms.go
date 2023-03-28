package executions

// AlgorithmRepository specifies algorithm persistence API.
type AlgorithmRepository interface {
	// Creates algorithm and stores it in persistant storage.
	Create(Algorithm) error

	// Update replaces existing algorithm with the new value.
	Update(Algorithm) error

	// One finds and returns algorithm by specified id.
	One(string) (Algorithm, error)
}

// Algorithm contains algorithm metadata that is required by AI service.
type Algorithm struct {
	ID         string
	ExternalID string
	Name       string
	Metadata   map[string]string
}

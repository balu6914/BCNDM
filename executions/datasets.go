package executions

// DatasetRepository specifies dataset persistence API.
type DatasetRepository interface {
	// Creates dataset and stores it in persistant storage.
	Create(Dataset) error

	// Update replaces existing dataset with new value.
	Update(Dataset) error

	// One finds and returns dataset by specified id.
	One(string) (Dataset, error)
}

// Dataset contains dataset metadata that is required by AI service.
type Dataset struct {
	ID         string
	ExternalID string
	Metadata   map[string]string
}

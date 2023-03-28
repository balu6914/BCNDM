package executions

const (
	// InProgress represents execution in progress state.
	InProgress State = "In progress"
	// Failed reprents unsuccessfully finished execution state.
	Failed State = "Failed"
	// Succeeded reprents successfully finished execution state.
	Succeeded State = "Succeeded"
)

// State represents execution state.
type State string

// Execution contains execution metadata.
type Execution struct {
	ID         string
	ExternalID string
	Name       string
	Owner      string
	Algo       string
	Data       string
	Metadata   map[string]interface{}
	State      State
}

// ExecutionRepository specifies execution persistence API.
type ExecutionRepository interface {
	// Creates new execution into database.
	Create(Execution) (string, error)

	// Update replaces the existing execution with the new value.
	Update(Execution) error

	// UpdateState updates current execution state.
	UpdateState(string, State) error

	// Execution finds and returns execution by owner and id.
	Execution(string, string) (Execution, error)

	// List returns a list of executions.
	List(string) ([]Execution, error)
}

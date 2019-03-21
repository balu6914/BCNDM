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

const (
	// Federated represents federated execution mode.
	Federated JobMode = "FEDERATED"
	// Centralized represents contralized execution mode.
	Centralized JobMode = "CENTRALIZED"
	// Distributed Rrepresents distributed execution mode.
	Distributed JobMode = "DISTRIBUTED"
)

// JobMode represents execution job mod.
type JobMode string

// Execution contains execution metadata.
type Execution struct {
	ID                       string
	Owner                    string
	Algo                     string
	Data                     string
	AdditionalLocalJobArgs   []string
	Type                     string
	GlobalTimeout            uint64
	LocalTimeout             uint64
	AdditionalPreprocessArgs []string
	Mode                     JobMode
	AdditionalGlobalJobArgs  []string
	AdditionalFiles          []string
	State                    State
	Token                    string
}

// ExecutionRepository specifies execution persistence API.
type ExecutionRepository interface {
	// Creates new execution into database.
	Create(Execution) (string, error)

	// UpdateToken sets existing execution token.
	UpdateToken(string, string) error

	// UpdateState updates current execution state.
	UpdateState(string, State) error

	// Execution finds and returns execution by owner and id.
	Execution(string, string) (Execution, error)

	// List returns a list of executions.
	List(string) ([]Execution, error)
}

package executions

const (
	// Executing represents execution in progress state.
	Executing State = "executing"
	// Done reprents finished execution state.
	Done State = "done"
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

	// Finish marks existing execution as done.
	Finish(string) error

	// Execution finds and returns execution by owner and id.
	Execution(string, string) (Execution, error)

	// List returns a list of executions.
	List(string) ([]Execution, error)
}

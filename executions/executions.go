package executions

// ExecutionRepository specifies exeuction persistence API.
type ExecutionRepository interface {
	// Creates new execution into database.
	Create(string, string, string, JobMode) (string, error)

	// Finish marks existing execution as done.
	Finish(string) error

	// Execution finds and returns execution by owner and id.
	Execution(string, string) (Execution, error)

	// List returns a list of executions.
	List(string) ([]Execution, error)
}

package executions

import "errors"

var (
	// ErrConflict indicates that given execution already exists.
	ErrConflict = errors.New("execution already exists")

	// ErrNotFound indicates that required execution doesn't exist.
	ErrNotFound = errors.New("execution not found")

	// ErrMalformedData indicates that method receiver invalid input data.
	ErrMalformedData = errors.New("invalid data received")
)

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
	Federated JobMode = "federated"
	// Centralized represents contralized execution mode.
	Centralized JobMode = "centralized"
	// Distributed Rrepresents distributed execution mode.
	Distributed JobMode = "distributed"
)

// JobMode represents execution job mod.
type JobMode string

// Execution contains execution metadata.
type Execution struct {
	ID    string
	Owner string
	State State
	Algo  string
	Data  string
	Mode  JobMode
}

// Service contains executions service API specification.
type Service interface {
	// Start starts execution of specified algorithm on given data set.
	Start(string, string, string, JobMode) (string, error)

	// Finish marks exsiting execution as done.
	Finish(string) error

	// Execution returns one execution metadata.
	Execution(string, string) (Execution, error)

	// List returns all of the executions for given owner.
	List(string) ([]Execution, error)
}

var _ Service = (*executionsService)(nil)

type executionsService struct {
	repo ExecutionRepository
}

// NewService instantiates the domain service implementation.
func NewService(repo ExecutionRepository) Service {
	return executionsService{
		repo: repo,
	}
}

func (es executionsService) Start(owner, algo, data string, mode JobMode) (string, error) {
	return es.repo.Create(owner, algo, data, mode)
}

func (es executionsService) Finish(id string) error {
	return es.repo.Finish(id)
}

func (es executionsService) Execution(owner, id string) (Execution, error) {
	return es.repo.Execution(owner, id)
}

func (es executionsService) List(owner string) ([]Execution, error) {
	return es.repo.List(owner)
}

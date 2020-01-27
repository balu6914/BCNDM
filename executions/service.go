package executions

import (
	"errors"
)

var (
	// ErrConflict indicates that given execution already exists.
	ErrConflict = errors.New("execution already exists")

	// ErrNotFound indicates that required execution doesn't exist.
	ErrNotFound = errors.New("execution not found")

	// ErrMalformedData indicates that method receiver invalid input data.
	ErrMalformedData = errors.New("invalid data received")

	// ErrCreateAlgoFailed indicates executions service failed to create an
	// algorithm on AI service.
	ErrCreateAlgoFailed = errors.New("failed to create an algorithm")

	// ErrExecutionFailed indicated that execution has finished and that it has failed.
	ErrExecutionFailed = errors.New("specified execution has failed")
)

// Service contains executions service API specification.
type Service interface {
	// Start starts execution of specified algorithm on given data set.
	Start(Execution) (string, error)

	// Result returns result of the execution if execution has finished.
	Result(string, string) (map[string]interface{}, error)

	// Execution returns one execution metadata.
	Execution(string, string) (Execution, error)

	// List returns all of the executions for given owner.
	List(string) ([]Execution, error)

	// CreateAlgorithm creates new algorithm on external AI service.
	CreateAlgorithm(Algorithm) error

	// CreateDataset creates new dataset on external AI serivce.
	CreateDataset(Dataset) error

	// ProcessEvents processes state change events for algorithms and datasets.
	ProcessEvents() error
}

var _ Service = (*executionsService)(nil)

type executionsService struct {
	execs ExecutionRepository
	algos AlgorithmRepository
	data  DatasetRepository
	ai    AIService
}

// NewService instantiates the domain service implementation.
func NewService(execs ExecutionRepository, algos AlgorithmRepository, data DatasetRepository, ai AIService) Service {
	return executionsService{
		execs: execs,
		algos: algos,
		data:  data,
		ai:    ai,
	}
}

func (es executionsService) Start(exec Execution) (string, error) {
	algo, err := es.algos.One(exec.Algo)
	if err != nil {
		return "", err
	}

	data, err := es.data.One(exec.Data)
	if err != nil {
		return "", err
	}

	exec.State = InProgress
	id, err := es.execs.Create(exec)
	if err != nil {
		return "", err
	}
	exec.ID = id

	uexec, err := es.ai.Start(exec, algo, data)
	if err != nil {
		return "", err
	}

	if err := es.execs.Update(uexec); err != nil {
		return "", err
	}

	return id, nil
}

func (es executionsService) Result(owner, id string) (map[string]interface{}, error) {
	exec, err := es.execs.Execution(owner, id)
	if err != nil {
		return nil, err
	}

	if exec.State == InProgress {
		return nil, ErrNotFound
	}

	return es.ai.Result(exec)
}

func (es executionsService) Execution(owner, id string) (Execution, error) {
	return es.execs.Execution(owner, id)
}

func (es executionsService) List(owner string) ([]Execution, error) {
	return es.execs.List(owner)
}

func (es executionsService) CreateAlgorithm(algo Algorithm) error {
	if err := es.algos.Create(algo); err != nil {
		return err
	}

	ualgo, err := es.ai.CreateAlgorithm(algo)
	if err != nil {
		return err
	}

	return es.algos.Update(ualgo)
}

func (es executionsService) CreateDataset(dataset Dataset) error {
	if err := es.data.Create(dataset); err != nil {
		return err
	}

	udata, err := es.ai.CreateDataset(dataset)
	if err != nil {
		return err
	}

	return es.data.Update(udata)
}

func (es executionsService) ProcessEvents() error {
	ch, err := es.ai.Events()
	if err != nil {
		return err
	}

	for event := range ch {
		if err := es.execs.UpdateState(event.ExternalID, event.Status); err != nil {
			continue
		}
	}

	return nil
}

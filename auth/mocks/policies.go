package mocks

import (
	"sync"

	"github.com/datapace/datapace/auth"
	"gopkg.in/mgo.v2/bson"
)

type policyRepositoryMock struct {
	mu       *sync.Mutex
	policies map[string]auth.Policy
}

// NewPolicyRepository returns a new policy repository mock.
func NewPolicyRepository(policies map[string]auth.Policy, mu *sync.Mutex) auth.PolicyRepository {
	return &policyRepositoryMock{policies: policies, mu: mu}
}

func (prm *policyRepositoryMock) Save(policy auth.Policy) (string, error) {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	policy.ID = bson.NewObjectId().Hex()
	prm.policies[policy.ID] = policy
	return policy.ID, nil
}

func (prm *policyRepositoryMock) OneByID(id string) (auth.Policy, error) {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	if p, ok := prm.policies[id]; ok {
		return p, nil
	}
	return auth.Policy{}, auth.ErrNotFound
}

func (prm *policyRepositoryMock) OneByName(name string) (auth.Policy, error) {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	for _, p := range prm.policies {
		if p.Name == name {
			return p, nil
		}
	}

	return auth.Policy{}, auth.ErrNotFound
}

func (prm *policyRepositoryMock) List(owner string) ([]auth.Policy, error) {
	prm.mu.Lock()
	defer prm.mu.Unlock()
	ret := []auth.Policy{}
	for _, p := range prm.policies {
		if p.Owner == owner {
			ret = append(ret, p)
		}
	}
	return ret, nil
}

func (prm *policyRepositoryMock) ListByIDs(ids []string) ([]auth.Policy, error) {
	prm.mu.Lock()
	defer prm.mu.Unlock()
	ret := []auth.Policy{}
	for _, id := range ids {
		if p, ok := prm.policies[id]; ok {
			ret = append(ret, p)
		}
	}
	return ret, nil
}

func (prm *policyRepositoryMock) Remove(id string) error {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	if _, ok := prm.policies[id]; !ok {
		return auth.ErrNotFound
	}
	delete(prm.policies, id)
	return nil
}

func (prm *policyRepositoryMock) Attach(policyID, userID string) error {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	return nil
}

func (prm *policyRepositoryMock) Detach(policyID, userID string) error {
	prm.mu.Lock()
	defer prm.mu.Unlock()

	return nil
}

package groups

import "strconv"

type groupsRepositoryMock struct {
	storage map[Gid]GroupMetadata
}

var _ GroupsRepository = (*groupsRepositoryMock)(nil)

func NewGroupsRepositoryMock() GroupsRepository {
	return groupsRepositoryMock{make(map[Gid]GroupMetadata)}
}

func (g groupsRepositoryMock) Create(metadata GroupMetadata) (Gid, error) {
	n := len(g.storage)
	gid := Gid(strconv.Itoa(n))
	g.storage[gid] = metadata
	return gid, nil
}

func (g groupsRepositoryMock) GetAll() ([]Group, error) {
	all := make([]Group, 0, len(g.storage))
	for gid, md := range g.storage {
		all = append(all, Group{md, gid})
	}
	return all, nil
}

func (g groupsRepositoryMock) Get(gid Gid) (*Group, error) {
	md, exists := g.storage[gid]
	if !exists {
		return nil, ErrNotFound
	}
	return &Group{md, gid}, nil
}

func (g groupsRepositoryMock) Exists(gid Gid) (bool, error) {
	_, exists := g.storage[gid]
	return exists, nil
}

func (g groupsRepositoryMock) Delete(gid Gid) error {
	delete(g.storage, gid)
	return nil
}

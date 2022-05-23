package groups

type usersRepositoryMock struct {
	storage map[Uid]map[Gid]bool
}

func NewUsersRepositoryMock() UsersRepository {
	return usersRepositoryMock{make(map[Uid]map[Gid]bool)}
}

func (u usersRepositoryMock) Add(uid Uid, gid Gid) error {
	gids, exists := u.storage[uid]
	if !exists {
		gids = make(map[Gid]bool, 0)
	}
	gids[gid] = true
	u.storage[uid] = gids
	return nil
}

func (u usersRepositoryMock) Delete(uid Uid, gid Gid) error {
	gids, exists := u.storage[uid]
	if !exists {
		return ErrNotFound
	}
	delete(gids, gid)
	return nil
}

func (u usersRepositoryMock) GetGroups(uid Uid) ([]Gid, error) {
	gids, exists := u.storage[uid]
	if !exists {
		return []Gid{}, ErrNotFound
	}
	results := make([]Gid, 0, len(gids))
	for gid := range gids {
		results = append(results, gid)
	}
	return results, nil
}

func (u usersRepositoryMock) GetUsersPage(gid Gid, startOffset, limit int) (*UidsPage, error) {
	uids := make([]Uid, 0)
	offset := 0
	count := 0
	complete := true
	for uid, gids := range u.storage {
		if gids[gid] {
			offset++
			if offset > startOffset {
				uids = append(uids, uid)
				count++
			}
			if count == limit {
				complete = false
				break
			}
		}
	}
	page := &UidsPage{
		Complete: complete,
		Offset:   offset,
		Uids:     uids,
	}
	return page, nil
}

func (u usersRepositoryMock) DeleteAll(gid Gid) error {
	for _, gids := range u.storage {
		delete(gids, gid)
	}
	return nil
}

var _ UsersRepository = (*usersRepositoryMock)(nil)

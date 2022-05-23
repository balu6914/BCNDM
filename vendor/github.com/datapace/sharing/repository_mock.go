package sharing

import "strings"

type repositoryMock struct {
	storageRef *map[string]Receivers
}

func NewRepositoryMock(storageRef *map[string]Receivers) Repository {
	return repositoryMock{storageRef: storageRef}
}

func (r repositoryMock) Create(s Sharing) error {
	key := string(s.SourceUserId) + "_" + string(s.StreamId)
	storage := *r.storageRef
	if _, exists := storage[key]; exists {
		return ErrConflict
	}
	version := uint64(1)
	v := Receivers{
		Version:  &version,
		UserIds:  s.Receivers.UserIds,
		GroupIds: s.Receivers.GroupIds,
	}
	storage[key] = v
	return nil
}

func (r repositoryMock) Delete(srcUserId UserId, streamId StreamId) error {
	key := string(srcUserId) + "_" + string(streamId)
	storage := *r.storageRef
	if _, exists := storage[key]; !exists {
		return ErrNotFound
	}
	delete(storage, key)
	return nil
}

func (r repositoryMock) UpdateConditionally(s Sharing) error {
	key := string(s.SourceUserId) + "_" + string(s.StreamId)
	storage := *r.storageRef
	receivers, exists := storage[key]
	if !exists {
		return ErrNotFound
	}
	if *receivers.Version != *s.Receivers.Version {
		return ErrNotFound
	}
	receivers = s.Receivers
	newVersion := (*receivers.Version) + 1
	receivers.Version = &newVersion
	storage[key] = receivers
	return nil
}

func (r repositoryMock) GetOne(srcUserId UserId, streamId StreamId) (*Sharing, error) {
	key := string(srcUserId) + "_" + string(streamId)
	storage := *r.storageRef
	receivers, exists := storage[key]
	if !exists {
		return nil, ErrNotFound
	}
	return &Sharing{
		SourceUserId: srcUserId,
		StreamId:     streamId,
		Receivers:    receivers,
	}, nil
}

func (r repositoryMock) GetAllByReceiverUser(rcvUserId UserId) ([]Sharing, error) {
	var results []Sharing
	storage := *r.storageRef
	for key, receivers := range storage {
		for _, uid := range receivers.UserIds {
			if rcvUserId == uid {
				keyParts := strings.Split(key, "_")
				result := Sharing{
					SourceUserId: UserId(keyParts[0]),
					StreamId:     StreamId(keyParts[1]),
					Receivers:    receivers,
				}
				results = append(results, result)
				break
			}
		}
	}
	return results, nil
}

func (r repositoryMock) GetAllByReceiverGroups(rcvGroupIds []GroupId) ([]Sharing, error) {
	var results []Sharing
	storage := *r.storageRef
	for key, receivers := range storage {
		match := false
		for _, gid := range receivers.GroupIds {
			for _, rcvGroupId := range rcvGroupIds {
				if rcvGroupId == gid {
					keyParts := strings.Split(key, "_")
					result := Sharing{
						SourceUserId: UserId(keyParts[0]),
						StreamId:     StreamId(keyParts[1]),
						Receivers:    receivers,
					}
					results = append(results, result)
					match = true
					break
				}
			}
			if match {
				break
			}
		}
	}
	return results, nil
}

func (r repositoryMock) GetAllByReceiverUserAndGroups(rcvUserId UserId, rcvGroupIds []GroupId) ([]Sharing, error) {
	var results []Sharing
	storage := *r.storageRef
	for key, receivers := range storage {
		match := false
		for _, uid := range receivers.UserIds {
			if rcvUserId == uid {
				keyParts := strings.Split(key, "_")
				result := Sharing{
					SourceUserId: UserId(keyParts[0]),
					StreamId:     StreamId(keyParts[1]),
					Receivers:    receivers,
				}
				results = append(results, result)
				match = true
				break
			}
		}
		if match {
			continue
		}
		for _, gid := range receivers.GroupIds {
			for _, rcvGroupId := range rcvGroupIds {
				if rcvGroupId == gid {
					keyParts := strings.Split(key, "_")
					result := Sharing{
						SourceUserId: UserId(keyParts[0]),
						StreamId:     StreamId(keyParts[1]),
						Receivers:    receivers,
					}
					results = append(results, result)
					match = true
					break
				}
			}
			if match {
				break
			}
		}
	}
	return results, nil
}

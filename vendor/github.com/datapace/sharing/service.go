package sharing

import (
	"errors"
	"github.com/datapace/sharing/auth"
)

var (
	// ErrConflict indicates the Sharing entity already exists
	ErrConflict = errors.New("create failed, entity already exists")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")

	// ErrBadRequest indicates invalid request (e.g. missing mandatory parameters)
	ErrBadRequest = errors.New("bad request")
)

type (
	StreamId string

	UserId string

	GroupId string

	// Receivers contains set of receivers (sharing destination), including users and groups.
	Receivers struct {

		// Version is Receivers entity version for conditional update purpose.
		// May be omitted (set to nil) for the unconditional update purpose (create/delete/replace).
		Version *uint64 `json:"version,omitempty"`

		// UserIds contains ids of the user receivers
		UserIds []UserId `json:"userIds"`

		// GroupIds contains ids of the group receivers
		GroupIds []GroupId `json:"groupIds"`
	}

	// Sharing is the entity addressed by the source UserId (who shared) and StreamId (what shared)
	Sharing struct {

		// SourceUserId is the source user who shared
		SourceUserId UserId `json:"sourceUserId"`

		// StreamId is the stream id (what has been shared)
		StreamId StreamId `json:"streamId"`

		// Receivers is the destination (whom shared to)
		Receivers Receivers `json:"receivers"`
	}

	// Service is the sharing service level API
	Service interface {

		// UpdateReceivers sets the new Receivers set if the version matches. Otherwise returns ErrNotFound.
		UpdateReceivers(sharing Sharing) error

		// GetReceivers queries for the Receivers.
		// A client should use the entity Version in the UpdateReceivers function.
		GetReceivers(srcUserId UserId, streamId StreamId) (*Receivers, error)

		// GetSharingsToGroups queries for the all Sharing entities where any receiver is any of the specified groupIds.
		GetSharingsToGroups(rcvGroupIds []GroupId) ([]Sharing, error)

		// GetSharings queries for the all Sharing entities where any receiver is either:
		//	(*) specified UserId,
		//	(*) any of the specified groupIds.
		GetSharings(rcvUserId UserId, rcvGroupIds []GroupId) ([]Sharing, error)
	}

	service struct {
		authSvc auth.AuthService
		repo    Repository
	}
)

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (svc service) UpdateReceivers(sharing Sharing) error {
	receivers := sharing.Receivers
	for _, rcvUserId := range receivers.UserIds {
		if rcvUserId == sharing.SourceUserId {
			return ErrBadRequest // sharing to self should not be allowed
		}
	}
	ver := receivers.Version
	if ver == nil {
		return svc.createOrDelete(sharing)
	}
	return svc.repo.UpdateConditionally(sharing)
}

func (svc service) createOrDelete(sharing Sharing) error {
	if len(sharing.Receivers.UserIds) == 0 && len(sharing.Receivers.GroupIds) == 0 {
		return svc.repo.Delete(sharing.SourceUserId, sharing.StreamId)
	}
	return svc.repo.Create(sharing)
}

func (svc service) GetReceivers(srcUserId UserId, streamId StreamId) (*Receivers, error) {
	sharing, err := svc.repo.GetOne(srcUserId, streamId)
	if err != nil {
		return nil, err
	}
	return &sharing.Receivers, nil
}

func (svc service) GetSharingsToGroups(rcvGroupIds []GroupId) ([]Sharing, error) {
	return svc.repo.GetAllByReceiverGroups(rcvGroupIds)
}

func (svc service) GetSharings(rcvUserId UserId, rcvGroupIds []GroupId) ([]Sharing, error) {
	if rcvGroupIds == nil || len(rcvGroupIds) == 0 {
		return svc.repo.GetAllByReceiverUser(rcvUserId)
	}
	return svc.repo.GetAllByReceiverUserAndGroups(rcvUserId, rcvGroupIds)
}

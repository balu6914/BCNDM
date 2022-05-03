package groups

import (
	"errors"
	"github.com/datapace/groups/sharing"
)

// Gid is the group id
type Gid string

// Uid is the user id
type Uid string

// GroupMetadata contains the group details with the exemption of Gid and group users
type GroupMetadata struct {

	// Name is the unique Group name
	Name string `json:"name"`

	// Role is the group role
	Role string `json:"role"`
}

// Group contains the group id and metadata with the exemption of its users
type Group struct {
	GroupMetadata

	// Group Id
	Gid Gid `json:"gid"`
}

// UidsPage contains the results page for the users listing
type UidsPage struct {

	// Complete represents if the page is last (end of results flag)
	Complete bool `json:"complete"`

	// Offset points to the next results page
	Offset int `json:"offset"`

	// Uids contains the list of Uid on the current results page
	Uids []Uid `json:"uids"`
}

var (
	// ErrConflict indicates the group/user already exists
	ErrConflict = errors.New("group already exists or user already added to group")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")

	// ErrBadRequest indicates invalid request (e.g. missing mandatory parameters)
	ErrBadRequest = errors.New("bad request")
)

// Service is the groups service level API
type Service interface {

	// Create a group with the specified GroupMetadata, returns new group id, error otherwise
	Create(GroupMetadata) (gid Gid, err error)

	// GetAll existing groups Gid array
	GetAll() ([]Gid, error)

	// GetMetadata returns GroupMetadata by its Gid
	GetMetadata(Gid) (*GroupMetadata, error)

	// Delete group by its Gid, also deletes all users from the group
	Delete(Gid) error

	// AddUser adds the user specified by its Uid to the Group specified by Gid
	AddUser(Uid, Gid) error

	// DeleteUser removes the user specified by its Uid from the Group specified by Gid
	DeleteUser(Uid, Gid) error

	// GetUsersPage returns a UidsPage who belong to the Group specified by the Gid
	GetUsersPage(gid Gid, offset, limit int) (*UidsPage, error)

	// GetUserGroups returns all groups those the specified user belongs to
	GetUserGroups(Uid) ([]Gid, error)
}

type groupsService struct {
	sharingSvc sharing.Service
	groupsRepo GroupsRepository
	usersRepo  UsersRepository
}

var _ Service = (*groupsService)(nil)

func NewService(sharingSvc sharing.Service, groupsRepo GroupsRepository, usersRepo UsersRepository) Service {
	return &groupsService{
		sharingSvc: sharingSvc,
		groupsRepo: groupsRepo,
		usersRepo:  usersRepo,
	}
}

func (gs groupsService) Create(md GroupMetadata) (Gid, error) {
	return gs.groupsRepo.Create(md)
}

func (gs groupsService) GetAll() ([]Gid, error) {
	all, err := gs.groupsRepo.GetAll()
	if err != nil {
		return nil, err
	}
	var gids []Gid
	for _, g := range all {
		gids = append(gids, g.Gid)
	}
	return gids, nil
}

func (gs groupsService) GetMetadata(gid Gid) (*GroupMetadata, error) {
	g, err := gs.groupsRepo.Get(gid)
	if err != nil {
		return nil, err
	}
	return &g.GroupMetadata, nil
}

func (gs groupsService) Delete(gid Gid) error {
	if err := gs.groupsRepo.Delete(gid); err != nil {
		return err
	}
	if err := gs.usersRepo.DeleteAll(gid); err != nil {
		return err
	}
	return gs.sharingSvc.DeleteReceivers(string(gid))
}

func (gs groupsService) AddUser(uid Uid, gid Gid) error {
	exists, err := gs.groupsRepo.Exists(gid)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return gs.usersRepo.Add(uid, gid)
}

func (gs groupsService) DeleteUser(uid Uid, gid Gid) error {
	return gs.usersRepo.Delete(uid, gid)
}

func (gs groupsService) GetUsersPage(gid Gid, offset, limit int) (*UidsPage, error) {
	exists, err := gs.groupsRepo.Exists(gid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}
	return gs.usersRepo.GetUsersPage(gid, offset, limit)
}

func (gs groupsService) GetUserGroups(uid Uid) ([]Gid, error) {
	return gs.usersRepo.GetGroups(uid)
}

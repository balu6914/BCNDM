package groups

// UsersRepository is the persistence interface for the users-to-groups mapping
type UsersRepository interface {

	// Add the user specified by groups.Uid to the group specified by groups.Gid
	Add(Uid, Gid) error

	// Delete the user specified by groups.Uid from the group specified by groups.Gid
	Delete(Uid, Gid) error

	// GetGroups returns all group's ids those the specified user belongs to
	GetGroups(Uid) ([]Gid, error)

	// GetUsersPage returns all user ids those belong to the group specified by the groups.Gid
	GetUsersPage(gid Gid, offset, limit int) (*UidsPage, error)

	// DeleteAll users from the group specified by the Gid
	DeleteAll(Gid) error
}

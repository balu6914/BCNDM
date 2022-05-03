package groups

type GroupsRepository interface {

	// Create a group with the specified groups.GroupMetadata, returns new groups.Gid, error otherwise
	Create(GroupMetadata) (Gid, error)

	// GetAll existing groups
	GetAll() ([]Group, error)

	// Get groups.Group by groups.Gid
	Get(Gid) (*Group, error)

	// Exists returns true if the group specified by Gid exists, false otherwise and error if was not able to check
	Exists(Gid) (bool, error)

	// Delete group by its groups.Gid
	Delete(Gid) error
}

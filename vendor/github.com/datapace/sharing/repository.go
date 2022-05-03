package sharing

type Repository interface {

	// Create a new Sharing entity, returns ErrConflict if exists already.
	Create(s Sharing) error

	// Delete the Sharing entity identified by source UserId and StreamId.
	Delete(srcUserId UserId, streamId StreamId) error

	// UpdateConditionally replaces the Receivers payload if supplied version matches the persisted one.
	// Returns ErrNotFound if Receivers version doesn't match or entity is not created yet.
	UpdateConditionally(s Sharing) error

	// GetOne returns single Sharing entity identified by source UserId and StreamId.
	GetOne(srcUserId UserId, streamId StreamId) (*Sharing, error)

	// GetAllByReceiverUser returns all Sharing entities where Receivers contains the specified UserId.
	GetAllByReceiverUser(rcvUserId UserId) ([]Sharing, error)

	// GetAllByReceiverGroups returns all Sharing entities where Receivers contains any of the specified group ids.
	GetAllByReceiverGroups(rcvGroupIds []GroupId) ([]Sharing, error)

	// GetAllByReceiverUserAndGroups returns all Sharing entities where Receivers contains either:
	//	(*) specified UserId,
	//	(*) any of the specified group ids.
	GetAllByReceiverUserAndGroups(rcvUserId UserId, rcvGroupIds []GroupId) ([]Sharing, error)
}

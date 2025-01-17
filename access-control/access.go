package access

const (
	// Pending represenets pending state of access request.
	Pending State = "pending"
	// Approved represents accepted state of access request.
	Approved State = "approved"
	// Revoked represents revoked state of access request.
	Revoked State = "revoked"
)

// State represents access request state.
type State string

// Request contains access request metadata.
type Request struct {
	ID       string
	Sender   string
	Receiver string
	State    State
}

// RequestRepository specifies access request persistance API.
type RequestRepository interface {
	// RequestAccess creates new connection between users that has pending
	// state.
	RequestAccess(string, string) (string, error)

	// ListSent returns a list of access requests that user with
	// specified id has sent.
	ListSent(string, State) ([]Request, error)

	// ListReceived returns a list of access requests that user
	// with specified id has received.
	ListReceived(string, State) ([]Request, error)

	// Approve updates status of access request to approved.
	Approve(string, string) error

	// Revoke updates status of access request to revoked.
	Revoke(string, string) error

	// One finds and returns access request by it's id.
	One(string) (Request, error)

	// GrantAccess combines the RequestAccess and Approve in a single call.
	// Returns the access request id to revoke later when needed.
	GrantAccess(srcUserId, dstUserId string) (string, error)

	// MigrateSameSenderAndReceiver scans the whole table and finds all records having same sender and receiver. These
	// records are then considered as a single access request record and discarded except the latest record id.
	MigrateSameSenderAndReceiver() error
}

// RequestLedger specifies access request writer API.
type RequestLedger interface {
	// RequestAccess creates new connection between users that has pending
	// state.
	RequestAccess(string, string) error

	// Approve updates status of access request to approved.
	Approve(string, string) error

	// Revoke updates status of access request to revoked.
	Revoke(string, string) error

	// Grant the access to the desitnation.
	Grant(src, dst string) error
}

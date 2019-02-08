package auth

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

// AccessRequest contains access request metadata.
type AccessRequest struct {
	ID       string
	Sender   string
	Receiver string
	State    State
}

// AccessRequestRepository specifies access request persistance API.
type AccessRequestRepository interface {
	// RequestAccess creates new connection between users that has pending
	// state.
	RequestAccess(string, string) (string, error)

	// ListSentAccessRequests returns a list of access requests that user with
	// specified id has sent.
	ListSentAccessRequests(string, State) ([]AccessRequest, error)

	// ListReceivedAccessRequests returns a list of access requests that user
	// with specified id has received.
	ListReceivedAccessRequests(string, State) ([]AccessRequest, error)

	// ApproveAccessRequest updates status of access request to approved.
	ApproveAccessRequest(string, string) error

	// RejectAccessRequest updates status of access request to rejected.
	RejectAccessRequest(string, string) error
}

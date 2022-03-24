package access

import "errors"

var (
	// ErrConflict indicates that sent access request already exists and is not
	// revoked.
	ErrConflict = errors.New("request already sent")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid id format).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// RequestAccess creates new connection between users that has pending
	// state.
	RequestAccess(string, string) (string, error)

	// ListSentAccessRequests returns a list of access requests that user with
	// specified key has sent.
	ListSentAccessRequests(string, State) ([]Request, error)

	// ListReceivedAccessRequests returns a list of access requests that user
	// with specified key has received.
	ListReceivedAccessRequests(string, State) ([]Request, error)

	// ListPartners fetches list of users that gave you access to their content.
	ListPartners(string) ([]string, error)

	// ListPotentialPartners fetches list of potential and actual partners.
	ListPotentialPartners(string) ([]string, error)

	// ApproveAccessRequest updates status of access request to approved.
	ApproveAccessRequest(string, string) error

	// RevokeAccessRequest updates status of access request to revoked.
	RevokeAccessRequest(string, string) error

	// GrantAccess combines both request access and approve it in a single call.
	// Returns the access request id to revoke later, if needed.
	GrantAccess(string, string) (string, error)
}

var _ Service = (*accessControlService)(nil)

type accessControlService struct {
	auth   AuthService
	repo   RequestRepository
	ledger RequestLedger
}

// New returns new access control service instance.
func New(auth AuthService, repo RequestRepository, ledger RequestLedger) Service {
	return accessControlService{
		auth:   auth,
		repo:   repo,
		ledger: ledger,
	}
}

func (acs accessControlService) RequestAccess(key, partner string) (string, error) {
	if err := acs.auth.Exists(partner); err != nil {
		return "", ErrNotFound
	}

	id, err := acs.auth.Identify(key)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := acs.ledger.RequestAccess(id, partner); err != nil {
		return "", err
	}

	return acs.repo.RequestAccess(id, partner)
}

func (acs accessControlService) ListSentAccessRequests(key string, state State) ([]Request, error) {
	id, err := acs.auth.Identify(key)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return acs.repo.ListSent(id, state)
}

func (acs accessControlService) ListReceivedAccessRequests(key string, state State) ([]Request, error) {
	id, err := acs.auth.Identify(key)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return acs.repo.ListReceived(id, state)
}

func (acs accessControlService) ListPartners(id string) ([]string, error) {
	requests, err := acs.repo.ListSent(id, Approved)
	if err != nil {
		return nil, err
	}

	partners := []string{}
	for _, req := range requests {
		partners = append(partners, req.Receiver)
	}

	return partners, nil
}

func (acs accessControlService) ListPotentialPartners(id string) ([]string, error) {
	pendingList, err := acs.repo.ListSent(id, Pending)
	if err != nil {
		return nil, err
	}

	approvedList, err := acs.repo.ListSent(id, Approved)
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, req := range pendingList {
		list = append(list, req.Receiver)
	}
	for _, req := range approvedList {
		list = append(list, req.Receiver)
	}

	return list, nil
}

func (acs accessControlService) ApproveAccessRequest(key, id string) error {
	uid, err := acs.auth.Identify(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	req, err := acs.repo.One(id)
	if err != nil {
		return err
	}

	if err := acs.ledger.Approve(uid, req.Sender); err != nil {
		return err
	}

	return acs.repo.Approve(uid, id)
}

func (acs accessControlService) RevokeAccessRequest(key, id string) error {
	uid, err := acs.auth.Identify(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	req, err := acs.repo.One(id)
	if err != nil {
		return err
	}

	if err := acs.ledger.Revoke(uid, req.Sender); err != nil {
		return err
	}

	return acs.repo.Revoke(uid, id)
}

func (acs accessControlService) GrantAccess(key, dstUserId string) (string, error) {
	srcUserId, err := acs.auth.Identify(key)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	if err := acs.ledger.GrantAccess(srcUserId, dstUserId); err != nil {
		return "", err
	}
	return acs.repo.GrantAccess(srcUserId, dstUserId)
}

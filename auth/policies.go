package auth

// Role constants
const (
	AdminRole       = "admin"
	AdminUserRole   = "admin_user"
	AdminWalletRole = "admin_wallet"
	UserRole        = "user"
	BuyerRole       = "buyer"
	SellerRole      = "seller"
)

// Policy consists of rules and constranits to validate Resource against.
type Policy struct {
	ID          string
	Version     string
	Name        string
	Owner       string
	Rules       []Rule
	Constraints []Constraint
}

// Attributes method returns Policy attributes.
func (p Policy) Attributes() map[string]string {
	return map[string]string{
		"id":      p.ID,
		"version": p.Version,
		"name":    p.Name,
		"owner":   p.Owner,
	}
}

// ResourceType returns Policy resource type.
func (p Policy) ResourceType() string {
	return "policy"
}

// Validate validates the access request against the Policy.
func (p Policy) Validate(subject Resource, action Action, object Resource) bool {
	rt := object.ResourceType()
	for _, rule := range p.Rules {
		if rule.Type == rt && (rule.Action == Any || rule.Action == action) {
			// Condition can be omitted.
			if rule.Condition == nil || rule.Condition.Validate(subject, object) {
				return true
			}
		}
	}
	return false
}

// Constraint is Policy constraint.
type Constraint interface {
	Validate(r Resource)
}

// Subject executes actions. Policies are attached to the subject.
type Subject interface {
	Policies() []Policy
}

// Resource is a resource an action is executed over.
// Resource prefix is added to avoid naming colision.
// For example, ID is common field name, and having
// field and method with the same name is not allowed.
type Resource interface {
	Attributes() map[string]string
	ResourceType() string
}

// PolicyRepository exposes Policy persistence API.
type PolicyRepository interface {
	// Save a single policy.
	Save(policy Policy) (string, error)

	// OneByID retrieves the Policy by its ID.
	OneByID(id string) (Policy, error)

	// OneByName retrieves the Policy by its name.
	OneByName(name string) (Policy, error)

	// List returns all the policies that belong to the owner.
	List(owner string) ([]Policy, error)

	// ListByIDs retrieves the Policy by its ID.
	ListByIDs(ids []string) ([]Policy, error)

	// Remove an existing policy.
	Remove(id string) error

	// Attach adds policy to the user.
	Attach(policyID, userID string) error

	// Detach removes policy from the user.
	Detach(policyID, userID string) error
}

package auth

import "time"

var _ Resource = (*User)(nil)

// User represents a Datapace user account. Each user is identified given its
// email and password.
type User struct {
	Email                     string
	ContactEmail              string
	Password                  string
	ID                        string
	FirstName                 string
	LastName                  string
	Company                   string
	Address                   string
	Country                   string
	Mobile                    string
	Phone                     string
	Role                      string
	Disabled                  bool
	Policies                  []Policy
	Locked                    bool
	Attempt                   int
	PasswordHistory           []string
	PasswordResetSecret       string
	PasswordResetTokenExpires int64
	CreatedDate               *time.Time
	Metadata                  map[string]interface{}
}

type AdminFilters struct {
	Roles    []string
	Locked   bool
	Disabled bool
}

// Attributes returns user's attributes.
func (u User) Attributes() map[string]string {
	return map[string]string{
		"id":           u.ID,
		"email":        u.Email,
		"company":      u.Company,
		"firstName":    u.FirstName,
		"lastName":     u.LastName,
		"address":      u.Address,
		"country":      u.Country,
		"mobile":       u.Mobile,
		"phone":        u.Phone,
		"contactEmail": u.ContactEmail,
		"role":         u.Role,
	}
}

// ResourceType returns User type string value.
func (u User) ResourceType() string {
	return "user"
}

// UserRepository specifies an account persistence API.
type UserRepository interface {
	// Save persists the user account. A non-nil error is returned to indicate
	// operation failure.
	Save(User) (string, error)

	// Retrieves user by its ID.
	OneByID(string) (User, error)

	// Retrieves user by its Email.
	OneByEmail(string) (User, error)

	// Updates user by its unique identifier.
	Update(User) error

	// Removes user by its unique identifier.
	Remove(string) error

	// List all users that are not in the specified list.
	AllExcept([]string) ([]User, error)

	// RetrieveAll retrieves a list with all users filtered by admin role.
	RetrieveAll(filers AdminFilters) ([]User, error)
}

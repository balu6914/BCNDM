package auth

// User represents a Monetasa user account. Each user is identified given its
// email and password.
type User struct {
	Email        string `json:"email"`
	ContactEmail string `json:"contact_email,omitempty"`
	Password     string `json:"password"`
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
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
}

package transactions

// User contains user account data.
type User struct {
	ID     string
	Secret string
}

// UserRepository defines API for user account persistance.
type UserRepository interface {
	// Save creates new user and saves it in storage.
	Save(User) error

	// Remove removes existing user from storage.
	Remove(string) error
}

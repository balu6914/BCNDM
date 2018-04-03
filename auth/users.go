package auth

import "github.com/asaskevich/govalidator"

// User represents a Mainflux user account. Each user is identified given its
// email and password.
type User struct {
	Email    string `gorm:"type:varchar(254);primary_key"`
	Password string `gorm:"type:char(60)"`
}

// Validate returns an error if user representation is invalid.
func (u *User) Validate() error {
	if u.Email == "" || u.Password == "" {
		return ErrMalformedEntity
	}

	if !govalidator.IsEmail(u.Email) {
		return ErrMalformedEntity
	}

	return nil
}

// UserRepository specifies an account persistence API.
type UserRepository interface {
	// Save persists the user account. A non-nil error is returned to indicate
	// operation failure.
	Save(User) error

	// One retrieves user by its unique identifier.
	One(string) (User, error)

	// All retrieves all users.
	All(string) (User, error)

	// Update updates user by its unique identifier.
	Update(string) (User, error)

	// Remove removes user by its unique identifier.
	Remove(string) (User, error)
}

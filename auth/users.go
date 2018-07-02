package auth

import (
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

// User represents a Monetasa user account. Each user is identified given its
// email and password.
type User struct {
	Email    string        `json:"email"`
	Password string        `json:"password"`
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Balance  uint64        `json:"balance"`
	PubCert  []byte        `json:"pub_cert"`
	PrivCert []byte        `json:"priv_cert"`
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

	// Retrieves user by its ID.
	OneByID(string) (User, error)

	// Retrieves user by its Email.
	OneByEmail(string) (User, error)

	// Update updates user by its unique identifier.
	Update(User) error

	// Remove removes user by its unique identifier.
	Remove(string) error
}

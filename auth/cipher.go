package auth

// Cipher interface implements Encrypt and Decrypt methods for private User data.
type Cipher interface {
	// Encrypt encrypts private User data.
	Encrypt(User) (User, error)

	// Decrypt decrypts private User data.
	Decrypt(User) (User, error)
}

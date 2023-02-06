package aes

import (
	"crypto/aes"
	c "crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/datapace/datapace/auth"
)

type aesCipher struct {
	key []byte
}

// NewCipher creates a new AES cipher.
func NewCipher(key []byte) auth.Cipher {
	return aesCipher{
		key: key,
	}
}

func (cipher aesCipher) Encrypt(user auth.User) (auth.User, error) {
	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		return auth.User{}, err
	}

	gcm, err := c.NewGCM(block)
	if err != nil {
		return auth.User{}, nil
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return auth.User{}, err
	}

	if user.Company != "" {
		user.Company = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(user.Company), nil))
	}

	if user.Address != "" {
		user.Address = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(user.Address), nil))
	}

	if user.Country != "" {
		user.Country = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(user.Country), nil))
	}

	if user.Mobile != "" {
		user.Mobile = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(user.Mobile), nil))
	}

	if user.Phone != "" {
		user.Phone = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(user.Phone), nil))
	}

	return user, nil
}

func (cipher aesCipher) Decrypt(user auth.User) (auth.User, error) {
	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		return auth.User{}, err
	}

	gcm, err := c.NewGCM(block)
	if err != nil {
		return auth.User{}, nil
	}

	if user.Address != "" {
		if user.Address, err = decrypt(user.Address, gcm); err != nil {
			return auth.User{}, err
		}
	}

	if user.Country != "" {
		if user.Country, err = decrypt(user.Country, gcm); err != nil {
			return auth.User{}, err
		}
	}

	if user.Company != "" {
		if user.Company, err = decrypt(user.Company, gcm); err != nil {
			return auth.User{}, err
		}
	}

	if user.Mobile != "" {
		if user.Mobile, err = decrypt(user.Mobile, gcm); err != nil {
			return auth.User{}, err
		}
	}

	if user.Phone != "" {
		if user.Phone, err = decrypt(user.Phone, gcm); err != nil {
			return auth.User{}, err
		}
	}

	return user, nil
}

func decrypt(cipherText string, gcm c.AEAD) (string, error) {
	nonceSize := gcm.NonceSize()

	bytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonce, input := bytes[:nonceSize], bytes[nonceSize:]
	plain, err := gcm.Open(nil, nonce, input, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

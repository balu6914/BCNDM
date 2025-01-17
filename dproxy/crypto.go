package dproxy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var errInvalidCyphertext = errors.New("ciphertext block size is too short!")

func encrypt(key []byte, data string) (string, error) {
	plainText := []byte(data)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(key []byte, data string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errInvalidCyphertext
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}

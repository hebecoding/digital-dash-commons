package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type EncryptionService struct {
	gcm cipher.AEAD
}

func NewEncryptionService(key []byte) (*EncryptionService, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &EncryptionService{gcm: gcm}, nil
}

func (e *EncryptionService) Encrypt(input string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("input string cannot be empty")
	}

	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := e.gcm.Seal(nonce, nonce, []byte(input), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *EncryptionService) Decrypt(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < e.gcm.NonceSize() {
		return "", errors.New("encrypted string is too short")
	}

	nonce, ciphertext := ciphertext[:e.gcm.NonceSize()], ciphertext[e.gcm.NonceSize():]
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

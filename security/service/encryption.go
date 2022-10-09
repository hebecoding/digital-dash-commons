package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

type EncryptionService struct {
	Key string
}

const (
	encryptionFailed = "failed to encrypt value"
	decryptionFailed = "failed to decrypt value"
)

func NewEncryptionService(key string) *EncryptionService {
	return &EncryptionService{Key: key}
}

func (s *EncryptionService) Encrypt(data string) (string, error) {

	block, err := aes.NewCipher([]byte(s.Key))
	if err != nil {
		log.Println(encryptionFailed)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(encryptionFailed)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(encryptionFailed)
		return "", err
	}

	encryptedData := gcm.Seal(nonce, nonce, []byte(data), nil)

	return string(encryptedData), nil
}

func (s *EncryptionService) Decrypt(data string) (string, error) {
	dataBytes := []byte(data)

	key := []byte(s.Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(decryptionFailed)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(decryptionFailed)
		return "", err
	}

	nonceSize := gcm.NonceSize()

	nonce, ciphertext := dataBytes[:nonceSize], dataBytes[nonceSize:]

	decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(decryptionFailed)
		return "", err

	}
	return string(decryptedText), nil

}

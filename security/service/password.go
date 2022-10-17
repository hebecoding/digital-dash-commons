package service

import (
	errs "digital-dash-commons/error"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	EncryptionCost int
}

func NewPasswordService(encryptionCost int) *PasswordService {
	return &PasswordService{EncryptionCost: encryptionCost}
}

func (p *PasswordService) Encrypt(password string) ([]byte, error) {
	if password == "" {
		return nil, errs.ErrEmptyInput
	}
	return bcrypt.GenerateFromPassword([]byte(password), p.EncryptionCost)
}

func (p *PasswordService) IsPassword(pwHash, password []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(pwHash, password); err != nil {
		return false, err
	}
	return true, nil
}

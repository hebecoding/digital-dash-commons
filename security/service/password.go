package service

import (
	"errors"
	"fmt"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"strings"
	"time"
)

type Config struct {
	EncryptionCost      int
	MinPasswordLen      int
	MaxPasswordAttempts int
	Blacklist           []string
	StrengthChecker     func(password string) error
}

type Option func(*Config)

func WithEncryptionCost(encryptionCost int) Option {
	return func(c *Config) {
		c.EncryptionCost = encryptionCost
	}
}

func WithMinPasswordLen(minPasswordLen int) Option {
	return func(c *Config) {
		c.MinPasswordLen = minPasswordLen
	}
}

func WithMaxPasswordAttempts(maxPasswordAttempts int) Option {
	return func(c *Config) {
		c.MaxPasswordAttempts = maxPasswordAttempts
	}
}

func WithBlacklist(blacklist []string) Option {
	return func(c *Config) {
		c.Blacklist = blacklist
	}
}

func WithStrengthChecker(strengthChecker func(password string) error) Option {
	return func(c *Config) {
		c.StrengthChecker = strengthChecker
	}
}

type PasswordService struct {
	config      Config
	rateLimiter *rate.Limiter
}

func NewPasswordService(opts ...Option) *PasswordService {
	cfg := &Config{
		EncryptionCost:      bcrypt.DefaultCost,
		MinPasswordLen:      8,
		MaxPasswordAttempts: 5,
		Blacklist:           []string{},
		StrengthChecker:     defaultStrengthChecker,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	rateLimiter := rate.NewLimiter(rate.Every(time.Minute), cfg.MaxPasswordAttempts)

	return &PasswordService{config: *cfg, rateLimiter: rateLimiter}
}

func (p *PasswordService) Encrypt(password string) ([]byte, error) {
	if err := p.config.StrengthChecker(password); err != nil {
		return nil, err
	}

	return bcrypt.GenerateFromPassword([]byte(password), p.config.EncryptionCost)
}

func (p *PasswordService) VerifyPassword(pwHash, password []byte) (bool, error) {
	if !p.rateLimiter.Allow() {
		return false, errors.New("too many password attempts")
	}

	err := bcrypt.CompareHashAndPassword(pwHash, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, fmt.Errorf("passwords do not match")
		}
		return false, err
	}
	return true, nil
}

func defaultStrengthChecker(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if strings.Contains(password, "password") || strings.Contains(password, "123456") {
		return fmt.Errorf("password is too weak")
	}

	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()-_=+{}[]|;:'\",.<>?/`~")

	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return fmt.Errorf("password must contain at least one uppercase, one lowercase, one digit, and one special character")
	}

	return nil
}

func goPasswordValidatorStrengthChecker(minEntropyBits float64) func(password string) error {
	return func(password string) error {
		entropyBits := passwordvalidator.GetEntropy(password)
		if entropyBits < minEntropyBits {
			return fmt.Errorf("password is too weak (entropy bits: %.2f), must have at least %.2f entropy bits", entropyBits, minEntropyBits)
		}
		return nil
	}
}

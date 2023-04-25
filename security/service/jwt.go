package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type TokenManager struct {
	config JWTConfig
}

type CustomClaims struct {
	jwt.StandardClaims
	UserID      string   `json:"user_id"`
	TenantID    string   `json:"tenant_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func NewTokenManager(config JWTConfig) *TokenManager {
	return &TokenManager{
		config: config,
	}
}

func (tm *TokenManager) GenerateAccessToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tm.config.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (tm *TokenManager) GenerateRefreshToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tm.config.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (tm *TokenManager) ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

package models

import "github.com/golang-jwt/jwt/v4"

type UserAUTHClaims struct {
	jwt.RegisteredClaims
	UserID    string    `json:"id"`
	UserRoles UserRoles `json:"userRoles,omitempty"`
	Payload   []byte    `json:"payload,omitempty"`
}

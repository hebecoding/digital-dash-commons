package service

import (
	"digital-dash-commons/security/models"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type JWTTokenService struct {
	issuer     string
	subject    string
	audience   string
	expiry     *time.Time
	privateKey []byte
	publicKey  []byte
}

type JWTConfiguration struct {
	Issuer     string
	Audience   string
	Expiry     *time.Time
	PublicKey  []byte
	PrivateKey []byte
}

var (
	errTokenTypeConversion = errors.New("unable to assert token claims type")
	errEmptyClaimsID       = errors.New("userID cannot be blank in claims")
)

func NewJWTTokenService(config *JWTConfiguration) *JWTTokenService {

	jwtTokenService := &JWTTokenService{
		issuer:     config.Issuer,
		audience:   config.Audience,
		expiry:     config.Expiry,
		privateKey: config.PrivateKey,
		publicKey:  config.PrivateKey,
	}

	switch {

	case config.PrivateKey == nil || config.PublicKey == nil:
		log.Fatal("missing private and public keys for token service.")

	case config.Issuer == "":
		jwtTokenService.issuer = "Digital Dash Security Service"
		fallthrough

	case config.Audience == "":
		jwtTokenService.audience = "Digital Dash Clients"
		fallthrough

	case config.Expiry == nil:
		expiry := time.Now().Add(15 * time.Minute)
		jwtTokenService.expiry = &expiry
	}

	return jwtTokenService
}

func (s *JWTTokenService) CreateToken(claims *models.UserAUTHClaims) (string, error) {
	privKey, err := jwt.ParseEdPrivateKeyFromPEM(s.privateKey)
	if err != nil {
		return "", err
	}

	claims.RegisteredClaims.Subject = claims.UserID

	// validate claims struct
	if err = s.validateTokenClaims(claims); err != nil {
		return "", err
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	return jwtClaims.SignedString(privKey)
}

func (s *JWTTokenService) ParseToken(token string) (*models.UserAUTHClaims, error) {

	var userClaims models.UserAUTHClaims

	pubKey, err := jwt.ParseEdPublicKeyFromPEM(s.publicKey)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.ParseWithClaims(token, &userClaims, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, err
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*models.UserAUTHClaims)
	if !ok {
		return nil, errTokenTypeConversion
	}

	return claims, nil
}

func (s *JWTTokenService) validateTokenClaims(claims *models.UserAUTHClaims) error {
	switch {
	case claims.UserID == "":
		return errEmptyClaimsID
	}

	return nil
}

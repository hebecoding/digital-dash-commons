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
	Subject    string
	Audience   string
	Expiry     *time.Time
	PublicKey  []byte
	PrivateKey []byte
}

var (
	errTokenTypeConversion = errors.New("unable to assert token claims type")
)

func NewJWTTokenService(config *JWTConfiguration) *JWTTokenService {

	jwtTokenService := &JWTTokenService{}

	switch {

	case config.PrivateKey == nil || config.PublicKey == nil:
		log.Fatal("missing private and public keys for token service.")

	case config.Issuer == "":
		jwtTokenService.issuer = "Digital Dash Security Service"
		fallthrough

	case config.Subject == "":
		jwtTokenService.subject = "Digital Dash JWT"
		fallthrough

	case config.Audience == "":
		jwtTokenService.audience = "Digital Dash Clients"
		fallthrough

	case config.Expiry == nil:
		expiry := time.Now().Add(30 * time.Minute)
		jwtTokenService.expiry = &expiry
	}

	jwtTokenService.privateKey = config.PrivateKey
	jwtTokenService.publicKey = config.PublicKey

	return jwtTokenService
}

func (s *JWTTokenService) CreateToken(claims *models.UserAUTHClaims) (string, error) {
	privKey, err := jwt.ParseEdPrivateKeyFromPEM(s.privateKey)
	if err != nil {
		return "", err
	}

	// validate claims struct
	s.validateTokenClaims(claims)

	//set default token claims if not preset
	s.setDefaultClaims(claims)

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
		return errors.New("userID cannot be blank in claims")
	}

	return nil
}

func (s *JWTTokenService) setDefaultClaims(claims *models.UserAUTHClaims) {
	switch {

	case claims.ExpiresAt == nil:
		claims.ExpiresAt = jwt.NewNumericDate(*s.expiry)
		fallthrough

	case claims.Issuer == "":
		claims.Issuer = s.issuer
		fallthrough

	case claims.Subject == "":
		claims.Subject = s.subject
		fallthrough

	case claims.Audience == nil:
		claims.Audience = []string{s.audience}
	}
}

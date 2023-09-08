package api

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	AccessKey string
)

// generateTokenPair generates access and refresh tokens with SHA512 algorithm
func generateTokenPair(guid string) (map[string]string, error) {
	// create access token
	access := jwt.New(jwt.SigningMethodHS512)

	// set access token claims
	aClaims := access.Claims.(jwt.MapClaims)
	aClaims["sub"] = guid
	aClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// sign access token
	a, err := access.SignedString([]byte(AccessKey))
	if err != nil {
		return nil, err
	}

	// create refresh token
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(guid), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// return pair of tokens
	return map[string]string{
		"access":  a,
		"refresh": string(refreshToken),
	}, nil
}

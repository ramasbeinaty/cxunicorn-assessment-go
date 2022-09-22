package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID         int
	Name           string
	Email          string
	Role           string
	StandardClaims jwt.StandardClaims
}

func (c Claims) verifyRole(required bool) bool {
	if c.Role == "" {
		return !required
	}
	return true
}

func (c Claims) Valid() error {
	// validate token claims
	var now = time.Now().UTC().Unix()

	if c.StandardClaims.VerifyExpiresAt(now, true) == false {
		return errors.New("Failed to validate token claims - Expiry At field is not defined")
	}

	if c.verifyRole(true) == false {
		return errors.New("Failed to validate token claims - Role field is not defined")
	}

	return nil
}

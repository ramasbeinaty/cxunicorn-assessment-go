package auth

import (
	"errors"
	"os"
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

func (s *service) GenerateJWT(claims *Claims) (string, error) {
	tokenStr := ""

	SECRET_KEY := os.Getenv("SECRET_KEY")
	EXPIRY_DURATION := os.Getenv("EXPIRY_DURATION")

	// parse and define an expiration duration of the token
	_expiryDuration, err := time.ParseDuration(EXPIRY_DURATION)

	if err != nil {
		return tokenStr, errors.New("ERROR: GenerateJWT - Failed to parse expiry string to type duration -" + err.Error())
	}

	expiresAt := time.Now().UTC().Add(time.Minute * time.Duration(_expiryDuration)).Unix()

	// update the token claim with the expiration time
	claims.StandardClaims.ExpiresAt = expiresAt

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err = _token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return tokenStr, errors.New("GenerateJWT - " + err.Error())
	}

	return tokenStr, nil
}

func (s *service) GetTokenFromString(tokenStr string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		SECRET_KEY := os.Getenv("SECRET_KEY")

		// validate that the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("GetTokenFromString - Unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing secret key, e.g. []byte("my_secret_key")
		return []byte(SECRET_KEY), nil
	})
}

func (s *service) VerifyJWT(tokenStr string) (bool, *Claims) {
	// validate generated jwt token

	claims := &Claims{}
	token, _ := s.GetTokenFromString(tokenStr, claims)

	if token.Valid {
		token.Claims.Valid()
		if err := claims.Valid(); err != nil {
			return false, claims
		}
	}

	return true, claims
}

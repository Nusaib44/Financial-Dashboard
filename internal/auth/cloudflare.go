package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken parses and validates the trusted Cloudflare token.
// Trusted = we trust the header presence because we are behind Access.
// But we still parse structure.
func ValidateToken(tokenString string) (sub string, email string, err error) {
	if tokenString == "" {
		return "", "", errors.New("missing token")
	}

	// Parse without signature verification (trusted upstream)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	sub, _ = claims["sub"].(string)
	email, _ = claims["email"].(string)

	if sub == "" || email == "" {
		return "", "", errors.New("missing subject or email in token")
	}

	return sub, email, nil
}

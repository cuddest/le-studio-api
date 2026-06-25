package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// TokenIssuer is the expected `iss` claim value.
const TokenIssuer = "le-studio-api"

// Claims defines access-token claims.
type Claims struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	jwtv5.RegisteredClaims
}

// GenerateAccessToken signs a JWT.
func GenerateAccessToken(secret string, subjectID uint, role, email string, ttl time.Duration) (string, error) {
	claims := Claims{
		Role:  role,
		Email: email,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", subjectID),
			Issuer:    TokenIssuer,
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// ParseAccessToken validates and parses JWT.
func ParseAccessToken(secret, token string) (*Claims, error) {
	t, err := jwtv5.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwtv5.Token) (any, error) { return []byte(secret), nil },
		jwtv5.WithValidMethods([]string{jwtv5.SigningMethodHS256.Alg()}),
		jwtv5.WithIssuer(TokenIssuer),
	)
	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return c, nil
}

// HashToken hashes opaque token.
func HashToken(token string) string {
	s := sha256.Sum256([]byte(token))
	return hex.EncodeToString(s[:])
}

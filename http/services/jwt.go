package services

import (
	"errors"
	"time"

	"faizalmaulana/lsp/conf"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(cfg *conf.Config, userID, sessionId, username, role string) (string, error) {
	secret := cfg.JWTSecret
	if secret == "" {
		return "", errors.New("jwt secret is empty: set JWT_SECRET environment variable")
	}

	ttl := cfg.JWTTTL
	if ttl <= 0 {
		ttl = 60
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":        userID,
		"role":       role,
		"session_id": sessionId,
		"iat":        now.Unix(),
		"exp":        now.Add(time.Duration(ttl) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const TokenTTL = 24 * time.Hour

var secret = loadSecret()

func loadSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		// Without a secret no token can be signed, so failing at startup is the only option.
		panic("failed to generate JWT secret: " + err.Error())
	}

	slog.Warn("JWT_SECRET is not set, generated a random secret; tokens become invalid on restart")

	return []byte(hex.EncodeToString(buf))
}

type TokenData struct {
	ID          int64  `json:"id"`
	ClientEmail string `json:"clientID"`
	ExpiresAt   int64  `json:"exp"`
}

func (t TokenData) Valid() error {
	if t.ExpiresAt == 0 {
		return errors.New("token has no expiration")
	}

	if time.Now().Unix() > t.ExpiresAt {
		return errors.New("token expired")
	}

	return nil
}

func ValidateToken(tokenString string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenData{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenData); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GenerateToken(data TokenData) (string, error) {
	if data.ExpiresAt == 0 {
		data.ExpiresAt = time.Now().Add(TokenTTL).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString(secret)
}

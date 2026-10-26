package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenData{},
		func(*jwt.Token) (any, error) { return secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenData)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("token has no expiration")
	}

	return claims, nil
}

func GenerateToken(data TokenData) (string, error) {
	if data.ExpiresAt == nil {
		data.ExpiresAt = jwt.NewNumericDate(time.Now().Add(TokenTTL))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString(secret)
}

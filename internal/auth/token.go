package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt"
)

var secret = []byte("dasdasdadsasdadadssda")

type TokenData struct {
	ID          int64  `json:"id"`
	ClientEmail string `json:"clientID"`
}

func (t TokenData) Valid() error {
	return nil
}

func ValidateToken(tokenString string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenData{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected method")
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString(secret)
}

func ParseToken(data string) *TokenData {
	var tokenJson []byte

	base64.StdEncoding.Decode(tokenJson, []byte(data))

	var tokenData TokenData

	json.Unmarshal(tokenJson, &tokenData)

	return &tokenData
}

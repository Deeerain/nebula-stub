package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ValidateTelegramInitData(botToken string, initDataRaw string, expDuration time.Duration) (url.Values, error) {
	queryParams, err := url.ParseQuery(initDataRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData: %w", err)
	}

	hashHex := queryParams.Get("hash")
	if hashHex == "" {
		return nil, errors.New("hash is missing in initData")
	}

	if expDuration > 0 {
		authDateStr := queryParams.Get("auth_date")
		if authDateStr == "" {
			return nil, errors.New("auth date is missing")
		}

		authTimestamp, err := strconv.ParseInt(authDateStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid auth_date format: %w", err)
		}

		authTime := time.Unix(authTimestamp, 0)
		if time.Since(authTime) > expDuration {
			return nil, errors.New("initData has expired")
		}
	}

	var keys []string
	for k := range queryParams {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	var dataPairs []string
	for _, k := range keys {
		dataPairs = append(dataPairs, fmt.Sprintf("%s=%s", k, queryParams.Get(k)))
	}
	dataCheckString := strings.Join(dataPairs, "\n")

	macSeed := hmac.New(sha256.New, []byte("WebAppData"))
	macSeed.Write([]byte(botToken))
	secretKey := macSeed.Sum(nil)

	macHash := hmac.New(sha256.New, secretKey)
	macHash.Write([]byte(dataCheckString))
	calculateHashBytes := macHash.Sum(nil)
	calculateHashHex := hex.EncodeToString(calculateHashBytes)

	if !hmac.Equal([]byte(calculateHashHex), []byte(hashHex)) {
		return nil, errors.New("invalid signature (hash mismatch)")
	}

	return queryParams, nil
}

package service

import (
	"net/url"
	"os"
	"strconv"
	"testing"
)

func TestGetClientByTgID(t *testing.T) {
	apiToken := os.Getenv("TEST_API_TOKEN")
	rawURL := os.Getenv("TEST_URL")
	tgID := os.Getenv("TEST_TG")
	tgIDint, err := strconv.ParseInt(tgID, 10, 64)

	if err != nil {
		panic(err)
	}

	url, _ := url.Parse(rawURL)
	s := NewXrayService(url, apiToken)

	_, err = s.GetClinetByTgID(tgIDint)
	if err != nil {
		t.Error(err)
	}
}

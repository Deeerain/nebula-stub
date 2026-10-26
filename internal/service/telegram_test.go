package service

import "testing"

func TestNewTelegramServiceRequiresToken(t *testing.T) {
	if _, err := NewTelegramService(""); err == nil {
		t.Fatal("NewTelegramService(\"\") returned nil error")
	}
}

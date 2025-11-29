package model

import (
	"testing"
	"time"
)

func TestAccessKeyIsExpired(t *testing.T) {
	t.Run("never expires when ExpiresAt is nil", func(t *testing.T) {
		accessKey := &AccessKey{
			ExpiresAt: nil,
		}
		if accessKey.IsExpired() {
			t.Error("expected IsExpired to return false when ExpiresAt is nil")
		}
	})

	t.Run("not expired when ExpiresAt is in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		accessKey := &AccessKey{
			ExpiresAt: &future,
		}
		if accessKey.IsExpired() {
			t.Error("expected IsExpired to return false when ExpiresAt is in the future")
		}
	})

	t.Run("expired when ExpiresAt is in the past", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		accessKey := &AccessKey{
			ExpiresAt: &past,
		}
		if !accessKey.IsExpired() {
			t.Error("expected IsExpired to return true when ExpiresAt is in the past")
		}
	})

	t.Run("expired when ExpiresAt is exactly now", func(t *testing.T) {
		now := time.Now().Add(-1 * time.Millisecond)
		accessKey := &AccessKey{
			ExpiresAt: &now,
		}
		if !accessKey.IsExpired() {
			t.Error("expected IsExpired to return true when ExpiresAt is in the past")
		}
	})
}

func TestGenToken(t *testing.T) {
	t.Run("generates token of correct length", func(t *testing.T) {
		lengths := []int{8, 16, 32, 64}
		for _, length := range lengths {
			token, err := GenToken(length)
			if err != nil {
				t.Errorf("GenToken(%d) returned error: %v", length, err)
			}
			if len(token) != length {
				t.Errorf("GenToken(%d) returned token of length %d, expected %d", length, len(token), length)
			}
		}
	})

	t.Run("generates token with valid characters", func(t *testing.T) {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		token, err := GenToken(100)
		if err != nil {
			t.Fatalf("GenToken(100) returned error: %v", err)
		}

		for _, char := range token {
			valid := false
			for _, validChar := range charset {
				if char == validChar {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("GenToken generated invalid character: %c", char)
			}
		}
	})

	t.Run("generates different tokens on subsequent calls", func(t *testing.T) {
		token1, err := GenToken(32)
		if err != nil {
			t.Fatalf("GenToken(32) returned error: %v", err)
		}

		token2, err := GenToken(32)
		if err != nil {
			t.Fatalf("GenToken(32) returned error: %v", err)
		}

		if token1 == token2 {
			t.Error("GenToken generated identical tokens on consecutive calls")
		}
	})

	t.Run("generates empty string for zero length", func(t *testing.T) {
		token, err := GenToken(0)
		if err != nil {
			t.Errorf("GenToken(0) returned error: %v", err)
		}
		if token != "" {
			t.Errorf("GenToken(0) returned non-empty token: %s", token)
		}
	})
}

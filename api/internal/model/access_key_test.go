package model

import (
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestGenAccessKeyTokenBasic(t *testing.T) {
	token, err := GenAccessKeyToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) != 64 {
		t.Fatalf("expected length 64, got %d", len(token))
	}
	allowed := make(map[rune]struct{}, len(charset))
	for _, r := range charset {
		allowed[r] = struct{}{}
	}
	for i, r := range token {
		if _, ok := allowed[r]; !ok {
			t.Fatalf("invalid character at pos %d: %q", i, r)
		}
	}
}

func TestGenAccessKeyTokenUniqueness(t *testing.T) {
	const n = 500
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		token, err := GenAccessKeyToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(token) != 64 {
			t.Fatalf("expected length 64, got %d", len(token))
		}
		if _, exists := seen[token]; exists {
			t.Fatalf("duplicate token found after %d iterations", i)
		}
		seen[token] = struct{}{}
	}
}

func TestGenAccessKeyTokenCharacterSet(t *testing.T) {
	allowed := make(map[rune]struct{}, len(charset))
	for _, r := range charset {
		allowed[r] = struct{}{}
	}
	for i := 0; i < 50; i++ {
		token, err := GenAccessKeyToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, r := range token {
			if _, ok := allowed[r]; !ok {
				t.Fatalf("token contains disallowed character: %q", r)
			}
		}
	}
}

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

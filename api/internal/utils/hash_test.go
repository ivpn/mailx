package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	secret := "mysecret"
	hash, err := Hash(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected a hash, got an empty string")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	if err != nil {
		t.Fatalf("expected hash to match the secret, got %v", err)
	}
}

func TestHashMatches(t *testing.T) {
	secret := "mysecret"
	hash, err := Hash(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !HashMatches(secret, hash) {
		t.Fatalf("expected hash to match the secret")
	}

	if HashMatches("wrongsecret", hash) {
		t.Fatalf("expected hash not to match the wrong secret")
	}
}

func TestHashPassword(t *testing.T) {
	secret := "mysecret"
	hash, err := HashPassword(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected a hash, got an empty string")
	}

	match, err := argon2id.ComparePasswordAndHash(secret, hash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !match {
		t.Fatalf("expected hash to match the secret")
	}
}

func TestHashMatchesPassword(t *testing.T) {
	secret := "mysecret"
	hash, err := HashPassword(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !HashMatchesPassword(secret, hash) {
		t.Fatalf("expected hash to match the secret")
	}

	if HashMatchesPassword("wrongsecret", hash) {
		t.Fatalf("expected hash not to match the wrong secret")
	}
}
func TestHashPGPKey(t *testing.T) {
	key := "mykey"
	expectedHash := sha256.Sum256([]byte(key))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	hash := HashPGPKey(key)
	if hash != expectedHashString {
		t.Fatalf("expected %v, got %v", expectedHashString, hash)
	}

	// Test with an empty key
	emptyKey := ""
	expectedEmptyHash := sha256.Sum256([]byte(emptyKey))
	expectedEmptyHashString := hex.EncodeToString(expectedEmptyHash[:])

	emptyHash := HashPGPKey(emptyKey)
	if emptyHash != expectedEmptyHashString {
		t.Fatalf("expected %v, got %v", expectedEmptyHashString, emptyHash)
	}
}

func TestTimingSafeEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{name: "equal strings", a: "mysecret", b: "mysecret", expected: true},
		{name: "different strings, same length", a: "mysecretA", b: "mysecretB", expected: false},
		{name: "different lengths", a: "short", b: "muchlongersecret", expected: false},
		{name: "empty vs empty", a: "", b: "", expected: true},
		{name: "empty vs non-empty", a: "", b: "mysecret", expected: false},
		{name: "case-sensitive mismatch", a: "MySecret", b: "mysecret", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimingSafeEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("TimingSafeEqual(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

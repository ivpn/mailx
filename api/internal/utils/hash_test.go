package utils

import (
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

package service

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateSlug(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := generateSlug()

	if len(result) != 8 {
		t.Errorf("Expected slug length of 8, but got %d", len(result))
	}

	for _, char := range result {
		if !contains(charset, byte(char)) {
			t.Errorf("Invalid character in slug: %c", char)
		}
	}
}

func contains(s string, char byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == char {
			return true
		}
	}

	return false
}

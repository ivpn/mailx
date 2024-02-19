package model

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateVerification(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789"
	result := GenerateVerification()

	if len(result) != 6 {
		t.Errorf("Expected slug length of 6, but got %d", len(result))
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

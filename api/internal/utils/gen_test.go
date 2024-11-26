package utils

import (
	"strings"
	"testing"
)

func TestRandomStringLength(t *testing.T) {
	length := 10
	result := RandomString(length, AlphaNumericUserFriendly)
	if len(result) != length {
		t.Errorf("Expected string length %d, but got %d", length, len(result))
	}
}

func TestRandomStringCharacters(t *testing.T) {
	length := 10
	result := RandomString(length, AlphaNumericUserFriendly)
	for _, char := range result {
		if !strings.ContainsRune(AlphaNumericUserFriendly, char) {
			t.Errorf("Unexpected character %c in result string", char)
		}
	}
}

func TestRandomStringDifferentResults(t *testing.T) {
	length := 10
	result1 := RandomString(length, AlphaNumericUserFriendly)
	result2 := RandomString(length, AlphaNumericUserFriendly)
	if result1 == result2 {
		t.Errorf("Expected different results, but got the same: %s", result1)
	}
}

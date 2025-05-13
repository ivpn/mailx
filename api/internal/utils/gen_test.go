package utils

import (
	"strings"
	"testing"
)

func TestRandomStringLength(t *testing.T) {
	length := 10
	result, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	if len(result) != length {
		t.Errorf("Expected string length %d, but got %d", length, len(result))
	}
}

func TestRandomStringCharacters(t *testing.T) {
	length := 10
	result, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	for _, char := range result {
		if !strings.ContainsRune(AlphaNumericUserFriendly, char) {
			t.Errorf("Unexpected character %c in result string", char)
		}
	}
}

func TestRandomStringDifferentResults(t *testing.T) {
	length := 10
	result1, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	result2, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	if result1 == result2 {
		t.Errorf("Expected different results, but got the same: %s", result1)
	}
	result3, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	result4, err := RandomString(length, AlphaNumericUserFriendly)
	if err != nil {
		t.Fatalf("RandomString returned an error: %v", err)
	}
	if result3 == result4 {
		t.Errorf("Expected different results, but got the same: %s", result3)
	}
}

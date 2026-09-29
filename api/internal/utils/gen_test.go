package utils

import (
	"encoding/base32"
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

func TestRandomBytesBase32DecodesToExactLength(t *testing.T) {
	n := 16
	result, err := RandomBytesBase32(n)
	if err != nil {
		t.Fatalf("RandomBytesBase32 returned an error: %v", err)
	}

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(result)
	if err != nil {
		t.Fatalf("failed to decode result: %v", err)
	}
	if len(decoded) != n {
		t.Errorf("Expected %d decoded bytes, but got %d", n, len(decoded))
	}
}

func TestRandomBytesBase32NoPadding(t *testing.T) {
	result, err := RandomBytesBase32(16)
	if err != nil {
		t.Fatalf("RandomBytesBase32 returned an error: %v", err)
	}
	if strings.Contains(result, "=") {
		t.Errorf("Expected no padding characters, but got: %s", result)
	}
}

func TestRandomBytesBase32DifferentResults(t *testing.T) {
	result1, err := RandomBytesBase32(16)
	if err != nil {
		t.Fatalf("RandomBytesBase32 returned an error: %v", err)
	}
	result2, err := RandomBytesBase32(16)
	if err != nil {
		t.Fatalf("RandomBytesBase32 returned an error: %v", err)
	}
	if result1 == result2 {
		t.Errorf("Expected different results, but got the same: %s", result1)
	}
}

func TestRandomBytesBase32InvalidLength(t *testing.T) {
	_, err := RandomBytesBase32(0)
	if err != ErrInvalidLength {
		t.Errorf("Expected ErrInvalidLength, but got: %v", err)
	}
}

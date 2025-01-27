package model

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAlias(t *testing.T) {
	t.Run("RandomChars", func(t *testing.T) {
		alias := GenerateAlias(AliasFormatRandomChars)
		if len(alias) != 8 {
			t.Errorf("Expected alias length to be 8, got %d", len(alias))
		}
		if !isAlphanumeric(alias) {
			t.Errorf("Expected alias to be alphanumeric, got %s", alias)
		}
	})

	t.Run("UUID", func(t *testing.T) {
		alias := GenerateAlias(AliasFormatUUID)
		if _, err := uuid.Parse(alias); err != nil {
			t.Errorf("Expected alias to be a valid UUID, got %s", alias)
		}
	})

	t.Run("RandomWords", func(t *testing.T) {
		alias := GenerateAlias("words")
		parts := strings.Split(alias, ".")
		if len(parts) != 2 {
			t.Errorf("Expected alias to have two parts separated by a dot, got %s", alias)
		}
		if !isAlphanumeric(parts[1]) {
			t.Errorf("Expected second part of alias to be alphanumeric, got %s", parts[1])
		}
	})
}

func TestGenerateCatchAllAlias(t *testing.T) {
	t.Run("SimpleSuffix", func(t *testing.T) {
		suffix := "test"
		result := GenerateCatchAllAlias(suffix)
		expected := "*+test"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("EmptySuffix", func(t *testing.T) {
		suffix := ""
		result := GenerateCatchAllAlias(suffix)
		expected := "*+"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("SpecialCharactersSuffix", func(t *testing.T) {
		suffix := "test@123"
		result := GenerateCatchAllAlias(suffix)
		expected := "*+test@123"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

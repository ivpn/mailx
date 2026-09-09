package model

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAlias(t *testing.T) {
	tests := []struct {
		format string
		suffix string
	}{
		{AliasFormatRandomChars, ""},
		{AliasFormatUUID, ""},
		{AliasFormatWildcard, "test"},
		{AliasFormatRandomWords, ""},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			alias := GenerateAlias(tt.format, tt.suffix)
			switch tt.format {
			case AliasFormatRandomChars:
				if len(alias) != 8 || !isAlphanumeric(alias) {
					t.Errorf("expected 8 alphanumeric characters, got %s", alias)
				}
			case AliasFormatUUID:
				if _, err := uuid.Parse(alias); err != nil {
					t.Errorf("expected valid UUID, got %s", alias)
				}
			case AliasFormatWildcard:
				expected := "*+" + tt.suffix
				if alias != expected {
					t.Errorf("expected %s, got %s", expected, alias)
				}
			case AliasFormatRandomWords:
				parts := strings.Split(alias, ".")
				if len(parts) != 2 || !isAlphanumeric(parts[1]) {
					t.Errorf("expected format adjective.noun, got %s", alias)
				}
			}
		})
	}
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func TestGenerateWildcardAlias(t *testing.T) {
	tests := []struct {
		name      string
		localPart string
		delimiter string
		expected  string
	}{
		{name: "plus delimiter", localPart: "news", delimiter: WildcardDelimiterPlus, expected: "*+news"},
		{name: "dot delimiter", localPart: "news", delimiter: WildcardDelimiterDot, expected: "*.news"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateWildcardAlias(tt.localPart, tt.delimiter)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestWildcardAliasDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected string
	}{
		{name: "plus wildcard", alias: "*+news@domain.com", expected: "+"},
		{name: "dot wildcard", alias: "*.news@domain.com", expected: "."},
		{name: "not a wildcard", alias: "news@domain.com", expected: ""},
		{name: "empty string", alias: "", expected: ""},
		{name: "just the asterisk", alias: "*", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WildcardAliasDelimiter(tt.alias)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestIsValidWildcardDelimiter(t *testing.T) {
	tests := []struct {
		delimiter string
		expected  bool
	}{
		{delimiter: "+", expected: true},
		{delimiter: ".", expected: true},
		{delimiter: "-", expected: false},
		{delimiter: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.delimiter, func(t *testing.T) {
			got := IsValidWildcardDelimiter(tt.delimiter)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

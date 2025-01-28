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
		{AliasFormatCatchAll, "test"},
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
			case AliasFormatCatchAll:
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

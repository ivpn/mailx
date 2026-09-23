package model

import (
	"slices"
	"testing"
)

func TestNormalizeBlocklistValue(t *testing.T) {
	tests := []struct {
		name     string
		itemType BlocklistType
		value    string
		want     string
	}{
		{"email lowercased and trimmed", BlocklistEmail, "  Spammer@Example.COM ", "spammer@example.com"},
		{"domain lowercased", BlocklistDomain, "Example.COM", "example.com"},
		{"domain leading @ stripped", BlocklistDomain, "@example.com", "example.com"},
		{"email keeps its @", BlocklistEmail, "@example.com", "@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeBlocklistValue(tt.itemType, tt.value); got != tt.want {
				t.Errorf("NormalizeBlocklistValue(%q, %q) = %q, want %q", tt.itemType, tt.value, got, tt.want)
			}
		})
	}
}

func TestBlocklistCandidates(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		wantEmail   string
		wantDomains []string
	}{
		{"plain domain", "Sender@Example.com", "sender@example.com", []string{"example.com"}},
		{"subdomain includes parents above TLD", "a@mail.news.example.co.uk", "a@mail.news.example.co.uk", []string{"mail.news.example.co.uk", "news.example.co.uk", "example.co.uk", "co.uk"}},
		{"no domain part", "sender", "sender", nil},
		{"empty domain part", "sender@", "sender@", nil},
		{"empty", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, domains := BlocklistCandidates(tt.from)
			if email != tt.wantEmail {
				t.Errorf("email = %q, want %q", email, tt.wantEmail)
			}
			if !slices.Equal(domains, tt.wantDomains) {
				t.Errorf("domains = %v, want %v", domains, tt.wantDomains)
			}
		})
	}
}

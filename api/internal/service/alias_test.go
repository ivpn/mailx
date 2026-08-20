package service

import (
	"testing"

	"ivpn.net/email/api/internal/model"
)

func TestAliasDomainPart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "standard email", input: "user@example.com", expected: "example.com"},
		{name: "multiple at signs", input: "user@foo@example.com", expected: "foo@example.com"},
		{name: "no at sign", input: "userexample.com", expected: ""},
		{name: "empty string", input: "", expected: ""},
		{name: "only at sign", input: "@", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aliasDomainPart(tt.input)
			if got != tt.expected {
				t.Errorf("aliasDomainPart(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAliasLocalPart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "standard email", input: "user@example.com", expected: "user"},
		{name: "multiple at signs", input: "user@foo@example.com", expected: "user"},
		{name: "no at sign", input: "userexample.com", expected: ""},
		{name: "empty string", input: "", expected: ""},
		{name: "only at sign", input: "@", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aliasLocalPart(tt.input)
			if got != tt.expected {
				t.Errorf("aliasLocalPart(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsCustomAliasDomain(t *testing.T) {
	tests := []struct {
		name              string
		domainPart        string
		predefinedDomains string
		expected          bool
	}{
		{name: "domain in predefined list", domainPart: "example.com", predefinedDomains: "example.com,other.com", expected: false},
		{name: "domain not in predefined list", domainPart: "custom.com", predefinedDomains: "example.com,other.com", expected: true},
		{name: "empty predefined domains", domainPart: "example.com", predefinedDomains: "", expected: true},
		{name: "empty domain part", domainPart: "", predefinedDomains: "example.com", expected: false},
		{name: "single match", domainPart: "other.com", predefinedDomains: "other.com", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCustomAliasDomain(tt.domainPart, tt.predefinedDomains)
			if got != tt.expected {
				t.Errorf("isCustomAliasDomain(%q, %q) = %v, want %v", tt.domainPart, tt.predefinedDomains, got, tt.expected)
			}
		})
	}
}

func TestIsCustomDomainEnabled(t *testing.T) {
	tests := []struct {
		name            string
		domainPart      string
		verifiedDomains []model.Domain
		expected        bool
	}{
		{
			name:       "domain found and enabled",
			domainPart: "example.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", Enabled: true},
				{Name: "other.com", Enabled: false},
			},
			expected: true,
		},
		{
			name:       "domain found and disabled",
			domainPart: "example.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", Enabled: false},
			},
			expected: false,
		},
		{
			name:       "domain not in list",
			domainPart: "missing.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", Enabled: true},
			},
			expected: false,
		},
		{name: "empty domain list", domainPart: "example.com", verifiedDomains: []model.Domain{}, expected: false},
		{name: "nil domain list", domainPart: "example.com", verifiedDomains: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCustomDomainEnabled(tt.domainPart, tt.verifiedDomains)
			if got != tt.expected {
				t.Errorf("isCustomDomainEnabled(%q, ...) = %v, want %v", tt.domainPart, got, tt.expected)
			}
		})
	}
}

func TestIsCreateAliasEnabled(t *testing.T) {
	tests := []struct {
		name            string
		domainPart      string
		verifiedDomains []model.Domain
		expected        bool
	}{
		{
			name:       "domain found and create alias enabled",
			domainPart: "example.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", CreateAlias: true},
				{Name: "other.com", CreateAlias: false},
			},
			expected: true,
		},
		{
			name:       "domain found and create alias disabled",
			domainPart: "example.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", CreateAlias: false},
			},
			expected: false,
		},
		{
			name:       "domain not in list",
			domainPart: "missing.com",
			verifiedDomains: []model.Domain{
				{Name: "example.com", CreateAlias: true},
			},
			expected: false,
		},
		{name: "empty domain list", domainPart: "example.com", verifiedDomains: []model.Domain{}, expected: false},
		{name: "nil domain list", domainPart: "example.com", verifiedDomains: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCreateAliasEnabled(tt.domainPart, tt.verifiedDomains)
			if got != tt.expected {
				t.Errorf("isCreateAliasEnabled(%q, ...) = %v, want %v", tt.domainPart, got, tt.expected)
			}
		})
	}
}

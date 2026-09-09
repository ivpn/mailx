package service

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestPostAlias_CustomDomainClassifiesStoreErrors(t *testing.T) {
	tests := []struct {
		name        string
		storeErr    error
		expectedErr error
	}{
		{name: "hourly inbound limit reached", storeErr: model.ErrInboundHourlyLimit, expectedErr: ErrPostInboundAlias},
		{name: "daily alias limit reached", storeErr: model.ErrDailyAliasLimit, expectedErr: ErrPostAliasLimit},
		{name: "generic store error", storeErr: errNotFound, expectedErr: ErrPostAlias},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newFakeStore()
			store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
			store.postAliasErr = tt.storeErr
			s := newTestService(store)

			_, err := s.PostAlias(context.Background(), model.Alias{UserID: "user-1", Origin: model.Inbound}, model.AliasFormatCustom, "customdomain.com", "newalias", "")
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestPostAlias_CustomDomainSucceeds(t *testing.T) {
	store := newFakeStore()
	store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
	s := newTestService(store)

	alias, err := s.PostAlias(context.Background(), model.Alias{UserID: "user-1", Origin: model.Inbound}, model.AliasFormatCustom, "customdomain.com", "newalias", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "newalias@customdomain.com" {
		t.Errorf("expected alias name newalias@customdomain.com, got %s", alias.Name)
	}
}

func TestPostAlias_WildcardSucceedsWithDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		expected  string
	}{
		{name: "plus delimiter", delimiter: model.WildcardDelimiterPlus, expected: "*+news@customdomain.com"},
		{name: "dot delimiter", delimiter: model.WildcardDelimiterDot, expected: "*.news@customdomain.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newFakeStore()
			store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
			s := newTestService(store)

			alias, err := s.PostAlias(context.Background(), model.Alias{UserID: "user-1"}, model.AliasFormatWildcard, "customdomain.com", "news", tt.delimiter)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if alias.Name != tt.expected {
				t.Errorf("expected alias name %s, got %s", tt.expected, alias.Name)
			}
			if !alias.Wildcard {
				t.Error("expected alias.Wildcard to be true")
			}
		})
	}
}

func TestPostAlias_WildcardDomainCapReachedRegardlessOfDelimiterMix(t *testing.T) {
	store := newFakeStore()
	store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
	store.aliases["*+news@customdomain.com"] = model.Alias{Name: "*+news@customdomain.com", UserID: "user-1", Wildcard: true}
	store.aliases["*.deals@customdomain.com"] = model.Alias{Name: "*.deals@customdomain.com", UserID: "user-1", Wildcard: true}
	s := newTestService(store)

	_, err := s.PostAlias(context.Background(), model.Alias{UserID: "user-1"}, model.AliasFormatWildcard, "customdomain.com", "third", model.WildcardDelimiterDot)
	if !errors.Is(err, model.ErrDuplicateAliasDomain) {
		t.Errorf("expected ErrDuplicateAliasDomain, got %v", err)
	}
}

func TestPostAlias_WildcardInvalidDelimiterRejected(t *testing.T) {
	store := newFakeStore()
	store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
	s := newTestService(store)

	_, err := s.PostAlias(context.Background(), model.Alias{UserID: "user-1"}, model.AliasFormatWildcard, "customdomain.com", "news", "-")
	if !errors.Is(err, ErrPostAlias) {
		t.Errorf("expected ErrPostAlias for an invalid delimiter, got %v", err)
	}
}

func TestGetWildcardDomainInfo(t *testing.T) {
	store := newFakeStore()
	store.aliases["*+news@customdomain.com"] = model.Alias{Name: "*+news@customdomain.com", UserID: "user-1", Wildcard: true}
	store.aliases["*.deals@customdomain.com"] = model.Alias{Name: "*.deals@customdomain.com", UserID: "user-1", Wildcard: true}
	store.aliases["*+other@otherdomain.com"] = model.Alias{Name: "*+other@otherdomain.com", UserID: "user-1", Wildcard: true}
	store.aliases["*+news2@customdomain.com"] = model.Alias{Name: "*+news2@customdomain.com", UserID: "user-2", Wildcard: true}
	s := newTestService(store)

	info, err := s.GetWildcardDomainInfo(context.Background(), "user-1", "customdomain.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info.Count != 2 {
		t.Errorf("expected count 2, got %d", info.Count)
	}
	if info.Limit != model.MaxWildcardAliasesPerDomain {
		t.Errorf("expected limit %d, got %d", model.MaxWildcardAliasesPerDomain, info.Limit)
	}
	if len(info.DelimitersUsed) != 2 {
		t.Fatalf("expected 2 delimiters used, got %+v", info.DelimitersUsed)
	}
	found := map[string]bool{}
	for _, d := range info.DelimitersUsed {
		found[d] = true
	}
	if !found["+"] || !found["."] {
		t.Errorf("expected delimiters used to include both + and ., got %+v", info.DelimitersUsed)
	}
}

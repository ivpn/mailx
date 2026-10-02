package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

func TestForgetAlias(t *testing.T) {
	t.Run("rejects alias on a built-in domain", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["random@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"},
			Name:      "random@mailx.net",
			UserID:    "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		svc := newTestService(store)

		err := svc.ForgetAlias(context.Background(), "alias-1", "user-1")
		if !errors.Is(err, ErrForgetAliasNotCustomDomain) {
			t.Errorf("expected ErrForgetAliasNotCustomDomain, got %v", err)
		}
	})

	t.Run("permanently removes an active custom-domain alias", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["custom@example.com"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"},
			Name:      "custom@example.com",
			UserID:    "user-1",
		}
		svc := newTestService(store)

		err := svc.ForgetAlias(context.Background(), "alias-2", "user-1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if _, ok := store.aliases["custom@example.com"]; ok {
			t.Errorf("expected alias to be removed from store")
		}
	})

	t.Run("permanently removes a deleted custom-domain alias", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["custom@example.com"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-3"},
			Name:      "custom@example.com",
			UserID:    "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		svc := newTestService(store)

		err := svc.ForgetAlias(context.Background(), "alias-3", "user-1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if _, ok := store.aliases["custom@example.com"]; ok {
			t.Errorf("expected alias to be removed from store")
		}
	})

	t.Run("rejects alias belonging to a different user", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["custom@example.com"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-4"},
			Name:      "custom@example.com",
			UserID:    "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		svc := newTestService(store)

		err := svc.ForgetAlias(context.Background(), "alias-4", "user-2")
		if !errors.Is(err, ErrGetAlias) {
			t.Errorf("expected ErrGetAlias, got %v", err)
		}
	})
}

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

func TestImportAliases_InactiveSubscriptionFailsFast(t *testing.T) {
	store := newFakeStore()
	store.subscription = model.Subscription{ActiveUntil: time.Now().Add(-time.Hour)}
	s := newTestService(store)

	reqs := []model.AliasImportReq{
		{Domain: "customdomain.com", LocalPart: "newalias", Recipients: "rcpt@example.com", Format: model.AliasFormatCustom},
	}

	aliases, err := s.ImportAliases(context.Background(), reqs, "user-1")
	if !errors.Is(err, ErrPostAliasInactiveSub) {
		t.Errorf("expected error %v, got %v", ErrPostAliasInactiveSub, err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected no aliases imported, got %d", len(aliases))
	}
}

func TestImportAliases_ActiveSubscriptionImportsValidRows(t *testing.T) {
	store := newFakeStore()
	store.subscription = model.Subscription{ActiveUntil: time.Now().Add(time.Hour)}
	store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-1", Enabled: true}
	store.verifiedRecipients["user-1"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	reqs := []model.AliasImportReq{
		{Domain: "customdomain.com", LocalPart: "newalias", Recipients: "rcpt@example.com", Format: model.AliasFormatCustom},
	}

	aliases, err := s.ImportAliases(context.Background(), reqs, "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(aliases) != 1 || aliases[0].Name != "newalias@customdomain.com" {
		t.Errorf("expected 1 imported alias newalias@customdomain.com, got %+v", aliases)
	}
}

func TestDedupeAliasIDs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{name: "no duplicates", input: []string{"a", "b", "c"}, expected: []string{"a", "b", "c"}},
		{name: "duplicates removed, order preserved", input: []string{"a", "b", "a", "c", "b"}, expected: []string{"a", "b", "c"}},
		{name: "empty input", input: []string{}, expected: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedupeAliasIDs(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("dedupeAliasIDs(%v) = %v, want %v", tt.input, got, tt.expected)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("dedupeAliasIDs(%v) = %v, want %v", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestBulkUpdateAliasEnabled(t *testing.T) {
	t.Run("skips aliases left without a recipient", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			Recipients: "user@example.com",
		}
		store.aliases["two@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1",
		}
		s := newTestService(store)

		count, err := s.BulkUpdateAliasEnabled(context.Background(), []string{"alias-1", "alias-2"}, "user-1", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 alias updated, got %d", count)
		}
		if !store.aliases["one@mailx.net"].Enabled {
			t.Errorf("expected alias-1 to be activated")
		}
		if store.aliases["two@mailx.net"].Enabled {
			t.Errorf("expected alias-2 without a recipient to stay disabled")
		}
	})

	t.Run("skips deleted aliases", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			Recipients: "user@example.com",
		}
		store.aliases["two@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1",
			Recipients: "user@example.com",
			DeletedAt:  gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		count, err := s.BulkUpdateAliasEnabled(context.Background(), []string{"alias-1", "alias-2"}, "user-1", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 alias updated, got %d", count)
		}
		if store.aliases["two@mailx.net"].Enabled {
			t.Errorf("expected deleted alias-2 to stay disabled")
		}
	})

	t.Run("rejects the request when no alias is eligible", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			Recipients: "user@example.com",
			DeletedAt:  gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		count, err := s.BulkUpdateAliasEnabled(context.Background(), []string{"alias-1"}, "user-1", true)
		if !errors.Is(err, ErrBulkUpdateAliasNotEligible) {
			t.Errorf("expected ErrBulkUpdateAliasNotEligible, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 aliases updated, got %d", count)
		}
	})

	t.Run("ignores aliases owned by another user", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-2",
			Recipients: "user@example.com",
		}
		s := newTestService(store)

		_, err := s.BulkUpdateAliasEnabled(context.Background(), []string{"alias-1"}, "user-1", true)
		if !errors.Is(err, ErrBulkUpdateAliasNotEligible) {
			t.Errorf("expected ErrBulkUpdateAliasNotEligible, got %v", err)
		}
		if store.aliases["one@mailx.net"].Enabled {
			t.Errorf("expected another user's alias to stay untouched")
		}
	})
}

func TestBulkUpdateAliasPinned(t *testing.T) {
	t.Run("skips deleted aliases", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1"}
		store.aliases["two@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		count, err := s.BulkUpdateAliasPinned(context.Background(), []string{"alias-1", "alias-2"}, "user-1", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 alias updated, got %d", count)
		}
		if !store.aliases["one@mailx.net"].Pinned {
			t.Errorf("expected alias-1 to be pinned")
		}
		if store.aliases["two@mailx.net"].Pinned {
			t.Errorf("expected deleted alias-2 to stay unpinned")
		}
	})

	t.Run("rejects the request when every alias is deleted", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		count, err := s.BulkUpdateAliasPinned(context.Background(), []string{"alias-1"}, "user-1", true)
		if !errors.Is(err, ErrBulkUpdateAliasNotEligible) {
			t.Errorf("expected ErrBulkUpdateAliasNotEligible, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 aliases updated, got %d", count)
		}
	})
}

func TestBulkDeleteAlias(t *testing.T) {
	t.Run("soft-deletes all non-deleted owned aliases", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1"}
		store.aliases["two@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1"}
		s := newTestService(store)

		err := s.BulkDeleteAlias(context.Background(), []string{"alias-1", "alias-2", "alias-1"}, "user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !store.aliases["one@mailx.net"].DeletedAt.Valid || !store.aliases["two@mailx.net"].DeletedAt.Valid {
			t.Errorf("expected both aliases to be soft-deleted")
		}
	})

	t.Run("rejects the whole batch when one alias is already deleted", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1"}
		store.aliases["two@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		err := s.BulkDeleteAlias(context.Background(), []string{"alias-1", "alias-2"}, "user-1")
		if !errors.Is(err, ErrBulkDeleteAliasNotEligible) {
			t.Errorf("expected ErrBulkDeleteAliasNotEligible, got %v", err)
		}
		if store.aliases["one@mailx.net"].DeletedAt.Valid {
			t.Errorf("expected alias-1 to remain untouched after a rejected batch")
		}
	})
}

func TestBulkRestoreAlias(t *testing.T) {
	t.Run("restores all deleted owned aliases", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		err := s.BulkRestoreAlias(context.Background(), []string{"alias-1"}, "user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if store.aliases["one@mailx.net"].DeletedAt.Valid {
			t.Errorf("expected alias to be restored")
		}
	})

	t.Run("rejects the whole batch when one alias is not deleted", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@mailx.net", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		store.aliases["two@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1"}
		s := newTestService(store)

		err := s.BulkRestoreAlias(context.Background(), []string{"alias-1", "alias-2"}, "user-1")
		if !errors.Is(err, ErrBulkRestoreAliasNotEligible) {
			t.Errorf("expected ErrBulkRestoreAliasNotEligible, got %v", err)
		}
		if !store.aliases["one@mailx.net"].DeletedAt.Valid {
			t.Errorf("expected alias-1 to remain untouched after a rejected batch")
		}
	})
}

func TestBulkForgetAlias(t *testing.T) {
	t.Run("permanently removes only custom-domain aliases", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@customdomain.com"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@customdomain.com", UserID: "user-1"}
		store.aliases["two@customdomain.com"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@customdomain.com", UserID: "user-1",
			DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		}
		s := newTestService(store)

		err := s.BulkForgetAlias(context.Background(), []string{"alias-1", "alias-2"}, "user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(store.aliases) != 0 {
			t.Errorf("expected both aliases to be permanently removed, got %+v", store.aliases)
		}
	})

	t.Run("rejects the whole batch when one alias is on a built-in domain", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@customdomain.com"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@customdomain.com", UserID: "user-1"}
		store.aliases["two@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@mailx.net", UserID: "user-1"}
		s := newTestService(store)

		err := s.BulkForgetAlias(context.Background(), []string{"alias-1", "alias-2"}, "user-1")
		if !errors.Is(err, ErrForgetAliasNotCustomDomain) {
			t.Errorf("expected ErrForgetAliasNotCustomDomain, got %v", err)
		}
		if len(store.aliases) != 2 {
			t.Errorf("expected no aliases to be removed after a rejected batch, got %+v", store.aliases)
		}
	})

	t.Run("rejects the whole batch when one alias doesn't belong to the user", func(t *testing.T) {
		store := newFakeStore()
		store.aliases["one@customdomain.com"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "one@customdomain.com", UserID: "user-1"}
		store.aliases["two@customdomain.com"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-2"}, Name: "two@customdomain.com", UserID: "user-2"}
		s := newTestService(store)

		err := s.BulkForgetAlias(context.Background(), []string{"alias-1", "alias-2"}, "user-1")
		if !errors.Is(err, ErrBulkForgetAlias) {
			t.Errorf("expected ErrBulkForgetAlias, got %v", err)
		}
		if len(store.aliases) != 2 {
			t.Errorf("expected no aliases to be removed after a rejected batch, got %+v", store.aliases)
		}
	})
}

func gracePeriodSubscription() model.Subscription {
	return model.Subscription{UpdatedAt: time.Now().Add(-60 * time.Hour), ActiveUntil: time.Now().Add(-time.Hour)}
}

func pendingDeleteSubscription() model.Subscription {
	return model.Subscription{ActiveUntil: time.Now().Add(30 * 24 * time.Hour), Terminated: true}
}

func newPauseTestService(t *testing.T) (*Service, *fakeStore) {
	t.Helper()

	store := newFakeStore()
	store.subscription = activeSubscription()
	store.aliases["promo@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "promo@mailx.net", UserID: "user-1", Enabled: true}
	store.aliases["gone@mailx.net"] = model.Alias{
		BaseModel: model.BaseModel{ID: "alias-deleted"}, Name: "gone@mailx.net", UserID: "user-1",
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}
	store.aliases["theirs@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-other"}, Name: "theirs@mailx.net", UserID: "user-2", Enabled: true}

	return newTestService(store), store
}

func TestPauseAlias(t *testing.T) {
	tests := []struct {
		value int
		unit  string
		want  time.Duration
	}{
		{value: 1, unit: "hours", want: time.Hour},
		{value: 10, unit: "hours", want: 10 * time.Hour},
		{value: 1, unit: "days", want: 24 * time.Hour},
		{value: 10, unit: "days", want: 10 * 24 * time.Hour},
		{value: 1, unit: "weeks", want: 7 * 24 * time.Hour},
		{value: 10, unit: "weeks", want: 70 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d %s", tt.value, tt.unit), func(t *testing.T) {
			s, store := newPauseTestService(t)

			before := time.Now().UTC().Truncate(time.Second)
			got, err := s.PauseAlias(context.Background(), "alias-1", "user-1", tt.value, tt.unit)
			after := time.Now().UTC()
			if err != nil {
				t.Fatalf("PauseAlias() error = %v", err)
			}

			if got.Location() != time.UTC {
				t.Errorf("paused_until %v is not UTC", got)
			}
			if got.Before(before.Add(tt.want)) || got.After(after.Add(tt.want)) {
				t.Errorf("paused_until = %v, want now + %v", got, tt.want)
			}
			if stored := store.aliases["promo@mailx.net"].PausedUntil; stored == nil || !stored.Equal(got) {
				t.Errorf("stored paused_until = %v, want %v", stored, got)
			}
		})
	}

	t.Run("re-pausing replaces the end time", func(t *testing.T) {
		s, store := newPauseTestService(t)

		if _, err := s.PauseAlias(context.Background(), "alias-1", "user-1", 10, "weeks"); err != nil {
			t.Fatalf("PauseAlias() error = %v", err)
		}
		got, err := s.PauseAlias(context.Background(), "alias-1", "user-1", 1, "hours")
		if err != nil {
			t.Fatalf("PauseAlias() error = %v", err)
		}
		if got.After(time.Now().Add(time.Hour)) {
			t.Errorf("paused_until = %v, want the new 1 hour pause, not one added to the old pause", got)
		}
		if stored := store.aliases["promo@mailx.net"].PausedUntil; !stored.Equal(got) {
			t.Errorf("stored paused_until = %v, want %v", stored, got)
		}
	})

	t.Run("invalid durations are rejected", func(t *testing.T) {
		s, store := newPauseTestService(t)

		for _, tt := range []struct {
			value int
			unit  string
		}{{0, "days"}, {11, "days"}, {1, "months"}} {
			if _, err := s.PauseAlias(context.Background(), "alias-1", "user-1", tt.value, tt.unit); !errors.Is(err, ErrPauseAliasDuration) {
				t.Errorf("PauseAlias(%d, %q) error = %v, want ErrPauseAliasDuration", tt.value, tt.unit, err)
			}
		}
		if store.aliases["promo@mailx.net"].PausedUntil != nil {
			t.Error("an invalid pause was stored")
		}
	})

	for name, id := range map[string]string{
		"another user's alias": "alias-other",
		"a deleted alias":      "alias-deleted",
		"an unknown ID":        "alias-missing",
	} {
		t.Run(name+" can't be paused", func(t *testing.T) {
			s, store := newPauseTestService(t)

			if _, err := s.PauseAlias(context.Background(), id, "user-1", 1, "days"); !errors.Is(err, ErrUpdateAlias) {
				t.Errorf("PauseAlias() error = %v, want ErrUpdateAlias", err)
			}
			for _, a := range store.aliases {
				if a.PausedUntil != nil {
					t.Errorf("alias %s was paused", a.ID)
				}
			}
		})
	}
}

func TestUnpauseAlias(t *testing.T) {
	t.Run("clears a running pause", func(t *testing.T) {
		s, _ := newPauseTestService(t)
		if _, err := s.PauseAlias(context.Background(), "alias-1", "user-1", 1, "days"); err != nil {
			t.Fatalf("PauseAlias() error = %v", err)
		}

		if err := s.UnpauseAlias(context.Background(), "alias-1", "user-1"); err != nil {
			t.Fatalf("UnpauseAlias() error = %v", err)
		}
		alias, err := s.GetAlias(context.Background(), "alias-1", "user-1")
		if err != nil || alias.PausedUntil != nil {
			t.Errorf("alias = %+v, %v, want paused_until nil", alias, err)
		}
	})

	t.Run("succeeds for an alias that isn't paused", func(t *testing.T) {
		s, _ := newPauseTestService(t)

		if err := s.UnpauseAlias(context.Background(), "alias-1", "user-1"); err != nil {
			t.Errorf("UnpauseAlias() error = %v", err)
		}
	})

	for name, id := range map[string]string{
		"another user's alias": "alias-other",
		"a deleted alias":      "alias-deleted",
		"an unknown ID":        "alias-missing",
	} {
		t.Run(name+" can't be unpaused", func(t *testing.T) {
			s, _ := newPauseTestService(t)

			if err := s.UnpauseAlias(context.Background(), id, "user-1"); !errors.Is(err, ErrUpdateAlias) {
				t.Errorf("UnpauseAlias() error = %v, want ErrUpdateAlias", err)
			}
		})
	}
}

func TestPauseAlias_Subscription(t *testing.T) {
	inactive := map[string]model.Subscription{
		"limited_access": {},
		"pending_delete": pendingDeleteSubscription(),
	}
	for name, sub := range inactive {
		t.Run(name+" rejects both and changes nothing", func(t *testing.T) {
			s, store := newPauseTestService(t)
			pausedUntil := time.Now().Add(time.Hour).UTC()
			alias := store.aliases["promo@mailx.net"]
			alias.PausedUntil = &pausedUntil
			store.aliases["promo@mailx.net"] = alias
			store.subscription = sub

			if _, err := s.PauseAlias(context.Background(), "alias-1", "user-1", 10, "weeks"); !errors.Is(err, ErrPauseAliasInactiveSub) {
				t.Errorf("PauseAlias() error = %v, want ErrPauseAliasInactiveSub", err)
			}
			if err := s.UnpauseAlias(context.Background(), "alias-1", "user-1"); !errors.Is(err, ErrUnpauseAliasInactiveSub) {
				t.Errorf("UnpauseAlias() error = %v, want ErrUnpauseAliasInactiveSub", err)
			}
			if stored := store.aliases["promo@mailx.net"].PausedUntil; stored == nil || !stored.Equal(pausedUntil) {
				t.Errorf("paused_until = %v, want the running pause kept", stored)
			}
		})
	}

	active := map[string]model.Subscription{
		"active":       activeSubscription(),
		"grace_period": gracePeriodSubscription(),
	}
	for name, sub := range active {
		t.Run(name+" allows both", func(t *testing.T) {
			s, store := newPauseTestService(t)
			store.subscription = sub

			if _, err := s.PauseAlias(context.Background(), "alias-1", "user-1", 1, "hours"); err != nil {
				t.Errorf("PauseAlias() error = %v", err)
			}
			if err := s.UnpauseAlias(context.Background(), "alias-1", "user-1"); err != nil {
				t.Errorf("UnpauseAlias() error = %v", err)
			}
		})
	}
}

func TestGetAliases_ClearsExpiredPause(t *testing.T) {
	s, store := newPauseTestService(t)
	past := time.Now().Add(-time.Minute).UTC()
	future := time.Now().Add(time.Hour).UTC()

	expired := store.aliases["promo@mailx.net"]
	expired.PausedUntil = &past
	store.aliases["promo@mailx.net"] = expired
	store.aliases["paused@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-paused"}, Name: "paused@mailx.net", UserID: "user-1", PausedUntil: &future}

	list, err := s.GetAliases(context.Background(), "user-1", 25, 1, "created_at", "DESC", "", "", "")
	if err != nil {
		t.Fatalf("GetAliases() error = %v", err)
	}
	for _, a := range list.Aliases {
		switch a.ID {
		case "alias-1":
			if a.PausedUntil != nil {
				t.Errorf("ended pause returned as %v, want nil", a.PausedUntil)
			}
		case "alias-paused":
			if a.PausedUntil == nil || !a.PausedUntil.Equal(future) {
				t.Errorf("running pause returned as %v, want %v", a.PausedUntil, future)
			}
		}
	}

	alias, err := s.GetAlias(context.Background(), "alias-1", "user-1")
	if err != nil || alias.PausedUntil != nil {
		t.Errorf("GetAlias() = %v, %v, want paused_until nil", alias.PausedUntil, err)
	}
}

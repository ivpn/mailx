package service

import (
	"context"
	"errors"
	"slices"
	"testing"

	"ivpn.net/email/api/internal/model"
)

// blocklistStore extends fakeStore with the blocklist and diagnostic log methods the blocklist
// tests exercise.
type blocklistStore struct {
	*fakeStore

	items         []model.BlocklistItem
	matchErr      error
	incremented   []string
	logs          []model.Log
	postedAliases []model.Alias
}

// PostAlias records auto-created (catch-all) aliases so tests can assert none were created.
func (b *blocklistStore) PostAlias(ctx context.Context, alias model.Alias, maxDaily int, maxInboundHourly int) (model.Alias, error) {
	b.postedAliases = append(b.postedAliases, alias)
	return alias, nil
}

func (b *blocklistStore) GetVerifiedDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	return b.GetDomains(ctx, userID)
}

func (b *blocklistStore) FindBlocklistMatches(ctx context.Context, userID string, email string, domains []string) ([]model.BlocklistItem, error) {
	if b.matchErr != nil {
		return nil, b.matchErr
	}

	var result []model.BlocklistItem
	for _, item := range b.items {
		if item.UserID != userID {
			continue
		}
		if (item.Type == model.BlocklistEmail && item.Value == email) ||
			(item.Type == model.BlocklistDomain && slices.Contains(domains, item.Value)) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (b *blocklistStore) IncrementBlocklistBlocked(ctx context.Context, ID string) error {
	b.incremented = append(b.incremented, ID)
	return nil
}

func (b *blocklistStore) PostLog(ctx context.Context, lg model.Log) error {
	b.logs = append(b.logs, lg)
	return nil
}

func newBlocklistTestService(items ...model.BlocklistItem) (*Service, *blocklistStore) {
	store := &blocklistStore{fakeStore: newFakeStore(), items: items}
	s := newTestService(store.fakeStore)
	s.Store = store
	return s, store
}

func TestFindBlocklistMatch(t *testing.T) {
	emailItem := model.BlocklistItem{BaseModel: model.BaseModel{ID: "email"}, UserID: "user-1", Type: model.BlocklistEmail, Value: "spammer@example.com"}
	otherEmailItem := model.BlocklistItem{BaseModel: model.BaseModel{ID: "other-email"}, UserID: "user-1", Type: model.BlocklistEmail, Value: "bad@elsewhere.org"}
	domainItem := model.BlocklistItem{BaseModel: model.BaseModel{ID: "domain"}, UserID: "user-1", Type: model.BlocklistDomain, Value: "example.com"}
	subdomainItem := model.BlocklistItem{BaseModel: model.BaseModel{ID: "subdomain"}, UserID: "user-1", Type: model.BlocklistDomain, Value: "news.example.com"}
	otherUserItem := model.BlocklistItem{BaseModel: model.BaseModel{ID: "other"}, UserID: "user-2", Type: model.BlocklistDomain, Value: "other.com"}

	tests := []struct {
		name    string
		from    string
		wantID  string
		blocked bool
	}{
		{"exact email match", "bad@elsewhere.org", "other-email", true},
		{"email on same domain not matched", "good@elsewhere.org", "", false},
		{"email preferred over its domain", "spammer@example.com", "email", true},
		{"email match is case-insensitive", "Spammer@EXAMPLE.com", "email", true},
		{"domain match", "someone@example.com", "domain", true},
		{"subdomain matches parent domain", "someone@mail.example.com", "domain", true},
		{"most specific domain preferred", "someone@news.example.com", "subdomain", true},
		{"lookalike domain not matched", "someone@notexample.com", "", false},
		{"other user's item not matched", "someone@other.com", "", false},
		{"no match", "friend@friendly.org", "", false},
	}

	s, _ := newBlocklistTestService(emailItem, otherEmailItem, domainItem, subdomainItem, otherUserItem)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, blocked := s.FindBlocklistMatch(context.Background(), "user-1", tt.from)
			if blocked != tt.blocked {
				t.Fatalf("blocked = %v, want %v", blocked, tt.blocked)
			}
			if item.ID != tt.wantID {
				t.Errorf("matched item = %q, want %q", item.ID, tt.wantID)
			}
		})
	}
}

func TestFindBlocklistMatch_LookupErrorDoesNotBlock(t *testing.T) {
	s, store := newBlocklistTestService(model.BlocklistItem{UserID: "user-1", Type: model.BlocklistDomain, Value: "example.com"})
	store.matchErr = errors.New("db down")

	if _, blocked := s.FindBlocklistMatch(context.Background(), "user-1", "someone@example.com"); blocked {
		t.Error("expected a failed lookup not to block the message")
	}
}

func TestBlockMessage(t *testing.T) {
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net", UserID: "user-1"}
	item := model.BlocklistItem{BaseModel: model.BaseModel{ID: "item-1"}, UserID: "user-1", Type: model.BlocklistDomain, Value: "example.com"}

	t.Run("with diagnostics enabled", func(t *testing.T) {
		s, store := newBlocklistTestService(item)
		store.settingsByUser["user-1"] = model.Settings{LogIssues: true}

		s.BlockMessage(alias, "someone@example.com", "myalias@mailx.net", item)

		if !slices.Equal(store.incremented, []string{"item-1"}) {
			t.Errorf("expected item-1 blocked count to be incremented, got %v", store.incremented)
		}
		if len(store.postedMessages) != 1 || store.postedMessages[0].Type != model.Block || store.postedMessages[0].AliasID != "alias-1" {
			t.Errorf("expected a Block message recorded for alias-1, got %+v", store.postedMessages)
		}
		if len(store.logs) != 1 {
			t.Fatalf("expected 1 diagnostic log, got %d", len(store.logs))
		}
		lg := store.logs[0]
		if lg.Type != model.BlockedSender || lg.From != "someone@example.com" || lg.Destination != "myalias@mailx.net" {
			t.Errorf("unexpected diagnostic log %+v", lg)
		}
		if lg.Message != "Sender is on your blocklist (domain: example.com)" {
			t.Errorf("unexpected diagnostic log message %q", lg.Message)
		}
	})

	t.Run("with diagnostics disabled", func(t *testing.T) {
		s, store := newBlocklistTestService(item)

		s.BlockMessage(alias, "someone@example.com", "myalias@mailx.net", item)

		if len(store.incremented) != 1 || len(store.postedMessages) != 1 {
			t.Errorf("expected the block to be counted regardless of diagnostics, got incremented=%v messages=%+v", store.incremented, store.postedMessages)
		}
		if len(store.logs) != 0 {
			t.Errorf("expected no diagnostic log, got %+v", store.logs)
		}
	})
}

// blockedSenderMsg is an authenticated message from a sender on the blocked.com domain.
func blockedSenderMsg(to string) []byte {
	return []byte("Authentication-Results: mx.test; dkim=pass header.d=blocked.com\r\n" +
		"From: Spammer <spammer@blocked.com>\r\n" +
		"To: " + to + "\r\n" +
		"Subject: hi\r\n" +
		"Message-ID: <1@blocked.com>\r\n\r\nbody\r\n")
}

func logTypes(logs []model.Log) []model.LogType {
	types := []model.LogType{}
	for _, lg := range logs {
		types = append(types, lg.Type)
	}
	return types
}

// A blocked sender on a catch-all domain is blocked, counted and logged before the catch-all
// alias would be auto-created.
func TestProcessMessage_BlockedSenderOnCatchAllDomain(t *testing.T) {
	item := model.BlocklistItem{BaseModel: model.BaseModel{ID: "item-1"}, UserID: "user-1", Type: model.BlocklistDomain, Value: "blocked.com"}
	s, store := newBlocklistTestService(item)
	store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-1", Enabled: true, CatchAll: true, CreateAlias: true, Recipient: "me@example.com"}
	store.settingsByUser["user-1"] = model.Settings{LogIssues: true}
	store.verifiedRecipients["user-1"] = []model.Recipient{{Email: "me@example.com", IsActive: true}}

	if err := s.ProcessMessage(blockedSenderMsg("random@customdomain.com"), "random@customdomain.com"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !slices.Equal(store.incremented, []string{"item-1"}) {
		t.Errorf("expected item-1 blocked count to be incremented, got %v", store.incremented)
	}
	if !slices.Equal(logTypes(store.logs), []model.LogType{model.BlockedSender}) {
		t.Errorf("expected a single blocked_sender log, got %v", logTypes(store.logs))
	}
	if len(store.postedMessages) != 1 || store.postedMessages[0].Type != model.Block || store.postedMessages[0].UserID != "user-1" {
		t.Errorf("expected one account-level Block message, got %+v", store.postedMessages)
	}
	if len(store.postedAliases) != 0 {
		t.Errorf("expected no catch-all alias to be created, got %+v", store.postedAliases)
	}
}

// A disabled alias discards the message before the blocklist is consulted, so the email gets a
// single reason: one disabled_alias log, one Block stat and no blocklist count.
func TestProcessMessage_DisabledAliasTakesPrecedenceOverBlocklist(t *testing.T) {
	item := model.BlocklistItem{BaseModel: model.BaseModel{ID: "item-1"}, UserID: "user-1", Type: model.BlocklistDomain, Value: "blocked.com"}
	s, store := newBlocklistTestService(item)
	store.aliases["myalias@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net", UserID: "user-1", Enabled: false, Recipients: "me@example.com"}
	store.settingsByUser["user-1"] = model.Settings{LogIssues: true}

	if err := s.ProcessMessage(blockedSenderMsg("myalias@mailx.net"), "myalias@mailx.net"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(store.incremented) != 0 {
		t.Errorf("expected no blocklist count for a disabled alias, got %v", store.incremented)
	}
	if !slices.Equal(logTypes(store.logs), []model.LogType{model.DisabledAlias}) {
		t.Errorf("expected a single disabled_alias log, got %v", logTypes(store.logs))
	}
	if len(store.postedMessages) != 1 {
		t.Errorf("expected exactly one Block message, got %+v", store.postedMessages)
	}
}

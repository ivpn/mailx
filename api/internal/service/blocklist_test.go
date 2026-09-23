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

	items       []model.BlocklistItem
	matchErr    error
	incremented []string
	logs        []model.Log
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

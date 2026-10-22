package service

import (
	"context"
	"testing"
	"time"

	"ivpn.net/email/api/internal/model"
)

// bounceLogStore extends blocklistStore with the methods ProcessBounceLog calls after
// building the log.
type bounceLogStore struct {
	*blocklistStore
}

func (b *bounceLogStore) SaveLogToFile(ctx context.Context, filename string, data []byte) error {
	return nil
}

func (b *bounceLogStore) GetMessagesByAlias(ctx context.Context, aliasID string) ([]model.Message, error) {
	return []model.Message{{ID: 1, Type: model.Send}}, nil
}

func (b *bounceLogStore) DeleteMessage(ctx context.Context, ID uint, userID string) error {
	return nil
}

// Only a bounce with a parseable Date gets an attempted_at; otherwise it's left unset, so the
// API returns null instead of the zero date.
func TestProcessBounceLog_AttemptedAt(t *testing.T) {
	attempted := time.Date(2026, 10, 21, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name string
		date string
		want *time.Time
	}{
		{"with date", "Date: Wed, 21 Oct 2026 10:30:00 +0000\n", &attempted},
		{"without date", "", nil},
		{"with unparseable date", "Date: not a date\n", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, store := newBlocklistTestService()
			store.settingsByUser["user-1"] = model.Settings{LogIssues: true}
			s.Store = &bounceLogStore{store}

			data := []byte("Final-Recipient: rfc822; someone@example.com\n" + tt.date + "Status: 5.1.1\n")
			if err := s.ProcessBounceLog("user-1", "alias-1", data, model.Msg{From: "me@example.com"}); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(store.logs) != 1 {
				t.Fatalf("expected 1 log, got %d", len(store.logs))
			}
			got := store.logs[0].AttemptedAt
			if (got == nil) != (tt.want == nil) || (got != nil && !got.Equal(*tt.want)) {
				t.Errorf("AttemptedAt = %v, want %v", got, tt.want)
			}
		})
	}
}

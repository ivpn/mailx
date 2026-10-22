package service

import (
	"slices"
	"testing"

	"ivpn.net/email/api/internal/model"
)

// A catch-all domain with no recipient to forward to drops the message, and with diagnostics
// enabled the drop is logged instead of being silent.
func TestProcessMessage_NoRecipientsIsLogged(t *testing.T) {
	tests := []struct {
		name      string
		logIssues bool
		wantLogs  []model.LogType
	}{
		{"with diagnostics enabled", true, []model.LogType{model.NoRecipients}},
		{"with diagnostics disabled", false, []model.LogType{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, store := newBlocklistTestService()
			store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-1", Enabled: true, CatchAll: true}
			store.settingsByUser["user-1"] = model.Settings{LogIssues: tt.logIssues}

			if err := s.ProcessMessage(blockedSenderMsg("random@customdomain.com"), "random@customdomain.com"); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !slices.Equal(logTypes(store.logs), tt.wantLogs) {
				t.Fatalf("expected logs %v, got %v", tt.wantLogs, logTypes(store.logs))
			}
			if tt.logIssues {
				lg := store.logs[0]
				if lg.Message != noRecipientsReason || lg.From != "spammer@blocked.com" || lg.Destination != "random@customdomain.com" {
					t.Errorf("unexpected diagnostic log %+v", lg)
				}
			}
			if len(store.postedAliases) != 0 {
				t.Errorf("expected no catch-all alias to be created, got %+v", store.postedAliases)
			}
		})
	}
}

// A message dropped by a disabled alias or domain is logged with a readable reason, not the bare
// sentinel error ("alias disabled:").
func TestProcessMessage_DisabledReasonIsLogged(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(store *blocklistStore)
		to         string
		wantType   model.LogType
		wantReason string
	}{
		{"disabled alias", func(store *blocklistStore) {
			store.aliases["myalias@mailx.net"] = model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net", UserID: "user-1", Enabled: false, Recipients: "me@example.com"}
		}, "myalias@mailx.net", model.DisabledAlias, disabledAliasReason},
		{"disabled catch-all domain", func(store *blocklistStore) {
			store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-1", Enabled: false, CatchAll: true}
		}, "random@customdomain.com", model.DisabledDomain, disabledDomainReason},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, store := newBlocklistTestService()
			store.settingsByUser["user-1"] = model.Settings{LogIssues: true}
			tt.setup(store)

			if err := s.ProcessMessage(blockedSenderMsg(tt.to), tt.to); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !slices.Equal(logTypes(store.logs), []model.LogType{tt.wantType}) {
				t.Fatalf("expected logs %v, got %v", []model.LogType{tt.wantType}, logTypes(store.logs))
			}
			if lg := store.logs[0]; lg.Message != tt.wantReason {
				t.Errorf("expected reason %q, got %q", tt.wantReason, lg.Message)
			}
		})
	}
}

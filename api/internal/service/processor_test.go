package service

import (
	"slices"
	"strings"
	"testing"
	"time"

	"ivpn.net/email/api/internal/model"
)

// pauseStoredAlias sets paused_until on an alias in the store to now + d.
func pauseStoredAlias(store *ruleStore, name string, d time.Duration) time.Time {
	pausedUntil := time.Now().UTC().Truncate(time.Second).Add(d)
	alias := store.aliases[name]
	alias.PausedUntil = &pausedUntil
	store.aliases[name] = alias
	return pausedUntil
}

func forwardEmailTo(to string, subject string) []byte {
	return []byte(strings.Replace(string(forwardEmail(subject)), "To: "+ruleAliasName, "To: "+to, 1))
}

func replyEmail() []byte {
	return []byte("From: Me <me@example.com>\r\n" +
		"To: shop+support=vendor.example@mailx.net\r\n" +
		"Subject: Re: Order\r\n" +
		"In-Reply-To: <1@vendor.example>\r\n" +
		"Authentication-Results: mx.mailx.net; dmarc=pass header.from=example.com\r\n" +
		"\r\nThanks\r\n")
}

// assertPauseDiscarded checks that the last email wasn't queued and was recorded as one
// discarded forward to alias-1.
func assertPauseDiscarded(t *testing.T, store *ruleStore, queued func() []queuedMessage) {
	t.Helper()

	if n := len(queued()); n != 0 {
		t.Errorf("queued %d messages, want the forward discarded", n)
	}
	if len(store.postedMessages) != 1 || store.postedMessages[0].Type != model.Block || store.postedMessages[0].AliasID != "alias-1" {
		t.Errorf("messages = %+v, want one Block for alias-1", store.postedMessages)
	}
}

func TestProcessMessage_PausedAliasDiscardsForwards(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
	pausedUntil := pauseStoredAlias(store, ruleAliasName, 3*time.Hour)

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	assertPauseDiscarded(t, store, queued)
	if len(store.logs) != 1 {
		t.Fatalf("logs = %+v, want one", store.logs)
	}
	lg := store.logs[0]
	wantMessage := "alias paused until: " + pausedUntil.Format(time.RFC3339)
	if lg.Type != model.PausedAlias || lg.Message != wantMessage || lg.From != "news@mail.shop.example" || lg.Destination != ruleAliasName || lg.AliasID != "alias-1" {
		t.Errorf("log = %+v, want a paused_alias log %q", lg, wantMessage)
	}

	t.Run("without diagnostics nothing is logged", func(t *testing.T) {
		store.settingsByUser[ruleUser] = model.Settings{}
		store.logs = nil
		store.postedMessages = nil

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
		if len(store.logs) != 0 {
			t.Errorf("logs = %+v, want none", store.logs)
		}
	})
}

func TestProcessMessage_EndedPauseDelivers(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	pauseStoredAlias(store, ruleAliasName, -time.Second)

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := queued(); len(got) != 1 || got[0].relayType != model.Forward || got[0].rcp.Email != "me@example.com" {
		t.Errorf("queued = %+v, want the forward delivered", got)
	}
}

func TestProcessMessage_PausedAliasStillRepliesAndSends(t *testing.T) {
	tests := []struct {
		name      string
		relayType model.MessageType
		email     []byte
	}{
		{name: "send", relayType: model.Send, email: sendEmail("Order")},
		{name: "reply", relayType: model.Reply, email: replyEmail()},
	}

	for _, tt := range tests {
		relayType, email := tt.relayType, tt.email
		t.Run(tt.name, func(t *testing.T) {
			s, store, queued := newRuleTestService(t)
			pauseStoredAlias(store, ruleAliasName, time.Hour)

			if err := s.ProcessMessage(email, "shop+support=vendor.example@mailx.net"); err != nil {
				t.Fatalf("ProcessMessage() error = %v", err)
			}
			got := queued()
			if len(got) != 1 || got[0].relayType != relayType || got[0].rcp.Email != "support@vendor.example" {
				t.Errorf("queued = %+v, want one message of type %v to support@vendor.example", got, relayType)
			}
		})
	}
}

func TestProcessMessage_PausedAliasMatchedAddresses(t *testing.T) {
	t.Run("a paused Wildcard Alias discards every address it matches", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		delete(store.aliases, ruleAliasName)
		store.aliases["*+news@mailx.net"] = model.Alias{
			BaseModel: model.BaseModel{ID: "alias-1"}, Name: "*+news@mailx.net", UserID: ruleUser,
			Enabled: true, Wildcard: true, Recipients: "me@example.com",
		}
		pauseStoredAlias(store, "*+news@mailx.net", time.Hour)

		for _, to := range []string{"anything+news@mailx.net", "other+news@mailx.net"} {
			store.postedMessages = nil
			if err := s.ProcessMessage(forwardEmailTo(to, "Hello"), to); err != nil {
				t.Fatalf("ProcessMessage(%s) error = %v", to, err)
			}
			assertPauseDiscarded(t, store, queued)
		}
	})

	t.Run("a +tag address follows its paused base alias", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		pauseStoredAlias(store, ruleAliasName, time.Hour)

		if err := s.ProcessMessage(forwardEmailTo("shop+deals@mailx.net", "Hello"), "shop+deals@mailx.net"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
	})

	t.Run("a catch-all address that doesn't exist yet is created and forwarded", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		pauseStoredAlias(store, ruleAliasName, time.Hour)
		store.domains["custom.example"] = model.Domain{Name: "custom.example", UserID: ruleUser, Enabled: true, CatchAll: true, CreateAlias: true, Recipient: "me@example.com"}

		if err := s.ProcessMessage(forwardEmailTo("new@custom.example", "Hello"), "new@custom.example"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(store.postedAlias) != 1 || store.postedAlias[0].Name != "new@custom.example" || store.postedAlias[0].PausedUntil != nil {
			t.Errorf("created %+v, want the unpaused on-the-fly alias", store.postedAlias)
		}
		if got := queued(); len(got) != 1 || got[0].rcp.Email != "me@example.com" {
			t.Errorf("queued = %+v, want the forward delivered", got)
		}
	})
}

func TestProcessMessage_PauseOrderWithOtherChecks(t *testing.T) {
	t.Run("the pause applies while the subscription is limited", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
		store.subscription = model.Subscription{}
		pauseStoredAlias(store, ruleAliasName, time.Hour)

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
		if types := logTypes(store.logs); len(types) != 1 || types[0] != model.PausedAlias {
			t.Errorf("logs = %v, want one paused_alias", types)
		}
	})

	t.Run("a paused alias that is also disabled logs disabled_alias", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
		pauseStoredAlias(store, ruleAliasName, time.Hour)
		alias := store.aliases[ruleAliasName]
		alias.Enabled = false
		store.aliases[ruleAliasName] = alias

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
		if types := logTypes(store.logs); len(types) != 1 || types[0] != model.DisabledAlias {
			t.Errorf("logs = %v, want one disabled_alias", types)
		}
	})

	t.Run("a blocklisted sender is logged as paused_alias only", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
		store.blocklist = []model.BlocklistItem{{BaseModel: model.BaseModel{ID: "block-1"}, UserID: ruleUser, Type: model.BlocklistEmail, Value: "news@mail.shop.example"}}
		pauseStoredAlias(store, ruleAliasName, time.Hour)

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
		if types := logTypes(store.logs); len(types) != 1 || types[0] != model.PausedAlias {
			t.Errorf("logs = %v, want one paused_alias", types)
		}
	})

	t.Run("rules aren't evaluated for a discarded forward", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		rule := store.addRule(model.Rule{Name: "Deactivate", Conditions: subjectContains("hello"),
			Actions: []model.RuleAction{{Type: model.RuleActionDeactivateAlias}, {Type: model.RuleActionBlocklistSender}}})
		pauseStoredAlias(store, ruleAliasName, time.Hour)

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		assertPauseDiscarded(t, store, queued)
		if appliedCount(store, rule.ID) != 0 || len(store.applied) != 0 {
			t.Errorf("applied = %d, want the rule skipped", appliedCount(store, rule.ID))
		}
		if len(store.blocklist) != 0 || !store.aliases[ruleAliasName].Enabled {
			t.Errorf("blocklist = %+v, enabled = %v, want no rule actions run", store.blocklist, store.aliases[ruleAliasName].Enabled)
		}
	})
}

func TestProcessMessage_RecoversQueuePanic(t *testing.T) {
	s, _, _ := newRuleTestService(t)
	queueMessage = func(s *Service, from string, fromName string, to string, toHeader string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType, settings model.Settings, listUnsubscribe model.ListUnsubscribeHeaders, outcome model.RuleOutcome) error {
		panic("boom")
	}

	err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Errorf("ProcessMessage() error = %v, want the panic as an error so Postfix retries", err)
	}
}

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

func TestProcessMessage_NoValidRecipientLogsNoRecipients(t *testing.T) {
	tests := []struct {
		name       string
		to         string
		recipients string // alias recipients
		catchAll   string // catch-all recipient, when to is on custom.example
		wantAlias  string
	}{
		{name: "alias recipient was deleted", to: ruleAliasName, recipients: "gone@example.com", wantAlias: "alias-1"},
		{name: "alias recipient is unverified", to: ruleAliasName, recipients: "pending@example.com", wantAlias: "alias-1"},
		{name: "alias recipient only contains a verified address", to: ruleAliasName, recipients: "name@example.com", wantAlias: "alias-1"},
		{name: "alias has no recipients", to: ruleAliasName, recipients: "", wantAlias: "alias-1"},
		{name: "catch-all recipient is unverified", to: "new@custom.example", catchAll: "pending@example.com"},
		{name: "catch-all recipient was deleted", to: "new@custom.example", catchAll: "gone@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, store, queued := newRuleTestService(t)
			store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
			alias := store.aliases[ruleAliasName]
			alias.Recipients = tt.recipients
			store.aliases[ruleAliasName] = alias
			store.domains["custom.example"] = model.Domain{Name: "custom.example", UserID: ruleUser, Enabled: true, CatchAll: true, CreateAlias: true, Recipient: tt.catchAll}

			if err := s.ProcessMessage(forwardEmailTo(tt.to, "Hello"), tt.to); err != nil {
				t.Fatalf("ProcessMessage() error = %v", err)
			}

			if got := queued(); len(got) != 0 {
				t.Errorf("queued = %+v, want the email dropped", got)
			}
			if len(store.postedAlias) != 0 {
				t.Errorf("created %+v, want no catch-all alias", store.postedAlias)
			}
			if len(store.logs) != 1 {
				t.Fatalf("logs = %+v, want one", store.logs)
			}
			lg := store.logs[0]
			if lg.Type != model.NoRecipients || lg.Message != noRecipientsReason || lg.UserID != ruleUser || lg.AliasID != tt.wantAlias || lg.Destination != tt.to {
				t.Errorf("log = %+v, want a no_recipients log for %s", lg, tt.to)
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

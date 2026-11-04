package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

// ruleStore extends fakeStore with the rule, blocklist, log and alias methods the rules
// feature uses, both through the API and in the mail pipeline.
type ruleStore struct {
	*fakeStore

	mu           sync.Mutex
	rules        []model.Rule
	nextID       int
	applied      [][]string
	blocklist    []model.BlocklistItem
	logs         []model.Log
	postedAlias  []model.Alias
	descriptions map[string]string
	activeErr    error
}

func (r *ruleStore) sorted(userID string, keep func(model.Rule) bool) []model.Rule {
	var result []model.Rule
	for _, rule := range r.rules {
		if rule.UserID == userID && keep(rule) {
			result = append(result, rule)
		}
	}
	slices.SortStableFunc(result, func(a, b model.Rule) int { return a.Order - b.Order })
	return result
}

func (r *ruleStore) GetRules(ctx context.Context, userID string) ([]model.Rule, error) {
	return r.sorted(userID, func(model.Rule) bool { return true }), nil
}

func (r *ruleStore) GetRule(ctx context.Context, ID string, userID string) (model.Rule, error) {
	for _, rule := range r.rules {
		if rule.ID == ID && rule.UserID == userID {
			return rule, nil
		}
	}
	return model.Rule{}, model.ErrRuleNotFound
}

func (r *ruleStore) GetActiveRules(ctx context.Context, userID string, relayType model.MessageType) ([]model.Rule, error) {
	if r.activeErr != nil {
		return nil, r.activeErr
	}
	return r.sorted(userID, func(rule model.Rule) bool { return rule.Active && rule.AppliesTo(relayType) }), nil
}

func (r *ruleStore) GetRuleCount(ctx context.Context, userID string) (int, error) {
	return len(r.sorted(userID, func(model.Rule) bool { return true })), nil
}

func (r *ruleStore) PostRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	r.nextID++
	rule.ID = fmt.Sprintf("rule-%d", r.nextID)
	rule.Order = len(r.sorted(rule.UserID, func(model.Rule) bool { return true }))
	r.rules = append(r.rules, rule)
	return rule, nil
}

func (r *ruleStore) UpdateRule(ctx context.Context, rule model.Rule) error {
	for i, existing := range r.rules {
		if existing.ID == rule.ID && existing.UserID == rule.UserID {
			r.rules[i] = rule
		}
	}
	return nil
}

func (r *ruleStore) UpdateRuleActive(ctx context.Context, ID string, userID string, active bool) error {
	for i, rule := range r.rules {
		if rule.ID == ID && rule.UserID == userID {
			r.rules[i].Active = active
		}
	}
	return nil
}

func (r *ruleStore) ReorderRules(ctx context.Context, userID string, ids []string) error {
	existing := r.sorted(userID, func(model.Rule) bool { return true })
	if len(existing) != len(ids) {
		return model.ErrRuleOrderMismatch
	}
	for _, rule := range existing {
		if !slices.Contains(ids, rule.ID) {
			return model.ErrRuleOrderMismatch
		}
	}
	for i, rule := range r.rules {
		if rule.UserID == userID {
			r.rules[i].Order = slices.Index(ids, rule.ID)
		}
	}
	return nil
}

func (r *ruleStore) IncrementRulesApplied(ctx context.Context, ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.applied = append(r.applied, ids)
	for i, rule := range r.rules {
		if slices.Contains(ids, rule.ID) {
			r.rules[i].Applied++
		}
	}
	return nil
}

func (r *ruleStore) DeleteRule(ctx context.Context, ID string, userID string) error {
	r.rules = slices.DeleteFunc(r.rules, func(rule model.Rule) bool { return rule.ID == ID && rule.UserID == userID })
	return nil
}

func (r *ruleStore) DeleteRulesByUserID(ctx context.Context, userID string) error {
	r.rules = slices.DeleteFunc(r.rules, func(rule model.Rule) bool { return rule.UserID == userID })
	return nil
}

// PostMessage is called from the goroutines that save delivery stats.
func (r *ruleStore) PostMessage(ctx context.Context, message model.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fakeStore.PostMessage(ctx, message)
}

func (r *ruleStore) PostLog(ctx context.Context, lg model.Log) error {
	r.logs = append(r.logs, lg)
	return nil
}

func (r *ruleStore) FindBlocklistMatches(ctx context.Context, userID string, email string, domains []string) ([]model.BlocklistItem, error) {
	var result []model.BlocklistItem
	for _, item := range r.blocklist {
		if item.UserID == userID && ((item.Type == model.BlocklistEmail && item.Value == email) ||
			(item.Type == model.BlocklistDomain && slices.Contains(domains, item.Value))) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *ruleStore) IncrementBlocklistBlocked(ctx context.Context, ID string) error {
	return nil
}

func (r *ruleStore) GetBlocklistCount(ctx context.Context, userID string) (int, error) {
	return len(r.blocklist), nil
}

func (r *ruleStore) PostBlocklistItem(ctx context.Context, item model.BlocklistItem) (model.BlocklistItem, error) {
	for _, existing := range r.blocklist {
		if existing.UserID == item.UserID && existing.Type == item.Type && existing.Value == item.Value {
			return model.BlocklistItem{}, &mysql.MySQLError{Number: 1062}
		}
	}
	item.ID = fmt.Sprintf("block-%d", len(r.blocklist)+1)
	r.blocklist = append(r.blocklist, item)
	return item, nil
}

func (r *ruleStore) UpdateAliasDescription(ctx context.Context, ID string, userID string, description string) error {
	r.descriptions[ID] = description
	return nil
}

func (r *ruleStore) GetVerifiedDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	return r.GetDomains(ctx, userID)
}

func (r *ruleStore) PostAlias(ctx context.Context, alias model.Alias, maxDaily int, maxInboundHourly int) (model.Alias, error) {
	alias.ID = "created-alias"
	r.postedAlias = append(r.postedAlias, alias)
	return alias, nil
}

type queuedMessage struct {
	fromName  string
	rcp       model.Recipient
	alias     model.Alias
	relayType model.MessageType
	settings  model.Settings
	outcome   model.RuleOutcome
}

const (
	ruleUser      = "user-1"
	ruleAliasName = "shop@mailx.net"
	rcpMain       = "rcp-main"
	rcpWork       = "rcp-work"
	rcpPending    = "rcp-pending"
)

func activeSubscription() model.Subscription {
	return model.Subscription{ActiveUntil: time.Now().Add(30 * 24 * time.Hour)}
}

// newRuleTestService sets up user-1 with the alias shop@mailx.net forwarding to me@example.com,
// two other recipients (one unverified) and an active subscription. Messages that would be
// handed to the mailer are captured instead.
func newRuleTestService(t *testing.T) (*Service, *ruleStore, func() []queuedMessage) {
	t.Helper()

	store := &ruleStore{fakeStore: newFakeStore(), descriptions: map[string]string{}}
	store.aliases[ruleAliasName] = model.Alias{
		BaseModel:   model.BaseModel{ID: "alias-1"},
		Name:        ruleAliasName,
		UserID:      ruleUser,
		Enabled:     true,
		Description: "Shopping",
		Recipients:  "me@example.com",
	}
	rcps := []model.Recipient{
		{BaseModel: model.BaseModel{ID: rcpMain}, UserID: ruleUser, Email: "me@example.com", IsActive: true},
		{BaseModel: model.BaseModel{ID: rcpWork}, UserID: ruleUser, Email: "work@example.com", IsActive: true},
		{BaseModel: model.BaseModel{ID: rcpPending}, UserID: ruleUser, Email: "pending@example.com", IsActive: false},
	}
	store.recipients[ruleUser] = rcps
	store.verifiedRecipients[ruleUser] = rcps[:2]
	store.recipients["user-2"] = []model.Recipient{{BaseModel: model.BaseModel{ID: "rcp-other"}, UserID: "user-2", Email: "other@example.com", IsActive: true}}
	store.subscription = activeSubscription()

	s := &Service{
		Cfg: config.Config{API: config.APIConfig{
			Domains:          "mailx.net",
			EmailAuthServIDs: []string{"mx.mailx.net"},
		}},
		Store: store,
	}

	var mu sync.Mutex
	var queued []queuedMessage
	orig := queueMessage
	queueMessage = func(s *Service, from string, fromName string, to string, toHeader string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType, settings model.Settings, listUnsubscribe model.ListUnsubscribeHeaders, outcome model.RuleOutcome) error {
		mu.Lock()
		defer mu.Unlock()
		queued = append(queued, queuedMessage{fromName, rcp, alias, msgType, settings, outcome})
		return nil
	}
	t.Cleanup(func() { queueMessage = orig })

	return s, store, func() []queuedMessage {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(queued)
	}
}

func (r *ruleStore) addRule(rule model.Rule) model.Rule {
	if rule.UserID == "" {
		rule.UserID = ruleUser
	}
	if rule.Operator == "" {
		rule.Operator = model.RuleAnd
	}
	if !rule.Forwards && !rule.Replies && !rule.Sends {
		rule.Forwards = true
	}
	rule.Active = true
	created, _ := r.PostRule(context.Background(), rule)
	return created
}

func subjectContains(value string) []model.RuleCondition {
	return []model.RuleCondition{{Type: model.RuleConditionSubject, Match: model.RuleMatchContains, Values: []string{value}}}
}

func forwardEmail(subject string, extraHeaders ...string) []byte {
	return []byte("From: Shop News <news@mail.shop.example>\r\n" +
		"To: " + ruleAliasName + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Message-ID: <1@shop.example>\r\n" +
		"Authentication-Results: mx.mailx.net; dmarc=pass header.from=shop.example\r\n" +
		strings.Join(extraHeaders, "") +
		"\r\nHello\r\n")
}

func sendEmail(subject string) []byte {
	return []byte("From: Me <me@example.com>\r\n" +
		"To: shop+support=vendor.example@mailx.net\r\n" +
		"Subject: " + subject + "\r\n" +
		"Authentication-Results: mx.mailx.net; dmarc=pass header.from=example.com\r\n" +
		"\r\nHello\r\n")
}

func appliedCount(store *ruleStore, id string) int {
	rule, _ := store.GetRule(context.Background(), id, ruleUser)
	return rule.Applied
}

func TestProcessMessage_NoMatchingRuleDeliversAsBefore(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.addRule(model.Rule{Name: "Bills", Conditions: subjectContains("invoice"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

	if err := s.ProcessMessage(forwardEmail("Hello there"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	got := queued()
	if len(got) != 1 || got[0].rcp.Email != "me@example.com" || got[0].relayType != model.Forward {
		t.Fatalf("queued = %+v, want one forward to me@example.com", got)
	}
	if got[0].outcome.Matched() || len(store.applied) != 0 {
		t.Errorf("outcome = %+v, applied = %v, want no rule applied", got[0].outcome, store.applied)
	}
}

func TestProcessMessage_ForwardToRecipient(t *testing.T) {
	// Acceptance: "if the subject contains invoice, then forward to recipient B"
	s, store, queued := newRuleTestService(t)
	rule := store.addRule(model.Rule{Name: "Bills", Conditions: subjectContains("invoice"),
		Actions: []model.RuleAction{{Type: model.RuleActionForwardTo, Value: rcpWork}}})

	if err := s.ProcessMessage(forwardEmail("Invoice #42"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	got := queued()
	if len(got) != 1 || got[0].rcp.Email != "work@example.com" {
		t.Fatalf("queued = %+v, want only work@example.com", got)
	}
	if n := appliedCount(store, rule.ID); n != 1 {
		t.Errorf("applied = %d, want 1", n)
	}
	if got[0].outcome.MatchedNames[0] != "Bills" {
		t.Errorf("matched = %v", got[0].outcome.MatchedNames)
	}
}

func TestProcessMessage_ForwardToMergesAndFallsBack(t *testing.T) {
	t.Run("recipients from all matching rules are merged", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.addRule(model.Rule{Name: "a", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionForwardTo, Value: rcpWork}}})
		store.addRule(model.Rule{Name: "b", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionForwardTo, Value: rcpMain + "," + rcpWork}}})

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}

		var emails []string
		for _, q := range queued() {
			emails = append(emails, q.rcp.Email)
		}
		slices.Sort(emails)
		if !slices.Equal(emails, []string{"me@example.com", "work@example.com"}) {
			t.Errorf("queued to %v", emails)
		}
		if len(store.applied) != 1 || len(store.applied[0]) != 2 {
			t.Errorf("applied = %v, want both rules counted once, in one update", store.applied)
		}
	})

	t.Run("deleted or unverified recipients fall back to the alias's", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.addRule(model.Rule{Name: "a", Conditions: subjectContains("hello"),
			Actions: []model.RuleAction{{Type: model.RuleActionForwardTo, Value: "deleted," + rcpPending}}})

		if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}

		got := queued()
		if len(got) != 1 || got[0].rcp.Email != "me@example.com" {
			t.Errorf("queued = %+v, want the alias's recipient", got)
		}
	})
}

func TestProcessMessage_AppliedOncePerEmailAndAlias(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	alias := store.aliases[ruleAliasName]
	alias.Recipients = "me@example.com,work@example.com"
	store.aliases[ruleAliasName] = alias
	rule := store.addRule(model.Rule{Name: "Tag", Conditions: subjectContains("hello"),
		Actions: []model.RuleAction{{Type: model.RuleActionReplaceSubject, Value: "[Tag] {subject}"}}})

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	if n := len(queued()); n != 2 {
		t.Fatalf("queued %d messages, want 2", n)
	}
	if n := appliedCount(store, rule.ID); n != 1 {
		t.Errorf("applied = %d, want 1 for two recipients", n)
	}
}

func TestProcessMessage_OrderAndStop(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	a := store.addRule(model.Rule{Name: "A", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionReplaceSubject, Value: "[A] {subject}"}}})
	b := store.addRule(model.Rule{Name: "B", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionReplaceSubject, Value: "[B] {subject}"}}})

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := queued()[0].outcome.Subject("Hello"); got != "[B] [A] Hello" {
		t.Errorf("subject = %q, want [B] [A] Hello", got)
	}

	// Acceptance: dragging the bottom rule to the top changes the next email's subject
	if err := s.ReorderRules(context.Background(), ruleUser, []string{b.ID, a.ID}); err != nil {
		t.Fatalf("ReorderRules() error = %v", err)
	}
	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := queued()[1].outcome.Subject("Hello"); got != "[A] [B] Hello" {
		t.Errorf("subject after reorder = %q, want [A] [B] Hello", got)
	}

	// stop on the top rule keeps the one below from being applied or counted
	for i, rule := range store.rules {
		if rule.ID == b.ID {
			store.rules[i].Actions = append(store.rules[i].Actions, model.RuleAction{Type: model.RuleActionStop})
		}
	}
	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := queued()[2].outcome; got.Subject("Hello") != "[B] Hello" || !slices.Equal(got.MatchedIDs, []string{b.ID}) {
		t.Errorf("outcome with stop = %+v", got)
	}
	if appliedCount(store, a.ID) != 2 || appliedCount(store, b.ID) != 3 {
		t.Errorf("applied a=%d b=%d, want 2 and 3", appliedCount(store, a.ID), appliedCount(store, b.ID))
	}
}

func TestProcessMessage_SkipsInactiveRulesAndOtherEmailTypes(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	inactive := store.addRule(model.Rule{Name: "inactive", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})
	sends := store.addRule(model.Rule{Name: "sends", Sends: true, Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})
	if err := store.UpdateRuleActive(context.Background(), inactive.ID, ruleUser, false); err != nil {
		t.Fatal(err)
	}

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	if len(queued()) != 1 || appliedCount(store, inactive.ID) != 0 || appliedCount(store, sends.ID) != 0 {
		t.Errorf("queued=%d inactive=%d sends=%d, want the forward delivered and no rule applied",
			len(queued()), appliedCount(store, inactive.ID), appliedCount(store, sends.ID))
	}

	// Acceptance: reactivating the rule applies it again, without resetting the count
	if err := s.UpdateRuleActive(context.Background(), inactive.ID, ruleUser, true); err != nil {
		t.Fatalf("UpdateRuleActive() error = %v", err)
	}
	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if len(queued()) != 1 || appliedCount(store, inactive.ID) != 1 {
		t.Errorf("after reactivating: queued=%d applied=%d, want the email blocked and counted", len(queued()), appliedCount(store, inactive.ID))
	}
}

func TestProcessMessage_BlockByRule(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
	rule := store.addRule(model.Rule{Name: `Drop "promos"`, Conditions: subjectContains("promo"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

	if err := s.ProcessMessage(forwardEmail("Big promo"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	if n := len(queued()); n != 0 {
		t.Fatalf("queued %d messages, want the email dropped", n)
	}
	if len(store.postedMessages) != 1 || store.postedMessages[0].Type != model.Block || store.postedMessages[0].AliasID != "alias-1" {
		t.Errorf("messages = %+v, want one Block for alias-1", store.postedMessages)
	}
	if len(store.logs) != 1 {
		t.Fatalf("logs = %+v, want one", store.logs)
	}
	lg := store.logs[0]
	if lg.Type != model.BlockedByRule || lg.Message != `Blocked by rule "Drop "promos""` || lg.From != "news@mail.shop.example" || lg.Destination != ruleAliasName {
		t.Errorf("log = %+v", lg)
	}
	if appliedCount(store, rule.ID) != 1 {
		t.Errorf("applied = %d, want the blocked email counted", appliedCount(store, rule.ID))
	}

	t.Run("without diagnostics nothing is logged", func(t *testing.T) {
		store.settingsByUser[ruleUser] = model.Settings{}
		store.logs = nil
		if err := s.ProcessMessage(forwardEmail("Big promo"), ruleAliasName); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(store.logs) != 0 || len(store.postedMessages) != 2 {
			t.Errorf("logs=%v messages=%d", store.logs, len(store.postedMessages))
		}
	})
}

func TestProcessMessage_BlocklistedSenderSkipsRules(t *testing.T) {
	// Acceptance: a rule adds the newsletter's domain to the blocklist and blocks it; later
	// emails from that domain are blocked by the blocklist without counting the rule again
	s, store, queued := newRuleTestService(t)
	rule := store.addRule(model.Rule{
		Name:       "Newsletters",
		Conditions: []model.RuleCondition{{Type: model.RuleConditionIsMailingList, Values: []string{}}},
		Actions:    []model.RuleAction{{Type: model.RuleActionBlocklistDomain}, {Type: model.RuleActionBlock}},
	})
	newsletter := forwardEmail("Weekly", "List-Unsubscribe: <https://shop.example/u>\r\n")

	if err := s.ProcessMessage(newsletter, ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if len(queued()) != 0 {
		t.Fatal("the newsletter was delivered")
	}
	if len(store.blocklist) != 1 || store.blocklist[0].Type != model.BlocklistDomain || store.blocklist[0].Value != "mail.shop.example" {
		t.Fatalf("blocklist = %+v, want the sender domain", store.blocklist)
	}

	if err := s.ProcessMessage(newsletter, ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if len(queued()) != 0 || appliedCount(store, rule.ID) != 1 || len(store.blocklist) != 1 {
		t.Errorf("queued=%d applied=%d blocklist=%d, want the blocklist to block it before rules run",
			len(queued()), appliedCount(store, rule.ID), len(store.blocklist))
	}
}

func TestProcessMessage_SideEffects(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.addRule(model.Rule{Name: "Housekeeping", Conditions: subjectContains("hello"), Actions: []model.RuleAction{
		{Type: model.RuleActionBlocklistSender},
		{Type: model.RuleActionDeactivateAlias},
		{Type: model.RuleActionSetAliasDescription, Value: "Spammed"},
	}})
	// Already on the blocklist under another type: the duplicate path must be silent
	store.blocklist = append(store.blocklist, model.BlocklistItem{UserID: "user-9", Type: model.BlocklistEmail, Value: "x@y.z"})

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	// The current email is still delivered
	if len(queued()) != 1 {
		t.Fatalf("queued = %d, want the email delivered", len(queued()))
	}
	if !slices.ContainsFunc(store.blocklist, func(i model.BlocklistItem) bool {
		return i.UserID == ruleUser && i.Type == model.BlocklistEmail && i.Value == "news@mail.shop.example"
	}) {
		t.Errorf("blocklist = %+v, want the sender email", store.blocklist)
	}
	if store.aliases[ruleAliasName].Enabled {
		t.Error("alias still enabled")
	}
	if store.descriptions["alias-1"] != "Spammed" {
		t.Errorf("descriptions = %v", store.descriptions)
	}
}

func TestApplyRuleSideEffects_IgnoresDuplicatesAndFullBlocklist(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	alias := store.aliases[ruleAliasName]
	outcome := model.RuleOutcome{BlocklistSender: true, BlocklistDomain: true}

	s.ApplyRuleSideEffects(context.Background(), alias, "News@Example.com", outcome)
	s.ApplyRuleSideEffects(context.Background(), alias, "news@example.com", outcome)
	if len(store.blocklist) != 2 {
		t.Fatalf("blocklist = %+v, want one email and one domain entry", store.blocklist)
	}

	for len(store.blocklist) < model.MaxBlocklistItems {
		store.blocklist = append(store.blocklist, model.BlocklistItem{UserID: ruleUser, Type: model.BlocklistEmail, Value: fmt.Sprint(len(store.blocklist))})
	}
	s.ApplyRuleSideEffects(context.Background(), alias, "new@other.com", outcome)
	if len(store.blocklist) != model.MaxBlocklistItems {
		t.Errorf("blocklist grew past its limit to %d", len(store.blocklist))
	}

	// Invalid senders are skipped
	s.ApplyRuleSideEffects(context.Background(), alias, "not-an-address", outcome)
}

func TestProcessMessage_BlockStopsOnTheFlyAlias(t *testing.T) {
	setup := func(t *testing.T) (*Service, *ruleStore) {
		s, store, _ := newRuleTestService(t)
		store.domains["custom.example"] = model.Domain{Name: "custom.example", UserID: ruleUser, Enabled: true, CatchAll: true, CreateAlias: true, Recipient: "me@example.com"}
		return s, store
	}
	email := func(subject string) []byte {
		return []byte(strings.Replace(string(forwardEmail(subject)), "To: "+ruleAliasName, "To: new@custom.example", 1))
	}

	t.Run("blocked", func(t *testing.T) {
		s, store := setup(t)
		store.addRule(model.Rule{Name: "Block", Conditions: subjectContains("promo"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

		if err := s.ProcessMessage(email("Promo"), "new@custom.example"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(store.postedAlias) != 0 {
			t.Errorf("created %+v, want no alias", store.postedAlias)
		}
	})

	t.Run("not blocked", func(t *testing.T) {
		s, store := setup(t)
		store.addRule(model.Rule{Name: "Describe", Conditions: subjectContains("promo"), Actions: []model.RuleAction{{Type: model.RuleActionSetAliasDescription, Value: "From promo"}}})

		if err := s.ProcessMessage(email("Promo"), "new@custom.example"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(store.postedAlias) != 1 || store.postedAlias[0].Name != "new@custom.example" {
			t.Fatalf("created %+v, want the on-the-fly alias", store.postedAlias)
		}
		// Alias actions run once the alias is saved
		if store.descriptions["created-alias"] != "From promo" {
			t.Errorf("descriptions = %v", store.descriptions)
		}
	})
}

func TestProcessMessage_InactiveSubscriptionStillAppliesRulesToForwards(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.subscription = model.Subscription{}
	rule := store.addRule(model.Rule{Name: "Block", Conditions: subjectContains("promo"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

	if err := s.ProcessMessage(forwardEmail("Promo"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if len(queued()) != 0 || appliedCount(store, rule.ID) != 1 {
		t.Errorf("queued=%d applied=%d, want the forward blocked", len(queued()), appliedCount(store, rule.ID))
	}
}

func TestProcessMessage_RepliesAndSends(t *testing.T) {
	t.Run("destination condition and overrides", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.settingsByUser[ruleUser] = model.Settings{FromName: "Default"}
		rule := store.addRule(model.Rule{Name: "Support", Sends: true,
			Conditions: []model.RuleCondition{{Type: model.RuleConditionDestination, Match: model.RuleMatchEndsWith, Values: []string{"@vendor.example"}}},
			Actions:    []model.RuleAction{{Type: model.RuleActionSetDisplayName, Value: "Customer"}, {Type: model.RuleActionReplaceSubject, Value: "[Ticket] {subject}"}},
		})

		if err := s.ProcessMessage(sendEmail("Help"), "shop+support=vendor.example@mailx.net"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}

		got := queued()
		if len(got) != 1 || got[0].relayType != model.Send || got[0].rcp.Email != "support@vendor.example" {
			t.Fatalf("queued = %+v, want one send to support@vendor.example", got)
		}
		if got[0].outcome.FromName != "Customer" || got[0].outcome.Subject("Help") != "[Ticket] Help" {
			t.Errorf("outcome = %+v", got[0].outcome)
		}
		if appliedCount(store, rule.ID) != 1 {
			t.Errorf("applied = %d", appliedCount(store, rule.ID))
		}
	})

	t.Run("forward rules don't apply to sends", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		rule := store.addRule(model.Rule{Name: "Block", Conditions: subjectContains("help"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

		if err := s.ProcessMessage(sendEmail("Help"), "shop+support=vendor.example@mailx.net"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(queued()) != 1 || appliedCount(store, rule.ID) != 0 {
			t.Errorf("queued=%d applied=%d", len(queued()), appliedCount(store, rule.ID))
		}
	})

	t.Run("a blocked send isn't sent and is logged", func(t *testing.T) {
		s, store, queued := newRuleTestService(t)
		store.settingsByUser[ruleUser] = model.Settings{LogIssues: true}
		store.addRule(model.Rule{Name: "No vendor", Sends: true, Replies: true, Conditions: subjectContains("help"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

		if err := s.ProcessMessage(sendEmail("Help"), "shop+support=vendor.example@mailx.net"); err != nil {
			t.Fatalf("ProcessMessage() error = %v", err)
		}
		if len(queued()) != 0 {
			t.Fatal("the blocked send was queued")
		}
		if len(store.logs) != 1 || store.logs[0].Type != model.BlockedByRule || store.logs[0].Message != `Blocked by rule "No vendor"` {
			t.Errorf("logs = %+v", store.logs)
		}
		if len(store.postedMessages) != 0 {
			t.Errorf("messages = %+v, want no Block stat for a send", store.postedMessages)
		}
	})
}

func TestProcessMessage_RuleLookupFailureDelivers(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.addRule(model.Rule{Name: "Block", Conditions: subjectContains("promo"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})
	store.activeErr = errors.New("db down")

	if err := s.ProcessMessage(forwardEmail("Promo"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if len(queued()) != 1 || queued()[0].outcome.Matched() {
		t.Errorf("queued = %+v, want the email delivered without rules", queued())
	}
}

func TestProcessMessage_BannerOverride(t *testing.T) {
	s, store, queued := newRuleTestService(t)
	store.settingsByUser[ruleUser] = model.Settings{RemoveHeader: false, BlockTrackers: true}
	store.addRule(model.Rule{Name: "Hide", Conditions: subjectContains("hello"), Actions: []model.RuleAction{{Type: model.RuleActionHideBanner}, {Type: model.RuleActionSetDisplayName, Value: "Shop"}}})

	if err := s.ProcessMessage(forwardEmail("Hello"), ruleAliasName); err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	got := queued()[0]
	if !got.settings.RemoveHeader || !got.settings.BlockTrackers {
		t.Errorf("settings = %+v, want the banner hidden and the rest kept", got.settings)
	}
	if got.outcome.FromName != "Shop" {
		t.Errorf("FromName = %q", got.outcome.FromName)
	}
}

func TestApplyRules_AttachmentsParsedOnlyWhenNeeded(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	alias := store.aliases[ruleAliasName]
	msg := model.Msg{From: "a@example.com", Subject: "Report"}
	attachment := []byte(strings.Join([]string{
		"Content-Type: multipart/mixed; boundary=b",
		"",
		"--b",
		"Content-Type: application/pdf",
		"Content-Disposition: attachment; filename=r.pdf",
		"",
		"x",
		"--b--",
		"",
	}, "\r\n"))

	rule := store.addRule(model.Rule{Name: "Files", Conditions: []model.RuleCondition{{Type: model.RuleConditionHasAttachments, Values: []string{}}},
		Actions: []model.RuleAction{{Type: model.RuleActionRemoveAttachments}}})

	outcome := s.ApplyRules(context.Background(), alias, msg, model.Forward, ruleAliasName, attachment)
	if !outcome.RemoveAttachments || appliedCount(store, rule.ID) != 1 {
		t.Errorf("outcome = %+v, want the attachment rule applied", outcome)
	}

	outcome = s.ApplyRules(context.Background(), alias, msg, model.Forward, ruleAliasName, []byte("Subject: x\r\n\r\nno attachments"))
	if outcome.Matched() {
		t.Errorf("outcome = %+v, want no match without attachments", outcome)
	}
}

func newRuleReq(name string) model.Rule {
	return model.Rule{
		UserID:     ruleUser,
		Name:       name,
		Operator:   model.RuleAnd,
		Conditions: subjectContains("invoice"),
		Actions:    []model.RuleAction{{Type: model.RuleActionForwardTo, Value: rcpWork}},
		Forwards:   true,
	}
}

func TestPostRule(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	ctx := context.Background()

	first, err := s.PostRule(ctx, newRuleReq("  First "))
	if err != nil {
		t.Fatalf("PostRule() error = %v", err)
	}
	second, err := s.PostRule(ctx, newRuleReq("Second"))
	if err != nil {
		t.Fatalf("PostRule() error = %v", err)
	}
	if first.Name != "First" || !first.Active || first.Order != 0 || second.Order != 1 {
		t.Errorf("created %+v and %+v", first, second)
	}

	t.Run("validation errors are returned", func(t *testing.T) {
		req := newRuleReq("Bad")
		req.Replies = true
		if _, err := s.PostRule(ctx, req); err == nil || err.Error() != "Forward to recipient(s) only applies to Forwards." {
			t.Errorf("PostRule() error = %v", err)
		}
	})

	t.Run("forward_to recipients must be the user's verified ones", func(t *testing.T) {
		for _, id := range []string{"rcp-other", rcpPending, "missing"} {
			req := newRuleReq("Bad")
			req.Actions[0].Value = rcpMain + "," + id
			if _, err := s.PostRule(ctx, req); !errors.Is(err, ErrRuleRecipient) {
				t.Errorf("PostRule() with %s error = %v, want ErrRuleRecipient", id, err)
			}
		}
	})

	t.Run("limit", func(t *testing.T) {
		for len(store.rules) < model.MaxRules {
			store.addRule(newRuleReq("filler"))
		}
		if _, err := s.PostRule(ctx, newRuleReq("One too many")); !errors.Is(err, ErrPostRuleLimit) {
			t.Errorf("PostRule() error = %v, want ErrPostRuleLimit", err)
		}
	})
}

func TestUpdateRule(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	ctx := context.Background()
	rule := store.addRule(newRuleReq("Original"))
	store.addRule(newRuleReq("Second"))
	if err := store.IncrementRulesApplied(ctx, []string{rule.ID}); err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateRuleActive(ctx, rule.ID, ruleUser, false); err != nil {
		t.Fatal(err)
	}

	update := newRuleReq("Renamed")
	update.ID = rule.ID
	update.Actions = []model.RuleAction{{Type: model.RuleActionBlock}}
	update.Order, update.Active, update.Applied = 5, true, 0

	got, err := s.UpdateRule(ctx, update)
	if err != nil {
		t.Fatalf("UpdateRule() error = %v", err)
	}
	if got.Name != "Renamed" || got.Actions[0].Type != model.RuleActionBlock || got.Order != 0 || got.Active || got.Applied != 1 {
		t.Errorf("updated = %+v, want the definition replaced and order/active/applied kept", got)
	}

	t.Run("another user's rule is not found", func(t *testing.T) {
		other := newRuleReq("Hijack")
		other.ID = rule.ID
		other.UserID = "user-2"
		other.Actions = []model.RuleAction{{Type: model.RuleActionBlock}}
		if _, err := s.UpdateRule(ctx, other); !errors.Is(err, model.ErrRuleNotFound) {
			t.Errorf("UpdateRule() error = %v, want ErrRuleNotFound", err)
		}
		if err := s.UpdateRuleActive(ctx, rule.ID, "user-2", true); !errors.Is(err, model.ErrRuleNotFound) {
			t.Errorf("UpdateRuleActive() error = %v, want ErrRuleNotFound", err)
		}
		if err := s.DeleteRule(ctx, rule.ID, "user-2"); !errors.Is(err, model.ErrRuleNotFound) {
			t.Errorf("DeleteRule() error = %v, want ErrRuleNotFound", err)
		}
		if len(store.rules) != 2 {
			t.Errorf("rules = %d, want 2", len(store.rules))
		}
	})
}

func TestReorderAndDeleteRules(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	ctx := context.Background()
	a := store.addRule(newRuleReq("a"))
	b := store.addRule(newRuleReq("b"))

	if err := s.ReorderRules(ctx, ruleUser, []string{a.ID}); !errors.Is(err, ErrReorderRulesStale) {
		t.Errorf("ReorderRules() with a missing ID error = %v, want ErrReorderRulesStale", err)
	}
	if err := s.ReorderRules(ctx, ruleUser, []string{b.ID, a.ID}); err != nil {
		t.Fatalf("ReorderRules() error = %v", err)
	}
	rules, _ := s.GetRules(ctx, ruleUser)
	if rules[0].ID != b.ID {
		t.Errorf("order = %v", ruleIDsOf(rules))
	}

	if err := s.DeleteRule(ctx, a.ID, ruleUser); err != nil {
		t.Fatalf("DeleteRule() error = %v", err)
	}
	rules, _ = s.GetRules(ctx, ruleUser)
	if len(rules) != 1 || rules[0].ID != b.ID {
		t.Errorf("rules after delete = %v", ruleIDsOf(rules))
	}
}

func ruleIDsOf(rules []model.Rule) []string {
	ids := make([]string, len(rules))
	for i, r := range rules {
		ids[i] = r.ID
	}
	return ids
}

func TestGetRules_FlagsInvalidRecipients(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	ok := store.addRule(newRuleReq("ok"))
	deleted := newRuleReq("deleted")
	deleted.Actions[0].Value = rcpMain + ",gone"
	deleted = store.addRule(deleted)
	unverified := newRuleReq("unverified")
	unverified.Actions[0].Value = rcpPending
	store.addRule(unverified)
	store.addRule(model.Rule{Name: "no forward", Conditions: subjectContains("x"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

	rules, err := s.GetRules(context.Background(), ruleUser)
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}

	want := map[string]bool{ok.ID: false, deleted.ID: true}
	for _, r := range rules {
		if expected, checked := want[r.ID]; checked && r.InvalidRecipients != expected {
			t.Errorf("rule %s invalid_recipients = %v, want %v", r.Name, r.InvalidRecipients, expected)
		}
		if r.Name == "unverified" && !r.InvalidRecipients {
			t.Error("a rule pointing at an unverified recipient should be flagged")
		}
		if r.Name == "no forward" && r.InvalidRecipients {
			t.Error("a rule without forward_to should not be flagged")
		}
	}
}

func TestRules_InactiveSubscriptionRejectsChanges(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	ctx := context.Background()
	rule := store.addRule(newRuleReq("Existing"))
	other := store.addRule(newRuleReq("Other"))
	store.subscription = model.Subscription{}

	update := newRuleReq("Changed")
	update.ID = rule.ID

	if _, err := s.PostRule(ctx, newRuleReq("New")); !errors.Is(err, ErrPostRuleInactiveSub) {
		t.Errorf("PostRule() error = %v", err)
	}
	if _, err := s.UpdateRule(ctx, update); !errors.Is(err, ErrUpdateRuleInactiveSub) {
		t.Errorf("UpdateRule() error = %v", err)
	}
	if err := s.UpdateRuleActive(ctx, rule.ID, ruleUser, false); !errors.Is(err, ErrUpdateRuleInactiveSub) {
		t.Errorf("UpdateRuleActive() error = %v", err)
	}
	if err := s.ReorderRules(ctx, ruleUser, []string{other.ID, rule.ID}); !errors.Is(err, ErrReorderRulesInactiveSub) {
		t.Errorf("ReorderRules() error = %v", err)
	}
	if err := s.DeleteRule(ctx, rule.ID, ruleUser); !errors.Is(err, ErrDeleteRuleInactiveSub) {
		t.Errorf("DeleteRule() error = %v", err)
	}

	// Nothing changed, and the rules can still be listed
	rules, err := s.GetRules(ctx, ruleUser)
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}
	if len(rules) != 2 || rules[0].ID != rule.ID || rules[0].Name != "Existing" || !rules[0].Active {
		t.Errorf("rules = %+v, want them unchanged", rules)
	}
}

func TestDeleteUser_DeletesRules(t *testing.T) {
	s, store, _ := newRuleTestService(t)
	store.addRule(newRuleReq("mine"))
	store.addRule(model.Rule{UserID: "user-2", Name: "theirs", Conditions: subjectContains("x"), Actions: []model.RuleAction{{Type: model.RuleActionBlock}}})

	if err := s.DeleteRulesByUserID(context.Background(), ruleUser); err != nil {
		t.Fatalf("DeleteRulesByUserID() error = %v", err)
	}
	if len(store.rules) != 1 || store.rules[0].UserID != "user-2" {
		t.Errorf("rules = %+v, want only user-2's", store.rules)
	}
}

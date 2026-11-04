package repository

import (
	"context"
	"errors"
	"slices"
	"testing"

	"ivpn.net/email/api/internal/model"
)

func newRuleTestDB(t *testing.T) *Database {
	t.Helper()

	d := newTestDB(t)
	if err := d.Client.AutoMigrate(&model.Rule{}); err != nil {
		t.Fatalf("failed to migrate rules: %v", err)
	}

	return d
}

func testRule(userID string, name string) model.Rule {
	return model.Rule{
		UserID:     userID,
		Name:       name,
		Operator:   model.RuleOr,
		Conditions: []model.RuleCondition{{Type: model.RuleConditionSubject, Match: model.RuleMatchContains, Values: []string{"invoice", "receipt"}}},
		Actions:    []model.RuleAction{{Type: model.RuleActionReplaceSubject, Value: "[Bills] {subject}"}, {Type: model.RuleActionStop}},
		Forwards:   true,
		Active:     true,
	}
}

func createRules(t *testing.T, d *Database, userID string, names ...string) []model.Rule {
	t.Helper()

	var rules []model.Rule
	for _, name := range names {
		rule, err := d.PostRule(context.Background(), testRule(userID, name))
		if err != nil {
			t.Fatalf("PostRule(%s) error = %v", name, err)
		}
		rules = append(rules, rule)
	}

	return rules
}

func ruleNames(rules []model.Rule) []string {
	names := make([]string, len(rules))
	for i, r := range rules {
		names[i] = r.Name
	}
	return names
}

func ruleIDs(rules []model.Rule) []string {
	ids := make([]string, len(rules))
	for i, r := range rules {
		ids[i] = r.ID
	}
	return ids
}

func TestRules_CreateAppendsAndRoundTrips(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()

	created := createRules(t, d, "user-1", "first", "second", "third")
	createRules(t, d, "user-2", "other")

	for i, r := range created {
		if r.ID == "" || r.Order != i {
			t.Errorf("rule %s: id=%q order=%d, want an ID and order %d", r.Name, r.ID, r.Order, i)
		}
	}

	rules, err := d.GetRules(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}
	if got := ruleNames(rules); !slices.Equal(got, []string{"first", "second", "third"}) {
		t.Fatalf("GetRules() = %v", got)
	}

	r := rules[0]
	if r.Operator != model.RuleOr || !r.Forwards || r.Replies || r.Sends || !r.Active || r.Applied != 0 || r.LastAppliedAt != nil {
		t.Errorf("stored rule = %+v", r)
	}
	if len(r.Conditions) != 1 || !slices.Equal(r.Conditions[0].Values, []string{"invoice", "receipt"}) || r.Conditions[0].Match != model.RuleMatchContains {
		t.Errorf("conditions = %+v", r.Conditions)
	}
	if len(r.Actions) != 2 || r.Actions[0].Value != "[Bills] {subject}" || r.Actions[1].Type != model.RuleActionStop {
		t.Errorf("actions = %+v", r.Actions)
	}

	count, err := d.GetRuleCount(ctx, "user-1")
	if err != nil || count != 3 {
		t.Errorf("GetRuleCount() = %d, %v, want 3", count, err)
	}
}

func TestRules_GetRuleChecksOwnership(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()
	rule := createRules(t, d, "user-1", "mine")[0]

	got, err := d.GetRule(ctx, rule.ID, "user-1")
	if err != nil || got.Name != "mine" {
		t.Errorf("GetRule() = %+v, %v", got, err)
	}

	if _, err := d.GetRule(ctx, rule.ID, "user-2"); !errors.Is(err, model.ErrRuleNotFound) {
		t.Errorf("GetRule() for another user error = %v, want ErrRuleNotFound", err)
	}
	if _, err := d.GetRule(ctx, "missing", "user-1"); !errors.Is(err, model.ErrRuleNotFound) {
		t.Errorf("GetRule() for a missing ID error = %v, want ErrRuleNotFound", err)
	}
}

func TestRules_UpdateKeepsOrderActiveAndCounters(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()
	rules := createRules(t, d, "user-1", "first", "second")

	if err := d.UpdateRuleActive(ctx, rules[1].ID, "user-1", false); err != nil {
		t.Fatalf("UpdateRuleActive() error = %v", err)
	}
	if err := d.IncrementRulesApplied(ctx, []string{rules[1].ID}); err != nil {
		t.Fatalf("IncrementRulesApplied() error = %v", err)
	}

	update := rules[1]
	update.Name = "renamed"
	update.Operator = model.RuleAnd
	update.Conditions = []model.RuleCondition{{Type: model.RuleConditionIsMailingList, Values: []string{}}}
	update.Actions = []model.RuleAction{{Type: model.RuleActionBlock}}
	update.Forwards = false
	update.Replies = true
	update.Sends = true
	// Fields the update must not touch
	update.Order = 99
	update.Active = true
	update.Applied = 0
	if err := d.UpdateRule(ctx, update); err != nil {
		t.Fatalf("UpdateRule() error = %v", err)
	}

	got, _ := d.GetRule(ctx, rules[1].ID, "user-1")
	if got.Name != "renamed" || got.Operator != model.RuleAnd || got.Forwards || !got.Replies || !got.Sends {
		t.Errorf("updated rule = %+v", got)
	}
	if len(got.Actions) != 1 || got.Actions[0].Type != model.RuleActionBlock || got.Conditions[0].Type != model.RuleConditionIsMailingList {
		t.Errorf("updated conditions/actions = %+v %+v", got.Conditions, got.Actions)
	}
	if got.Order != 1 || got.Active || got.Applied != 1 || got.LastAppliedAt == nil {
		t.Errorf("order=%d active=%v applied=%d last=%v, want them kept", got.Order, got.Active, got.Applied, got.LastAppliedAt)
	}

	// Another user's ID changes nothing
	other := got
	other.UserID = "user-2"
	other.Name = "hijacked"
	if err := d.UpdateRule(ctx, other); err != nil {
		t.Fatalf("UpdateRule() for another user error = %v", err)
	}
	if err := d.UpdateRuleActive(ctx, got.ID, "user-2", true); err != nil {
		t.Fatalf("UpdateRuleActive() for another user error = %v", err)
	}
	got, _ = d.GetRule(ctx, rules[1].ID, "user-1")
	if got.Name != "renamed" || got.Active {
		t.Errorf("another user changed the rule: %+v", got)
	}
}

func TestRules_GetActiveRules(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()

	forwards := testRule("user-1", "forwards")
	inactive := testRule("user-1", "inactive")
	replies := testRule("user-1", "replies and sends")
	replies.Forwards, replies.Replies, replies.Sends = false, true, true
	for _, r := range []model.Rule{forwards, inactive, replies} {
		if _, err := d.PostRule(ctx, r); err != nil {
			t.Fatalf("PostRule() error = %v", err)
		}
	}
	createRules(t, d, "user-2", "other user")

	all, _ := d.GetRules(ctx, "user-1")
	if err := d.UpdateRuleActive(ctx, all[1].ID, "user-1", false); err != nil {
		t.Fatalf("UpdateRuleActive() error = %v", err)
	}
	// Reverse the order to check the active rules come back in evaluation order
	if err := d.ReorderRules(ctx, "user-1", []string{all[2].ID, all[1].ID, all[0].ID}); err != nil {
		t.Fatalf("ReorderRules() error = %v", err)
	}

	tests := []struct {
		relayType model.MessageType
		want      []string
	}{
		{model.Forward, []string{"forwards"}},
		{model.Reply, []string{"replies and sends"}},
		{model.Send, []string{"replies and sends"}},
		{model.FailBounce, []string{}},
	}
	for _, tt := range tests {
		rules, err := d.GetActiveRules(ctx, "user-1", tt.relayType)
		if err != nil {
			t.Fatalf("GetActiveRules(%d) error = %v", tt.relayType, err)
		}
		if got := ruleNames(rules); !slices.Equal(got, tt.want) {
			t.Errorf("GetActiveRules(%d) = %v, want %v", tt.relayType, got, tt.want)
		}
	}

	// Ordering across several active rules
	if err := d.UpdateRuleActive(ctx, all[1].ID, "user-1", true); err != nil {
		t.Fatalf("UpdateRuleActive() error = %v", err)
	}
	rules, _ := d.GetActiveRules(ctx, "user-1", model.Forward)
	if got := ruleNames(rules); !slices.Equal(got, []string{"inactive", "forwards"}) {
		t.Errorf("GetActiveRules() order = %v", got)
	}
}

func TestRules_Reorder(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()
	rules := createRules(t, d, "user-1", "a", "b", "c")
	other := createRules(t, d, "user-2", "x")[0]
	ids := ruleIDs(rules)

	if err := d.ReorderRules(ctx, "user-1", []string{ids[2], ids[0], ids[1]}); err != nil {
		t.Fatalf("ReorderRules() error = %v", err)
	}
	got, _ := d.GetRules(ctx, "user-1")
	if names := ruleNames(got); !slices.Equal(names, []string{"c", "a", "b"}) {
		t.Fatalf("order after reorder = %v", names)
	}
	for i, r := range got {
		if r.Order != i {
			t.Errorf("rule %s order = %d, want %d", r.Name, r.Order, i)
		}
	}

	rejected := map[string][]string{
		"missing ID":        {ids[0], ids[1]},
		"extra ID":          {ids[0], ids[1], ids[2], "extra"},
		"another user's ID": {ids[0], ids[1], other.ID},
		"duplicate ID":      {ids[0], ids[1], ids[1]},
		"empty":             {},
	}
	for name, attempt := range rejected {
		t.Run(name, func(t *testing.T) {
			if err := d.ReorderRules(ctx, "user-1", attempt); !errors.Is(err, model.ErrRuleOrderMismatch) {
				t.Errorf("ReorderRules() error = %v, want ErrRuleOrderMismatch", err)
			}
			got, _ := d.GetRules(ctx, "user-1")
			if names := ruleNames(got); !slices.Equal(names, []string{"c", "a", "b"}) {
				t.Errorf("a rejected reorder changed the order to %v", names)
			}
		})
	}

	// A rule created after a reorder is still appended at the bottom
	appended := createRules(t, d, "user-1", "d")[0]
	if appended.Order != 3 {
		t.Errorf("appended order = %d, want 3", appended.Order)
	}
}

func TestRules_IncrementApplied(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()
	rules := createRules(t, d, "user-1", "a", "b", "c")

	if err := d.IncrementRulesApplied(ctx, []string{rules[0].ID, rules[2].ID}); err != nil {
		t.Fatalf("IncrementRulesApplied() error = %v", err)
	}
	if err := d.IncrementRulesApplied(ctx, []string{rules[0].ID}); err != nil {
		t.Fatalf("IncrementRulesApplied() error = %v", err)
	}

	got, _ := d.GetRules(ctx, "user-1")
	want := []int{2, 0, 1}
	for i, r := range got {
		if r.Applied != want[i] {
			t.Errorf("rule %s applied = %d, want %d", r.Name, r.Applied, want[i])
		}
		if (r.LastAppliedAt != nil) != (want[i] > 0) {
			t.Errorf("rule %s last_applied_at = %v", r.Name, r.LastAppliedAt)
		}
	}
}

func TestRules_Delete(t *testing.T) {
	d := newRuleTestDB(t)
	ctx := context.Background()
	rules := createRules(t, d, "user-1", "a", "b")
	createRules(t, d, "user-2", "x")

	// Another user's ID deletes nothing
	if err := d.DeleteRule(ctx, rules[0].ID, "user-2"); err != nil {
		t.Fatalf("DeleteRule() error = %v", err)
	}
	if count, _ := d.GetRuleCount(ctx, "user-1"); count != 2 {
		t.Errorf("count after another user's delete = %d, want 2", count)
	}

	if err := d.DeleteRule(ctx, rules[0].ID, "user-1"); err != nil {
		t.Fatalf("DeleteRule() error = %v", err)
	}
	got, _ := d.GetRules(ctx, "user-1")
	if names := ruleNames(got); !slices.Equal(names, []string{"b"}) {
		t.Errorf("rules after delete = %v", names)
	}

	if err := d.DeleteRulesByUserID(ctx, "user-1"); err != nil {
		t.Fatalf("DeleteRulesByUserID() error = %v", err)
	}
	if count, _ := d.GetRuleCount(ctx, "user-1"); count != 0 {
		t.Errorf("user-1 count after account deletion = %d, want 0", count)
	}
	if count, _ := d.GetRuleCount(ctx, "user-2"); count != 1 {
		t.Errorf("user-2 count = %d, want 1", count)
	}
}

func TestUpdateAliasDescription(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()

	alias := model.Alias{Name: "a@mailx.net", UserID: "user-1", Enabled: true, Description: "old", Recipients: "r@example.com"}
	if err := d.Client.Create(&alias).Error; err != nil {
		t.Fatalf("create alias error = %v", err)
	}

	if err := d.UpdateAliasDescription(ctx, alias.ID, "user-2", "hijacked"); err != nil {
		t.Fatalf("UpdateAliasDescription() error = %v", err)
	}
	if err := d.UpdateAliasDescription(ctx, alias.ID, "user-1", ""); err != nil {
		t.Fatalf("UpdateAliasDescription() error = %v", err)
	}

	var got model.Alias
	d.Client.First(&got, "id = ?", alias.ID)
	if got.Description != "" || got.Recipients != "r@example.com" || !got.Enabled {
		t.Errorf("alias = %+v, want only the description cleared", got)
	}
}

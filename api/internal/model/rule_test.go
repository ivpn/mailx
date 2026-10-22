package model

import (
	"strings"
	"testing"
)

func textCondition(t RuleConditionType, match RuleMatch, values ...string) RuleCondition {
	return RuleCondition{Type: t, Match: match, Values: values}
}

func boolCondition(t RuleConditionType) RuleCondition {
	return RuleCondition{Type: t, Values: []string{}}
}

func matchRule(op RuleOperator, conditions ...RuleCondition) Rule {
	return Rule{Operator: op, Conditions: conditions}
}

func TestRuleMatches_TextMatchTypes(t *testing.T) {
	in := RuleInput{Subject: "Your Invoice #42 is ready"}

	tests := []struct {
		match  RuleMatch
		values []string
		want   bool
	}{
		{RuleMatchIs, []string{"your invoice #42 is ready"}, true},
		{RuleMatchIs, []string{"invoice"}, false},
		{RuleMatchIsNot, []string{"invoice"}, true},
		{RuleMatchIsNot, []string{"YOUR INVOICE #42 IS READY"}, false},
		{RuleMatchContains, []string{"invoice"}, true},
		{RuleMatchContains, []string{"receipt"}, false},
		{RuleMatchNotContains, []string{"receipt"}, true},
		{RuleMatchNotContains, []string{"invoice"}, false},
		{RuleMatchStartsWith, []string{"your"}, true},
		{RuleMatchStartsWith, []string{"invoice"}, false},
		{RuleMatchNotStartsWith, []string{"invoice"}, true},
		{RuleMatchNotStartsWith, []string{"your"}, false},
		{RuleMatchEndsWith, []string{"ready"}, true},
		{RuleMatchEndsWith, []string{"your"}, false},
		{RuleMatchNotEndsWith, []string{"your"}, true},
		{RuleMatchNotEndsWith, []string{"ready"}, false},
		{RuleMatchRegex, []string{`invoice #\d+`}, true},
		{RuleMatchRegex, []string{`^invoice`}, false},
		{RuleMatchNotRegex, []string{`^invoice`}, true},
		{RuleMatchNotRegex, []string{`#\d+`}, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.match)+" "+strings.Join(tt.values, ","), func(t *testing.T) {
			rule := matchRule(RuleAnd, textCondition(RuleConditionSubject, tt.match, tt.values...))
			if got := rule.Matches(in); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleMatches_CaseInsensitiveAndTrimmed(t *testing.T) {
	in := RuleInput{Sender: "  News@Substack.COM ", Subject: "BIG Sale"}

	tests := []struct {
		name string
		cond RuleCondition
	}{
		{"is", textCondition(RuleConditionSender, RuleMatchIs, " news@SUBSTACK.com ")},
		{"ends with", textCondition(RuleConditionSender, RuleMatchEndsWith, "@substack.com")},
		{"contains", textCondition(RuleConditionSubject, RuleMatchContains, "big sale")},
		{"regex gets (?i)", textCondition(RuleConditionSubject, RuleMatchRegex, "^big sale$")},
		{"uppercase regex", textCondition(RuleConditionSubject, RuleMatchRegex, "SALE")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !matchRule(RuleAnd, tt.cond).Matches(in) {
				t.Error("Matches() = false, want true")
			}
		})
	}
}

func TestRuleMatches_ValuesAreORed(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		match   RuleMatch
		want    bool
	}{
		{"positive: one of the values matches", "Summer promo", RuleMatchContains, true},
		{"positive: no value matches", "Hello", RuleMatchContains, false},
		{"negative: no value matches", "Hello", RuleMatchNotContains, true},
		{"negative: one value matches", "Summer promo", RuleMatchNotContains, false},
		{"negative: every value matches", "Sale and promo", RuleMatchNotContains, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := matchRule(RuleAnd, textCondition(RuleConditionSubject, tt.match, "sale", "promo"))
			if got := rule.Matches(RuleInput{Subject: tt.subject}); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleMatches_EmptyField(t *testing.T) {
	// No sender name: positive matches are false and negative ones true, even for patterns that
	// would match an empty string
	for match, m := range ruleMatchPositive {
		t.Run(string(match), func(t *testing.T) {
			rule := matchRule(RuleAnd, textCondition(RuleConditionSenderName, match, ".*"))
			if got := rule.Matches(RuleInput{SenderName: "  "}); got != m.negated {
				t.Errorf("Matches() = %v, want %v", got, m.negated)
			}
		})
	}
}

func TestRuleMatches_Fields(t *testing.T) {
	in := RuleInput{
		Sender:           "jane@mail.example.com",
		SenderName:       "Jane Doe",
		Subject:          "Hello",
		Alias:            "shop.alias@mailx.net",
		AliasDescription: "Online shopping",
		Destination:      "support@vendor.com",
	}

	tests := []struct {
		name string
		cond RuleCondition
		want bool
	}{
		{"sender", textCondition(RuleConditionSender, RuleMatchIs, "jane@mail.example.com"), true},
		{"sender domain", textCondition(RuleConditionSenderDomain, RuleMatchIs, "mail.example.com"), true},
		{"sender domain is matches subdomains", textCondition(RuleConditionSenderDomain, RuleMatchIs, "example.com"), true},
		{"sender domain with @", textCondition(RuleConditionSenderDomain, RuleMatchIs, "@example.com"), true},
		{"sender domain lookalike", textCondition(RuleConditionSenderDomain, RuleMatchIs, "ample.com"), false},
		{"sender domain is not excludes subdomains", textCondition(RuleConditionSenderDomain, RuleMatchIsNot, "example.com"), false},
		{"sender domain other match types don't", textCondition(RuleConditionSenderDomain, RuleMatchStartsWith, "example"), false},
		{"sender name", textCondition(RuleConditionSenderName, RuleMatchStartsWith, "jane"), true},
		{"subject", textCondition(RuleConditionSubject, RuleMatchIs, "hello"), true},
		{"alias", textCondition(RuleConditionAlias, RuleMatchStartsWith, "shop."), true},
		{"alias description", textCondition(RuleConditionAliasDescription, RuleMatchContains, "shopping"), true},
		{"destination", textCondition(RuleConditionDestination, RuleMatchEndsWith, "@vendor.com"), true},
		{"unknown condition", textCondition("body", RuleMatchContains, "x"), false},
		{"unknown match", textCondition(RuleConditionSubject, "like", "hello"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchRule(RuleAnd, tt.cond).Matches(in); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleMatches_SenderDomainWithoutAt(t *testing.T) {
	rule := matchRule(RuleAnd, textCondition(RuleConditionSenderDomain, RuleMatchIsNot, "example.com"))
	if !rule.Matches(RuleInput{Sender: "not-an-address"}) {
		t.Error("an address without a domain has an empty domain, so a negative match should be true")
	}
}

func TestRuleMatches_BooleanConditions(t *testing.T) {
	attachmentsCalled := 0
	in := RuleInput{
		MailingList:   true,
		AliasWildcard: false,
		HasAttachments: func() bool {
			attachmentsCalled++
			return true
		},
	}

	tests := []struct {
		cond RuleConditionType
		want bool
	}{
		{RuleConditionIsMailingList, true},
		{RuleConditionHasAttachments, true},
		{RuleConditionAliasIsWildcard, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.cond), func(t *testing.T) {
			if got := matchRule(RuleAnd, boolCondition(tt.cond)).Matches(in); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}

	if attachmentsCalled != 1 {
		t.Errorf("HasAttachments called %d times, want 1", attachmentsCalled)
	}

	if matchRule(RuleAnd, boolCondition(RuleConditionHasAttachments)).Matches(RuleInput{}) {
		t.Error("a nil HasAttachments should count as no attachments")
	}
}

func TestRuleMatches_AttachmentsCheckedOnlyWhenNeeded(t *testing.T) {
	in := RuleInput{Subject: "Hello", HasAttachments: func() bool {
		t.Fatal("HasAttachments called although the AND rule already failed")
		return false
	}}

	rule := matchRule(RuleAnd, textCondition(RuleConditionSubject, RuleMatchIs, "other"), boolCondition(RuleConditionHasAttachments))
	if rule.Matches(in) {
		t.Error("Matches() = true, want false")
	}
}

func TestRuleMatches_Operator(t *testing.T) {
	yes := textCondition(RuleConditionSubject, RuleMatchContains, "hello")
	no := textCondition(RuleConditionSubject, RuleMatchContains, "bye")
	in := RuleInput{Subject: "Hello"}

	tests := []struct {
		name string
		rule Rule
		want bool
	}{
		{"AND all match", matchRule(RuleAnd, yes, yes), true},
		{"AND one fails", matchRule(RuleAnd, yes, no), false},
		{"OR one matches", matchRule(RuleOr, no, yes), true},
		{"OR none match", matchRule(RuleOr, no, no), false},
		{"no conditions", matchRule(RuleAnd), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Matches(in); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleMatches_InvalidStoredRegexNeverMatches(t *testing.T) {
	for _, match := range []RuleMatch{RuleMatchRegex, RuleMatchNotRegex} {
		rule := matchRule(RuleAnd, textCondition(RuleConditionSubject, match, "(unclosed"))
		got := rule.Matches(RuleInput{Subject: "(unclosed"})
		if want := match == RuleMatchNotRegex; got != want {
			t.Errorf("%s: Matches() = %v, want %v", match, got, want)
		}
	}
}

func validRule() Rule {
	return Rule{
		Name:       "Newsletters",
		Operator:   RuleOr,
		Conditions: []RuleCondition{boolCondition(RuleConditionIsMailingList), textCondition(RuleConditionSender, RuleMatchEndsWith, "@substack.com")},
		Actions:    []RuleAction{{Type: RuleActionReplaceSubject, Value: "[News] {subject}"}, {Type: RuleActionHideBanner}},
		Forwards:   true,
	}
}

func TestRuleValidate(t *testing.T) {
	tooManyValues := make([]string, MaxRuleConditionValues+1)
	for i := range tooManyValues {
		tooManyValues[i] = strings.Repeat("v", i+1)
	}

	tests := []struct {
		name   string
		modify func(r *Rule)
		err    string
	}{
		{"valid", func(r *Rule) {}, ""},
		{"empty name", func(r *Rule) { r.Name = "  " }, "Name is required."},
		{"long name", func(r *Rule) { r.Name = strings.Repeat("ž", MaxRuleNameLength+1) }, "Name must be at most 50 characters."},
		{"name at limit", func(r *Rule) { r.Name = strings.Repeat("ž", MaxRuleNameLength) }, ""},
		{"name with newline", func(r *Rule) { r.Name = "a\r\nX-Injected: 1" }, "Name can't contain line breaks or control characters."},
		{"bad operator", func(r *Rule) { r.Operator = "XOR" }, "Conditions must be joined with AND or OR."},
		{"lowercase operator", func(r *Rule) { r.Operator = "and" }, ""},
		{"no email types", func(r *Rule) { r.Forwards = false }, "Choose at least one email type to apply the rule to."},
		{"no conditions", func(r *Rule) { r.Conditions = nil }, "Add at least one condition."},
		{"too many conditions", func(r *Rule) {
			r.Conditions = make([]RuleCondition, MaxRuleConditions+1)
			for i := range r.Conditions {
				r.Conditions[i] = boolCondition(RuleConditionIsMailingList)
			}
		}, "A rule can have at most 5 conditions."},
		{"unknown condition", func(r *Rule) { r.Conditions[1].Type = "body" }, `Unknown condition "body".`},
		{"missing match", func(r *Rule) { r.Conditions[1].Match = "" }, "Choose how Sender email should match."},
		{"no values", func(r *Rule) { r.Conditions[1].Values = nil }, "Sender email needs at least one value."},
		{"empty value", func(r *Rule) { r.Conditions[1].Values = []string{"a", " "} }, "Sender email values can't be empty."},
		{"too many values", func(r *Rule) { r.Conditions[1].Values = tooManyValues }, "Sender email can have at most 20 values."},
		{"long value", func(r *Rule) { r.Conditions[1].Values = []string{strings.Repeat("a", MaxRuleConditionValueLength+1)} }, "Sender email values must be at most 255 characters."},
		{"duplicate values ignore case", func(r *Rule) { r.Conditions[1].Values = []string{"Sale", " sale"} }, `Sender email has the value "sale" more than once.`},
		{"regex values keep case", func(r *Rule) {
			r.Conditions[1] = textCondition(RuleConditionSubject, RuleMatchRegex, `\d`, `\D`)
		}, ""},
		{"invalid regex", func(r *Rule) {
			r.Conditions[1] = textCondition(RuleConditionSubject, RuleMatchNotRegex, "(unclosed")
		}, `Invalid regex "(unclosed": error parsing regexp: missing closing ): ` + "`(unclosed`"},
		{"boolean values are dropped", func(r *Rule) { r.Conditions[0].Values = []string{"x"}; r.Conditions[0].Match = "is" }, ""},
		{"no actions", func(r *Rule) { r.Actions = nil }, "Add at least one action."},
		{"too many actions", func(r *Rule) {
			for range MaxRuleActions {
				r.Actions = append(r.Actions, RuleAction{Type: RuleActionForwardTo, Value: "id"})
			}
		}, "A rule can have at most 5 actions."},
		{"unknown action", func(r *Rule) { r.Actions[1].Type = "delete_alias" }, `Unknown action "delete_alias".`},
		{"duplicate action", func(r *Rule) { r.Actions = append(r.Actions, RuleAction{Type: RuleActionHideBanner}) }, "Hide the Mailx banner can only be used once in a rule."},
		{"forward_to more than once", func(r *Rule) {
			r.Actions = []RuleAction{{Type: RuleActionForwardTo, Value: "a"}, {Type: RuleActionForwardTo, Value: "b"}}
		}, ""},
		{"forward_to without recipients", func(r *Rule) { r.Actions = []RuleAction{{Type: RuleActionForwardTo, Value: " , "}} }, "Forward to recipient(s) needs at least one recipient."},
		{"subject required", func(r *Rule) { r.Actions[0].Value = "" }, "Replace the subject with needs a value."},
		{"subject too long", func(r *Rule) { r.Actions[0].Value = strings.Repeat("a", MaxRuleActionValueLength+1) }, "Replace the subject with must be at most 100 characters."},
		{"subject with CRLF", func(r *Rule) { r.Actions[0].Value = "a\r\nBcc: x@example.com" }, "Replace the subject with can't contain line breaks or control characters."},
		{"display name with LF", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionSetDisplayName, Value: "a\nb"} }, "Set the sender name to can't contain line breaks or control characters."},
		{"description with CR", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionSetAliasDescription, Value: "a\rb"} }, "Set the alias description to can't contain line breaks or control characters."},
		{"name with line separator", func(r *Rule) { r.Name = "a b" }, "Name can't contain line breaks or control characters."},
		{"subject with line separator", func(r *Rule) { r.Actions[0].Value = "a b" }, "Replace the subject with can't contain line breaks or control characters."},
		{"display name with paragraph separator", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionSetDisplayName, Value: "a b"} }, "Set the sender name to can't contain line breaks or control characters."},
		{"description with line separator", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionSetAliasDescription, Value: "a b"} }, "Set the alias description to can't contain line breaks or control characters."},
		{"subject with non-ASCII text", func(r *Rule) { r.Actions[0].Value = "[Newsletter] Žltá – {subject}" }, ""},
		{"empty description clears it", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionSetAliasDescription} }, ""},
		{"hide and show banner", func(r *Rule) { r.Actions[0] = RuleAction{Type: RuleActionShowBanner} }, "Hide the Mailx banner and Show the Mailx banner can't be used together."},
		{"block with blocklist", func(r *Rule) {
			r.Actions = []RuleAction{{Type: RuleActionBlocklistDomain}, {Type: RuleActionBlock}, {Type: RuleActionDeactivateAlias}, {Type: RuleActionStop}}
		}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validRule()
			tt.modify(&r)
			r.Normalize()

			err := r.Validate()
			if tt.err == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil || err.Error() != tt.err {
				t.Fatalf("Validate() error = %v, want %q", err, tt.err)
			}
		})
	}
}

func TestRuleValidate_RE2RejectsLookarounds(t *testing.T) {
	// JS accepts lookarounds, so the client-side check lets them through; the server is authoritative
	r := validRule()
	r.Conditions[1] = textCondition(RuleConditionSubject, RuleMatchRegex, "(?<=x)y")
	r.Normalize()

	err := r.Validate()
	if err == nil || !strings.HasPrefix(err.Error(), `Invalid regex "(?<=x)y": error parsing regexp: `) {
		t.Errorf("Validate() error = %v, want an invalid regex error", err)
	}
}

func TestRuleValidate_BlockConflicts(t *testing.T) {
	for _, conflict := range ruleBlockConflicts {
		t.Run(string(conflict), func(t *testing.T) {
			r := validRule()
			r.Actions = []RuleAction{{Type: RuleActionBlock}, {Type: conflict, Value: "x"}}
			r.Normalize()

			err := r.Validate()
			want := "Block the email can't be combined with " + ruleActionSpecs[conflict].Label + ", because the email isn't delivered."
			if err == nil || err.Error() != want {
				t.Errorf("Validate() error = %v, want %q", err, want)
			}
		})
	}
}

func TestRuleValidate_EmailTypes(t *testing.T) {
	tests := []struct {
		name                     string
		forwards, replies, sends bool
		cond                     RuleCondition
		action                   RuleAction
		err                      string
	}{
		{"blocklist sender on replies", false, true, false, boolCondition(RuleConditionHasAttachments), RuleAction{Type: RuleActionBlocklistSender}, "Add sender email to the blocklist only applies to Forwards."},
		{"blocklist sender on forwards and sends", true, false, true, boolCondition(RuleConditionHasAttachments), RuleAction{Type: RuleActionBlocklistSender}, "Add sender email to the blocklist only applies to Forwards."},
		{"mailing list on sends", false, false, true, boolCondition(RuleConditionIsMailingList), RuleAction{Type: RuleActionStop}, "Email is from a mailing list / newsletter only applies to Forwards."},
		{"sender name on replies", false, true, false, textCondition(RuleConditionSenderName, RuleMatchIs, "x"), RuleAction{Type: RuleActionStop}, "Sender name only applies to Forwards."},
		{"destination on forwards", true, true, false, textCondition(RuleConditionDestination, RuleMatchIs, "x"), RuleAction{Type: RuleActionStop}, "Destination email only applies to Replies and Sends."},
		{"destination on replies and sends", false, true, true, textCondition(RuleConditionDestination, RuleMatchIs, "x"), RuleAction{Type: RuleActionBlock}, ""},
		{"subject on all types", true, true, true, textCondition(RuleConditionSubject, RuleMatchIs, "x"), RuleAction{Type: RuleActionReplaceSubject, Value: "y"}, ""},
		{"forward_to on sends", false, false, true, boolCondition(RuleConditionHasAttachments), RuleAction{Type: RuleActionForwardTo, Value: "id"}, "Forward to recipient(s) only applies to Forwards."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{Name: "r", Operator: RuleAnd, Forwards: tt.forwards, Replies: tt.replies, Sends: tt.sends,
				Conditions: []RuleCondition{tt.cond}, Actions: []RuleAction{tt.action}}
			r.Normalize()

			err := r.Validate()
			if tt.err == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil || err.Error() != tt.err {
				t.Fatalf("Validate() error = %v, want %q", err, tt.err)
			}
		})
	}
}

func TestRuleNormalize(t *testing.T) {
	r := Rule{
		Name:     "  Shops  ",
		Operator: " or ",
		Conditions: []RuleCondition{
			{Type: RuleConditionSubject, Match: RuleMatchContains, Values: []string{" sale ", "promo"}},
			{Type: RuleConditionHasAttachments, Match: RuleMatchIs, Values: []string{"x"}},
		},
		Actions: []RuleAction{
			{Type: RuleActionBlock, Value: "ignored"},
			{Type: RuleActionForwardTo, Value: " a, b ,a,,c "},
			{Type: RuleActionReplaceSubject, Value: "  [Shop] {subject}  "},
		},
	}
	r.Normalize()

	if r.Name != "Shops" || r.Operator != RuleOr {
		t.Errorf("name/operator = %q/%q", r.Name, r.Operator)
	}
	if got := strings.Join(r.Conditions[0].Values, "|"); got != "sale|promo" {
		t.Errorf("text values = %q", got)
	}
	if r.Conditions[1].Match != "" || r.Conditions[1].Values == nil || len(r.Conditions[1].Values) != 0 {
		t.Errorf("boolean condition = %+v, want no match and empty values", r.Conditions[1])
	}
	if r.Actions[0].Value != "" {
		t.Errorf("flag action value = %q, want empty", r.Actions[0].Value)
	}
	if r.Actions[1].Value != "a,b,c" {
		t.Errorf("forward_to value = %q, want a,b,c", r.Actions[1].Value)
	}
	if r.Actions[2].Value != "[Shop] {subject}" {
		t.Errorf("text action value = %q", r.Actions[2].Value)
	}
}

func TestRuleAppliesTo(t *testing.T) {
	r := Rule{Forwards: true, Sends: true}
	if !r.AppliesTo(Forward) || r.AppliesTo(Reply) || !r.AppliesTo(Send) || r.AppliesTo(Block) {
		t.Errorf("AppliesTo() wrong for %+v", r)
	}
}

func subjectRule(id string, subject string, extra ...RuleAction) Rule {
	return Rule{
		BaseModel:  BaseModel{ID: id},
		Name:       "rule " + id,
		Operator:   RuleAnd,
		Conditions: []RuleCondition{textCondition(RuleConditionSubject, RuleMatchContains, subject)},
		Actions:    extra,
	}
}

func TestEvaluateRules(t *testing.T) {
	in := RuleInput{Subject: "Original"}

	t.Run("no match", func(t *testing.T) {
		o := EvaluateRules([]Rule{subjectRule("a", "other", RuleAction{Type: RuleActionBlock})}, in)
		if o.Matched() || o.Block {
			t.Errorf("outcome = %+v, want nothing", o)
		}
	})

	t.Run("matching rules merge top to bottom", func(t *testing.T) {
		rules := []Rule{
			subjectRule("a", "orig",
				RuleAction{Type: RuleActionReplaceSubject, Value: "[A] {subject}"},
				RuleAction{Type: RuleActionSetDisplayName, Value: "First"},
				RuleAction{Type: RuleActionForwardTo, Value: "r1,r2"},
				RuleAction{Type: RuleActionHideBanner},
			),
			subjectRule("skip", "nope", RuleAction{Type: RuleActionBlock}),
			subjectRule("b", "orig",
				// Conditions see the original subject, not "[A] Original"
				RuleAction{Type: RuleActionReplaceSubject, Value: "[B] {subject}"},
				RuleAction{Type: RuleActionSetDisplayName, Value: "Second"},
				RuleAction{Type: RuleActionForwardTo, Value: "r2,r3"},
				RuleAction{Type: RuleActionShowBanner},
				RuleAction{Type: RuleActionSetAliasDescription, Value: "Seen"},
			),
		}

		o := EvaluateRules(rules, in)

		if strings.Join(o.MatchedIDs, ",") != "a,b" || strings.Join(o.MatchedNames, ",") != "rule a,rule b" {
			t.Errorf("matched = %v %v", o.MatchedIDs, o.MatchedNames)
		}
		if got := o.Subject("Original"); got != "[B] [A] Original" {
			t.Errorf("Subject() = %q, want stacked prefixes", got)
		}
		if o.FromName != "Second" {
			t.Errorf("FromName = %q, want the last match", o.FromName)
		}
		if strings.Join(o.ForwardTo, ",") != "r1,r2,r3" {
			t.Errorf("ForwardTo = %v, want merged and deduplicated", o.ForwardTo)
		}
		if o.RemoveHeader == nil || *o.RemoveHeader {
			t.Errorf("RemoveHeader = %v, want the last match (show)", o.RemoveHeader)
		}
		if o.AliasDescription == nil || *o.AliasDescription != "Seen" {
			t.Errorf("AliasDescription = %v", o.AliasDescription)
		}
		if o.Block {
			t.Error("Block set by a rule that didn't match")
		}
	})

	t.Run("stop ends evaluation after the current rule", func(t *testing.T) {
		rules := []Rule{
			subjectRule("a", "orig", RuleAction{Type: RuleActionReplaceSubject, Value: "[A] {subject}"}, RuleAction{Type: RuleActionStop}),
			subjectRule("b", "orig", RuleAction{Type: RuleActionBlock}),
		}

		o := EvaluateRules(rules, in)
		if strings.Join(o.MatchedIDs, ",") != "a" || o.Block {
			t.Errorf("outcome = %+v, want only rule a", o)
		}
		if got := o.Subject("Original"); got != "[A] Original" {
			t.Errorf("Subject() = %q", got)
		}
	})

	t.Run("block keeps the first blocking rule and side effects of later ones", func(t *testing.T) {
		rules := []Rule{
			subjectRule("a", "orig", RuleAction{Type: RuleActionBlock}),
			subjectRule("b", "orig", RuleAction{Type: RuleActionBlock}, RuleAction{Type: RuleActionBlocklistDomain}, RuleAction{Type: RuleActionBlocklistSender}),
			subjectRule("c", "orig", RuleAction{Type: RuleActionDeactivateAlias}, RuleAction{Type: RuleActionDisableEncryption}, RuleAction{Type: RuleActionRemoveAttachments}),
		}

		o := EvaluateRules(rules, in)
		if !o.Block || o.BlockedBy != "rule a" {
			t.Errorf("Block = %v by %q, want rule a", o.Block, o.BlockedBy)
		}
		if !o.BlocklistDomain || !o.BlocklistSender || !o.DeactivateAlias || !o.DisableEncryption || !o.RemoveAttachments {
			t.Errorf("outcome = %+v, want every flag set", o)
		}
		if len(o.MatchedIDs) != 3 {
			t.Errorf("matched = %v, want all three counted", o.MatchedIDs)
		}
	})
}

func TestRuleOutcomeSubject(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		in     string
		want   string
	}{
		{"unchanged", nil, "Hello", "Hello"},
		{"replace", []string{"Fixed"}, "Hello", "Fixed"},
		{"prefix", []string{"[News] {subject}"}, "Hello", "[News] Hello"},
		{"placeholder twice", []string{"{subject} / {subject}"}, "Hi", "Hi / Hi"},
		{"stacked", []string{"[A] {subject}", "[B] {subject}"}, "Hi", "[B] [A] Hi"},
		{"empty subject", []string{"[News] {subject}"}, "", "[News] "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (RuleOutcome{SubjectValues: tt.values}).Subject(tt.in); got != tt.want {
				t.Errorf("Subject() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRuleOutcomeSubject_CappedAtLineLimit(t *testing.T) {
	// Fifty rules doubling the subject would otherwise grow it exponentially
	values := make([]string, MaxRules)
	for i := range values {
		values[i] = "{subject}{subject}"
	}

	got := (RuleOutcome{SubjectValues: values}).Subject("žluťoučký")
	if n := len([]rune(got)); n != MaxRuleSubjectLength {
		t.Errorf("subject length = %d runes, want %d", n, MaxRuleSubjectLength)
	}
	if !strings.HasPrefix(got, "žluťoučkýžluťoučký") {
		t.Errorf("subject = %q...", got[:40])
	}
}

func TestRuleOutcomeSettings(t *testing.T) {
	settings := Settings{RemoveHeader: false, BlockTrackers: true}

	if got := (RuleOutcome{}).Settings(settings); got != settings {
		t.Errorf("Settings() without override = %+v", got)
	}

	hide := true
	if got := (RuleOutcome{RemoveHeader: &hide}).Settings(settings); !got.RemoveHeader || !got.BlockTrackers {
		t.Errorf("Settings() with hide = %+v", got)
	}

	show := false
	if got := (RuleOutcome{RemoveHeader: &show}).Settings(Settings{RemoveHeader: true}); got.RemoveHeader {
		t.Errorf("Settings() with show = %+v", got)
	}
}

func TestHasAttachments(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{"plain text", "From: a@example.com\r\nSubject: Hi\r\n\r\nHello\r\n", false},
		{"single-part pdf", "From: a@example.com\r\nContent-Type: application/pdf; name=a.pdf\r\n\r\nJVBERi0=\r\n", true},
		{"alternative bodies", multipartEmail("multipart/alternative",
			"Content-Type: text/plain\r\n\r\nHi",
			"Content-Type: text/html\r\n\r\n<p>Hi</p>"), false},
		{"attachment", multipartEmail("multipart/mixed",
			"Content-Type: text/plain\r\n\r\nHi",
			"Content-Type: application/pdf\r\nContent-Disposition: attachment; filename=\"a.pdf\"\r\n\r\nJVBERi0="), true},
		{"text attachment", multipartEmail("multipart/mixed",
			"Content-Type: text/plain\r\n\r\nHi",
			"Content-Type: text/plain\r\nContent-Disposition: ATTACHMENT; filename=\"notes.txt\"\r\n\r\nnotes"), true},
		{"file without disposition in mixed", multipartEmail("multipart/mixed",
			"Content-Type: text/plain\r\n\r\nHi",
			"Content-Type: image/png; name=a.png\r\n\r\niVBORw0="), true},
		{"cid image in related", multipartEmail("multipart/related",
			"Content-Type: text/html\r\n\r\n<img src=\"cid:logo\">",
			"Content-Type: image/png\r\nContent-ID: <logo>\r\n\r\niVBORw0="), false},
		{"inline image", multipartEmail("multipart/mixed",
			"Content-Type: text/html\r\n\r\n<img src=\"cid:logo\">",
			"Content-Type: image/png\r\nContent-Disposition: inline\r\nContent-ID: <logo>\r\n\r\niVBORw0="), false},
		{"nested attachment", multipartEmail("multipart/mixed",
			"Content-Type: multipart/alternative; boundary=\"inner\"\r\n\r\n--inner\r\nContent-Type: text/plain\r\n\r\nHi\r\n--inner--",
			"Content-Type: application/zip\r\nContent-Disposition: attachment; filename=\"a.zip\"\r\n\r\nUEs="), true},
		{"malformed", "not an email", false},
		{"multipart without boundary", "From: a@example.com\r\nContent-Type: multipart/mixed\r\n\r\nx\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAttachments([]byte(tt.data)); got != tt.want {
				t.Errorf("HasAttachments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func multipartEmail(mediaType string, parts ...string) string {
	var b strings.Builder
	b.WriteString("From: a@example.com\r\nMIME-Version: 1.0\r\nContent-Type: " + mediaType + "; boundary=\"outer\"\r\n\r\n")
	for _, part := range parts {
		b.WriteString("--outer\r\n" + part + "\r\n")
	}
	b.WriteString("--outer--\r\n")
	return b.String()
}

func TestParseMsg_MailingList(t *testing.T) {
	tests := []struct {
		name    string
		headers string
		want    bool
	}{
		{"List-Id", "List-Id: News <news.example.com>\r\n", true},
		{"List-Unsubscribe", "List-Unsubscribe: <https://example.com/u>\r\n", true},
		{"neither", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := "From: news@example.com\r\nTo: alias@mailx.net\r\nSubject: Hi\r\n" + tt.headers + "\r\nHello\r\n"
			msg, err := ParseMsg([]byte(data))
			if err != nil {
				t.Fatalf("ParseMsg() error = %v", err)
			}
			if got := msg.IsMailingList(); got != tt.want {
				t.Errorf("IsMailingList() = %v, want %v", got, tt.want)
			}
		})
	}
}

package model

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	// MaxRules is the per-user cap on rules.
	MaxRules                    = 50
	MaxRuleConditions           = 5
	MaxRuleActions              = 5
	MaxRuleConditionValues      = 20
	MaxRuleConditionValueLength = 255
	MaxRuleNameLength           = 50
	MaxRuleActionValueLength    = 100
	// MaxRuleSubjectLength caps a replaced subject at the RFC 5322 line limit.
	MaxRuleSubjectLength = 998

	// RuleSubjectPlaceholder is replaced with the current subject in a replace_subject value.
	RuleSubjectPlaceholder = "{subject}"

	maxMIMEDepth = 10
)

var (
	ErrRuleNotFound      = errors.New("Rule not found.")
	ErrRuleOrderMismatch = errors.New("rule IDs don't match the user's rules")
)

type RuleOperator string

const (
	RuleAnd RuleOperator = "AND"
	RuleOr  RuleOperator = "OR"
)

type RuleConditionType string

const (
	RuleConditionSender           RuleConditionType = "sender"
	RuleConditionSenderDomain     RuleConditionType = "sender_domain"
	RuleConditionSenderName       RuleConditionType = "sender_name"
	RuleConditionSubject          RuleConditionType = "subject"
	RuleConditionAlias            RuleConditionType = "alias"
	RuleConditionAliasDescription RuleConditionType = "alias_description"
	RuleConditionDestination      RuleConditionType = "destination"
	RuleConditionIsMailingList    RuleConditionType = "is_mailing_list"
	RuleConditionHasAttachments   RuleConditionType = "has_attachments"
	RuleConditionAliasIsWildcard  RuleConditionType = "alias_is_wildcard"
)

type RuleMatch string

const (
	RuleMatchIs            RuleMatch = "is"
	RuleMatchIsNot         RuleMatch = "is_not"
	RuleMatchContains      RuleMatch = "contains"
	RuleMatchNotContains   RuleMatch = "not_contains"
	RuleMatchStartsWith    RuleMatch = "starts_with"
	RuleMatchNotStartsWith RuleMatch = "not_starts_with"
	RuleMatchEndsWith      RuleMatch = "ends_with"
	RuleMatchNotEndsWith   RuleMatch = "not_ends_with"
	RuleMatchRegex         RuleMatch = "regex"
	RuleMatchNotRegex      RuleMatch = "not_regex"
)

// ruleMatchPositive maps every match type to its positive form; negated is true for the "not" ones.
var ruleMatchPositive = map[RuleMatch]struct {
	positive RuleMatch
	negated  bool
}{
	RuleMatchIs:            {RuleMatchIs, false},
	RuleMatchIsNot:         {RuleMatchIs, true},
	RuleMatchContains:      {RuleMatchContains, false},
	RuleMatchNotContains:   {RuleMatchContains, true},
	RuleMatchStartsWith:    {RuleMatchStartsWith, false},
	RuleMatchNotStartsWith: {RuleMatchStartsWith, true},
	RuleMatchEndsWith:      {RuleMatchEndsWith, false},
	RuleMatchNotEndsWith:   {RuleMatchEndsWith, true},
	RuleMatchRegex:         {RuleMatchRegex, false},
	RuleMatchNotRegex:      {RuleMatchRegex, true},
}

type RuleActionType string

const (
	RuleActionBlock               RuleActionType = "block"
	RuleActionBlocklistSender     RuleActionType = "blocklist_sender"
	RuleActionBlocklistDomain     RuleActionType = "blocklist_domain"
	RuleActionReplaceSubject      RuleActionType = "replace_subject"
	RuleActionSetDisplayName      RuleActionType = "set_display_name"
	RuleActionForwardTo           RuleActionType = "forward_to"
	RuleActionDisableEncryption   RuleActionType = "disable_encryption"
	RuleActionHideBanner          RuleActionType = "hide_banner"
	RuleActionShowBanner          RuleActionType = "show_banner"
	RuleActionDeactivateAlias     RuleActionType = "deactivate_alias"
	RuleActionSetAliasDescription RuleActionType = "set_alias_description"
	RuleActionRemoveAttachments   RuleActionType = "remove_attachments"
	RuleActionStop                RuleActionType = "stop"
)

// ruleRelayTypes are the email types a condition or action can be used with.
type ruleRelayTypes struct {
	Forwards, Replies, Sends bool
}

var (
	allRelayTypes   = ruleRelayTypes{Forwards: true, Replies: true, Sends: true}
	forwardsOnly    = ruleRelayTypes{Forwards: true}
	repliesAndSends = ruleRelayTypes{Replies: true, Sends: true}
)

type ruleConditionSpec struct {
	Label   string
	Boolean bool
	ruleRelayTypes
}

var ruleConditionSpecs = map[RuleConditionType]ruleConditionSpec{
	RuleConditionSender:           {"Sender email", false, allRelayTypes},
	RuleConditionSenderDomain:     {"Sender domain", false, forwardsOnly},
	RuleConditionSenderName:       {"Sender name", false, forwardsOnly},
	RuleConditionSubject:          {"Subject", false, allRelayTypes},
	RuleConditionAlias:            {"Alias", false, allRelayTypes},
	RuleConditionAliasDescription: {"Alias description", false, allRelayTypes},
	RuleConditionDestination:      {"Destination email", false, repliesAndSends},
	RuleConditionIsMailingList:    {"Email is from a mailing list / newsletter", true, forwardsOnly},
	RuleConditionHasAttachments:   {"Email has attachments", true, allRelayTypes},
	RuleConditionAliasIsWildcard:  {"Alias is a wildcard (catch-all) alias", true, allRelayTypes},
}

type ruleActionValue int

const (
	ruleActionValueNone ruleActionValue = iota
	ruleActionValueText
	ruleActionValueRecipients
)

type ruleActionSpec struct {
	Label string
	Value ruleActionValue
	// Required is set for text actions that make no sense with an empty value.
	Required bool
	ruleRelayTypes
}

var ruleActionSpecs = map[RuleActionType]ruleActionSpec{
	RuleActionBlock:               {"Block the email", ruleActionValueNone, false, allRelayTypes},
	RuleActionBlocklistSender:     {"Add sender email to the blocklist", ruleActionValueNone, false, forwardsOnly},
	RuleActionBlocklistDomain:     {"Add sender domain to the blocklist", ruleActionValueNone, false, forwardsOnly},
	RuleActionReplaceSubject:      {"Replace the subject with", ruleActionValueText, true, allRelayTypes},
	RuleActionSetDisplayName:      {"Set the sender name to", ruleActionValueText, true, allRelayTypes},
	RuleActionForwardTo:           {"Forward to recipient(s)", ruleActionValueRecipients, true, forwardsOnly},
	RuleActionDisableEncryption:   {"Turn off PGP encryption", ruleActionValueNone, false, forwardsOnly},
	RuleActionHideBanner:          {"Hide the Mailx banner", ruleActionValueNone, false, forwardsOnly},
	RuleActionShowBanner:          {"Show the Mailx banner", ruleActionValueNone, false, forwardsOnly},
	RuleActionDeactivateAlias:     {"Deactivate the alias", ruleActionValueNone, false, forwardsOnly},
	RuleActionSetAliasDescription: {"Set the alias description to", ruleActionValueText, false, forwardsOnly},
	RuleActionRemoveAttachments:   {"Remove attachments", ruleActionValueNone, false, forwardsOnly},
	RuleActionStop:                {"Stop processing more rules", ruleActionValueNone, false, allRelayTypes},
}

// ruleBlockConflicts are the actions that have no effect on an email the rule also blocks.
var ruleBlockConflicts = []RuleActionType{
	RuleActionForwardTo,
	RuleActionReplaceSubject,
	RuleActionSetDisplayName,
	RuleActionDisableEncryption,
	RuleActionHideBanner,
	RuleActionShowBanner,
	RuleActionRemoveAttachments,
}

// Rule acts on the emails passing through the user's aliases. Rules are evaluated in Order,
// lowest first, against the email types they apply to.
type Rule struct {
	BaseModel
	UserID        string          `json:"-" gorm:"size:191;index:idx_rules_user_id_order,priority:1"`
	Name          string          `json:"name" gorm:"size:50"`
	Order         int             `json:"order" gorm:"not null;default:0;index:idx_rules_user_id_order,priority:2"`
	Operator      RuleOperator    `json:"operator" gorm:"size:3;default:'AND'"`
	Conditions    []RuleCondition `json:"conditions" gorm:"type:json;serializer:json"`
	Actions       []RuleAction    `json:"actions" gorm:"type:json;serializer:json"`
	Forwards      bool            `json:"forwards" gorm:"not null;default:false"`
	Replies       bool            `json:"replies" gorm:"not null;default:false"`
	Sends         bool            `json:"sends" gorm:"not null;default:false"`
	Active        bool            `json:"active" gorm:"not null;default:true"`
	Applied       int             `json:"applied" gorm:"not null;default:0"`
	LastAppliedAt *time.Time      `json:"last_applied_at"`
	// InvalidRecipients is set when a forward_to action points at a recipient that has since
	// been deleted or deactivated.
	InvalidRecipients bool `json:"invalid_recipients,omitempty" gorm:"-"`
}

type RuleCondition struct {
	Type   RuleConditionType `json:"type"`
	Match  RuleMatch         `json:"match"`
	Values []string          `json:"values"` // ORed together; empty for boolean types
}

type RuleAction struct {
	Type  RuleActionType `json:"type"`
	Value string         `json:"value,omitempty"` // empty for flag actions
}

// Normalize trims the user-entered parts of a rule and clears the ones its condition and action
// types don't use, so a rule is stored the way it's validated and matched.
func (r *Rule) Normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Operator = RuleOperator(strings.ToUpper(strings.TrimSpace(string(r.Operator))))

	for i, c := range r.Conditions {
		spec, ok := ruleConditionSpecs[c.Type]
		values := []string{}
		if !ok || !spec.Boolean {
			for _, v := range c.Values {
				values = append(values, strings.TrimSpace(v))
			}
		} else {
			c.Match = ""
		}
		c.Values = values
		r.Conditions[i] = c
	}

	for i, a := range r.Actions {
		spec, ok := ruleActionSpecs[a.Type]
		switch {
		case ok && spec.Value == ruleActionValueNone:
			a.Value = ""
		case ok && spec.Value == ruleActionValueRecipients:
			a.Value = strings.Join(a.RecipientIDs(), ",")
		default:
			a.Value = strings.TrimSpace(a.Value)
		}
		r.Actions[i] = a
	}
}

// Validate reports the first problem with a normalized rule, as a message for the user. It
// doesn't check that forward_to recipients belong to the user, which needs the database.
func (r Rule) Validate() error {
	if r.Name == "" {
		return errors.New("Name is required.")
	}
	if utf8.RuneCountInString(r.Name) > MaxRuleNameLength {
		return fmt.Errorf("Name must be at most %d characters.", MaxRuleNameLength)
	}
	if hasControlChars(r.Name) {
		return errors.New("Name can't contain line breaks or control characters.")
	}

	if r.Operator != RuleAnd && r.Operator != RuleOr {
		return errors.New("Conditions must be joined with AND or OR.")
	}

	if !r.Forwards && !r.Replies && !r.Sends {
		return errors.New("Choose at least one email type to apply the rule to.")
	}

	if len(r.Conditions) == 0 {
		return errors.New("Add at least one condition.")
	}
	if len(r.Conditions) > MaxRuleConditions {
		return fmt.Errorf("A rule can have at most %d conditions.", MaxRuleConditions)
	}
	for _, c := range r.Conditions {
		if err := r.validateCondition(c); err != nil {
			return err
		}
	}

	if len(r.Actions) == 0 {
		return errors.New("Add at least one action.")
	}
	if len(r.Actions) > MaxRuleActions {
		return fmt.Errorf("A rule can have at most %d actions.", MaxRuleActions)
	}
	seen := map[RuleActionType]bool{}
	for _, a := range r.Actions {
		if err := r.validateAction(a); err != nil {
			return err
		}
		if seen[a.Type] && a.Type != RuleActionForwardTo {
			return fmt.Errorf("%s can only be used once in a rule.", ruleActionSpecs[a.Type].Label)
		}
		seen[a.Type] = true
	}

	if seen[RuleActionBlock] {
		for _, t := range ruleBlockConflicts {
			if seen[t] {
				return fmt.Errorf("%s can't be combined with %s, because the email isn't delivered.", ruleActionSpecs[RuleActionBlock].Label, ruleActionSpecs[t].Label)
			}
		}
	}
	if seen[RuleActionHideBanner] && seen[RuleActionShowBanner] {
		return fmt.Errorf("%s and %s can't be used together.", ruleActionSpecs[RuleActionHideBanner].Label, ruleActionSpecs[RuleActionShowBanner].Label)
	}

	return nil
}

func (r Rule) validateCondition(c RuleCondition) error {
	spec, ok := ruleConditionSpecs[c.Type]
	if !ok {
		return fmt.Errorf("Unknown condition %q.", c.Type)
	}
	if err := r.checkRelayTypes(spec.Label, spec.ruleRelayTypes); err != nil {
		return err
	}
	if spec.Boolean {
		return nil
	}

	if _, ok := ruleMatchPositive[c.Match]; !ok {
		return fmt.Errorf("Choose how %s should match.", spec.Label)
	}
	if len(c.Values) == 0 {
		return fmt.Errorf("%s needs at least one value.", spec.Label)
	}
	if len(c.Values) > MaxRuleConditionValues {
		return fmt.Errorf("%s can have at most %d values.", spec.Label, MaxRuleConditionValues)
	}

	isRegex := ruleMatchPositive[c.Match].positive == RuleMatchRegex
	seen := map[string]bool{}
	for _, v := range c.Values {
		if v == "" {
			return fmt.Errorf("%s values can't be empty.", spec.Label)
		}
		if utf8.RuneCountInString(v) > MaxRuleConditionValueLength {
			return fmt.Errorf("%s values must be at most %d characters.", spec.Label, MaxRuleConditionValueLength)
		}

		// Text matching ignores case; regex patterns are kept as typed, since e.g. \d and \D differ
		key := v
		if !isRegex {
			key = strings.ToLower(v)
		}
		if seen[key] {
			return fmt.Errorf("%s has the value %q more than once.", spec.Label, v)
		}
		seen[key] = true

		if isRegex {
			if _, err := regexp.Compile(v); err != nil {
				return fmt.Errorf("Invalid regex %q: %s", v, err.Error())
			}
		}
	}

	return nil
}

func (r Rule) validateAction(a RuleAction) error {
	spec, ok := ruleActionSpecs[a.Type]
	if !ok {
		return fmt.Errorf("Unknown action %q.", a.Type)
	}
	if err := r.checkRelayTypes(spec.Label, spec.ruleRelayTypes); err != nil {
		return err
	}

	switch spec.Value {
	case ruleActionValueText:
		if spec.Required && a.Value == "" {
			return fmt.Errorf("%s needs a value.", spec.Label)
		}
		if utf8.RuneCountInString(a.Value) > MaxRuleActionValueLength {
			return fmt.Errorf("%s must be at most %d characters.", spec.Label, MaxRuleActionValueLength)
		}
		if hasControlChars(a.Value) {
			return fmt.Errorf("%s can't contain line breaks or control characters.", spec.Label)
		}
	case ruleActionValueRecipients:
		if len(a.RecipientIDs()) == 0 {
			return fmt.Errorf("%s needs at least one recipient.", spec.Label)
		}
	}

	return nil
}

// checkRelayTypes rejects a condition or action that doesn't support every email type the rule
// applies to.
func (r Rule) checkRelayTypes(label string, supported ruleRelayTypes) error {
	if (r.Forwards && !supported.Forwards) || (r.Replies && !supported.Replies) || (r.Sends && !supported.Sends) {
		return fmt.Errorf("%s only applies to %s.", label, supported.String())
	}

	return nil
}

func (t ruleRelayTypes) String() string {
	var names []string
	if t.Forwards {
		names = append(names, "Forwards")
	}
	if t.Replies {
		names = append(names, "Replies")
	}
	if t.Sends {
		names = append(names, "Sends")
	}

	switch len(names) {
	case 0:
		return "no email types"
	case 1:
		return names[0]
	default:
		return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
	}
}

// hasControlChars reports whether s has a CR, LF or another control character, or a Unicode line
// or paragraph separator (U+2028, U+2029), none of which belong in a header.
func hasControlChars(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return unicode.IsControl(r) || unicode.In(r, unicode.Zl, unicode.Zp)
	}) != -1
}

// RecipientIDs returns the distinct recipient IDs of a forward_to value, in the order given.
func (a RuleAction) RecipientIDs() []string {
	var ids []string
	for id := range strings.SplitSeq(a.Value, ",") {
		id = strings.TrimSpace(id)
		if id != "" && !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}

	return ids
}

// AppliesTo reports whether the rule runs for the given email type.
func (r Rule) AppliesTo(relayType MessageType) bool {
	switch relayType {
	case Forward:
		return r.Forwards
	case Reply:
		return r.Replies
	case Send:
		return r.Sends
	}

	return false
}

// RuleInput is the email a rule's conditions are matched against, as it arrived: actions of
// earlier rules never change it.
type RuleInput struct {
	Sender           string
	SenderName       string
	Subject          string
	Alias            string
	AliasDescription string
	AliasWildcard    bool
	Destination      string
	MailingList      bool
	// HasAttachments is only called for a rule with an attachment condition, so the body of an
	// email isn't parsed unless a rule needs it. A nil func means no attachments.
	HasAttachments func() bool
}

// Matches reports whether the email matches the rule's conditions: all of them with AND, at
// least one with OR. Regexes are compiled once per rule evaluation.
func (r Rule) Matches(in RuleInput) bool {
	regexps := map[string]*regexp.Regexp{}

	for _, c := range r.Conditions {
		matched := c.matches(in, regexps)
		if r.Operator == RuleOr && matched {
			return true
		}
		if r.Operator != RuleOr && !matched {
			return false
		}
	}

	return r.Operator != RuleOr && len(r.Conditions) > 0
}

func (c RuleCondition) matches(in RuleInput, regexps map[string]*regexp.Regexp) bool {
	switch c.Type {
	case RuleConditionIsMailingList:
		return in.MailingList
	case RuleConditionHasAttachments:
		return in.HasAttachments != nil && in.HasAttachments()
	case RuleConditionAliasIsWildcard:
		return in.AliasWildcard
	case RuleConditionSender:
		return c.matchText(in.Sender, regexps)
	case RuleConditionSenderDomain:
		domain := ""
		if at := strings.LastIndex(in.Sender, "@"); at != -1 {
			domain = in.Sender[at+1:]
		}
		return c.matchText(domain, regexps)
	case RuleConditionSenderName:
		return c.matchText(in.SenderName, regexps)
	case RuleConditionSubject:
		return c.matchText(in.Subject, regexps)
	case RuleConditionAlias:
		return c.matchText(in.Alias, regexps)
	case RuleConditionAliasDescription:
		return c.matchText(in.AliasDescription, regexps)
	case RuleConditionDestination:
		return c.matchText(in.Destination, regexps)
	}

	return false
}

// matchText matches a field against the condition's values, ignoring case and surrounding space.
// A positive match is true if any value matches, a negative one only if none does. A missing
// (empty) field never matches positively, so it always matches negatively.
func (c RuleCondition) matchText(field string, regexps map[string]*regexp.Regexp) bool {
	m, ok := ruleMatchPositive[c.Match]
	if !ok {
		return false
	}

	field = strings.ToLower(strings.TrimSpace(field))
	if field == "" {
		return m.negated
	}

	for _, value := range c.Values {
		if c.matchValue(m.positive, field, value, regexps) {
			return !m.negated
		}
	}

	return m.negated
}

func (c RuleCondition) matchValue(match RuleMatch, field string, value string, regexps map[string]*regexp.Regexp) bool {
	if match == RuleMatchRegex {
		re, ok := regexps[value]
		if !ok {
			// A pattern is checked when the rule is saved; one that no longer compiles never matches
			re, _ = regexp.Compile("(?i)" + strings.TrimSpace(value))
			regexps[value] = re
		}
		return re != nil && re.MatchString(field)
	}

	value = strings.ToLower(strings.TrimSpace(value))
	if c.Type == RuleConditionSenderDomain {
		value = strings.TrimPrefix(value, "@")
	}
	if value == "" {
		return false
	}

	switch match {
	case RuleMatchIs:
		// Like the blocklist, a domain also matches its subdomains
		if c.Type == RuleConditionSenderDomain {
			return field == value || strings.HasSuffix(field, "."+value)
		}
		return field == value
	case RuleMatchContains:
		return strings.Contains(field, value)
	case RuleMatchStartsWith:
		return strings.HasPrefix(field, value)
	case RuleMatchEndsWith:
		return strings.HasSuffix(field, value)
	}

	return false
}

// RuleOutcome is what the matching rules do to an email, merged top to bottom.
type RuleOutcome struct {
	// MatchedIDs and MatchedNames list every matching rule in evaluation order.
	MatchedIDs   []string
	MatchedNames []string

	Block bool
	// BlockedBy is the name of the first rule that blocked the email.
	BlockedBy string

	BlocklistSender bool
	BlocklistDomain bool

	// SubjectValues are the replace_subject values in the order they're applied.
	SubjectValues []string
	// FromName is the sender name to use instead of the email's own; the last match wins.
	FromName string
	// ForwardTo are the recipient IDs that replace the alias's recipients, merged across rules.
	ForwardTo []string

	DisableEncryption bool
	RemoveAttachments bool
	// RemoveHeader overrides settings.RemoveHeader when set; the last match wins.
	RemoveHeader *bool

	DeactivateAlias bool
	// AliasDescription is the new alias description when set; the last match wins.
	AliasDescription *string
}

// Matched reports whether any rule matched.
func (o RuleOutcome) Matched() bool {
	return len(o.MatchedIDs) > 0
}

// EvaluateRules applies every rule that matches the email, in the given order, until one with a
// stop action. Rules must already be filtered to the active ones for the email's type.
func EvaluateRules(rules []Rule, in RuleInput) RuleOutcome {
	var o RuleOutcome

	for _, r := range rules {
		if !r.Matches(in) {
			continue
		}

		o.MatchedIDs = append(o.MatchedIDs, r.ID)
		o.MatchedNames = append(o.MatchedNames, r.Name)

		stop := false
		for _, a := range r.Actions {
			switch a.Type {
			case RuleActionBlock:
				if !o.Block {
					o.Block = true
					o.BlockedBy = r.Name
				}
			case RuleActionBlocklistSender:
				o.BlocklistSender = true
			case RuleActionBlocklistDomain:
				o.BlocklistDomain = true
			case RuleActionReplaceSubject:
				o.SubjectValues = append(o.SubjectValues, a.Value)
			case RuleActionSetDisplayName:
				o.FromName = a.Value
			case RuleActionForwardTo:
				for _, id := range a.RecipientIDs() {
					if !slices.Contains(o.ForwardTo, id) {
						o.ForwardTo = append(o.ForwardTo, id)
					}
				}
			case RuleActionDisableEncryption:
				o.DisableEncryption = true
			case RuleActionRemoveAttachments:
				o.RemoveAttachments = true
			case RuleActionHideBanner, RuleActionShowBanner:
				remove := a.Type == RuleActionHideBanner
				o.RemoveHeader = &remove
			case RuleActionDeactivateAlias:
				o.DeactivateAlias = true
			case RuleActionSetAliasDescription:
				description := a.Value
				o.AliasDescription = &description
			case RuleActionStop:
				stop = true
			}
		}

		if stop {
			break
		}
	}

	return o
}

// Subject returns the subject after every replace_subject action, in order. Each value's
// {subject} stands for the subject as it is at that point, so two prefixing rules stack.
func (o RuleOutcome) Subject(subject string) string {
	for _, value := range o.SubjectValues {
		// Capped at every step, so stacked {subject}s can't grow it exponentially
		subject = truncateRunes(strings.ReplaceAll(value, RuleSubjectPlaceholder, subject), MaxRuleSubjectLength)
	}

	return subject
}

// Settings returns a copy of settings with the rules' banner override applied.
func (o RuleOutcome) Settings(settings Settings) Settings {
	if o.RemoveHeader != nil {
		settings.RemoveHeader = *o.RemoveHeader
	}

	return settings
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	return string([]rune(s)[:max])
}

// HasAttachments reports whether a raw email has a part that's sent as an attachment, the same
// way the mailer tells attachments apart from bodies and inline files such as cid: images.
func HasAttachments(data []byte) bool {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return false
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return false
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		return multipartHasAttachments(msg.Body, mediaType, params["boundary"], 0)
	}

	return !isBodyMediaType(mediaType)
}

func multipartHasAttachments(body io.Reader, parentType string, boundary string, depth int) bool {
	if boundary == "" || depth > maxMIMEDepth {
		return false
	}

	mr := multipart.NewReader(body, boundary)
	for {
		part, err := mr.NextRawPart()
		if err != nil {
			return false
		}

		disposition, _, _ := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		if disposition == "attachment" {
			return true
		}

		mediaType, params, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			mediaType = "text/plain"
		}

		switch {
		case isBodyMediaType(mediaType):
			continue
		case strings.HasPrefix(mediaType, "multipart/"):
			if multipartHasAttachments(part, mediaType, params["boundary"], depth+1) {
				return true
			}
		case disposition == "inline" || parentType == "multipart/related":
			continue
		default:
			return true
		}
	}
}

func isBodyMediaType(mediaType string) bool {
	return mediaType == "text/plain" || mediaType == "text/html" || mediaType == "text/enriched"
}

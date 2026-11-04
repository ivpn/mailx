package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetRules                = errors.New("Unable to retrieve rules.")
	ErrPostRule                = errors.New("Unable to create rule. Please try again.")
	ErrPostRuleLimit           = errors.New("You’ve reached the maximum number of rules.")
	ErrPostRuleInactiveSub     = errors.New("Your subscription is not active. Please renew to create new rules.")
	ErrUpdateRule              = errors.New("Unable to update rule. Please try again.")
	ErrUpdateRuleInactiveSub   = errors.New("Your subscription is not active. Please renew to update rules.")
	ErrReorderRules            = errors.New("Unable to reorder rules. Please try again.")
	ErrReorderRulesStale       = errors.New("Your rules have changed since this page was loaded. Please reload and try again.")
	ErrReorderRulesInactiveSub = errors.New("Your subscription is not active. Please renew to reorder rules.")
	ErrDeleteRule              = errors.New("Unable to delete rule. Please try again.")
	ErrDeleteRuleInactiveSub   = errors.New("Your subscription is not active. Please renew to delete rules.")
	ErrDeleteRulesByUserID     = errors.New("Unable to delete rules for this user.")
	ErrRuleRecipient           = errors.New("Forward to recipient(s) can only use your own verified recipients.")
)

type RuleStore interface {
	GetRules(context.Context, string) ([]model.Rule, error)
	GetRule(context.Context, string, string) (model.Rule, error)
	GetActiveRules(context.Context, string, model.MessageType) ([]model.Rule, error)
	GetRuleCount(context.Context, string) (int, error)
	PostRule(context.Context, model.Rule) (model.Rule, error)
	UpdateRule(context.Context, model.Rule) error
	UpdateRuleActive(context.Context, string, string, bool) error
	ReorderRules(context.Context, string, []string) error
	IncrementRulesApplied(context.Context, []string) error
	DeleteRule(context.Context, string, string) error
	DeleteRulesByUserID(context.Context, string) error
}

// GetRules returns the user's rules in evaluation order. A rule whose forward_to points at a
// recipient that's since been deleted or deactivated is flagged, so the UI can warn about it.
func (s *Service) GetRules(ctx context.Context, userID string) ([]model.Rule, error) {
	rules, err := s.Store.GetRules(ctx, userID)
	if err != nil {
		log.Printf("error fetching rules: %s", err.Error())
		return nil, ErrGetRules
	}

	if !hasForwardTo(rules) {
		return rules, nil
	}

	rcps, err := s.Store.GetRecipients(ctx, userID)
	if err != nil {
		log.Printf("error fetching recipients for rules: %s", err.Error())
		return rules, nil
	}

	active := activeRecipientIDs(rcps)
	for i, rule := range rules {
		for _, id := range forwardToIDs(rule) {
			if !active[id] {
				rules[i].InvalidRecipients = true
				break
			}
		}
	}

	return rules, nil
}

func (s *Service) PostRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	if err := s.checkRuleSubscription(ctx, rule.UserID, ErrPostRuleInactiveSub, ErrPostRule); err != nil {
		return model.Rule{}, err
	}

	count, err := s.Store.GetRuleCount(ctx, rule.UserID)
	if err != nil {
		log.Printf("error fetching rule count: %s", err.Error())
		return model.Rule{}, ErrPostRule
	}

	if count >= model.MaxRules {
		return model.Rule{}, ErrPostRuleLimit
	}

	rule, err = s.prepareRule(ctx, rule, ErrPostRule)
	if err != nil {
		return model.Rule{}, err
	}

	rule.ID = ""
	rule.Active = true
	rule.Applied = 0
	rule.LastAppliedAt = nil

	created, err := s.Store.PostRule(ctx, rule)
	if err != nil {
		log.Printf("error creating rule: %s", err.Error())
		return model.Rule{}, ErrPostRule
	}

	return created, nil
}

// UpdateRule replaces a rule's definition. Its order, active flag and counters are kept.
func (s *Service) UpdateRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	if err := s.checkRuleSubscription(ctx, rule.UserID, ErrUpdateRuleInactiveSub, ErrUpdateRule); err != nil {
		return model.Rule{}, err
	}

	existing, err := s.getRule(ctx, rule.ID, rule.UserID, ErrUpdateRule)
	if err != nil {
		return model.Rule{}, err
	}

	rule, err = s.prepareRule(ctx, rule, ErrUpdateRule)
	if err != nil {
		return model.Rule{}, err
	}

	existing.Name = rule.Name
	existing.Operator = rule.Operator
	existing.Conditions = rule.Conditions
	existing.Actions = rule.Actions
	existing.Forwards = rule.Forwards
	existing.Replies = rule.Replies
	existing.Sends = rule.Sends

	err = s.Store.UpdateRule(ctx, existing)
	if err != nil {
		log.Printf("error updating rule: %s", err.Error())
		return model.Rule{}, ErrUpdateRule
	}

	return existing, nil
}

func (s *Service) UpdateRuleActive(ctx context.Context, ID string, userID string, active bool) error {
	if err := s.checkRuleSubscription(ctx, userID, ErrUpdateRuleInactiveSub, ErrUpdateRule); err != nil {
		return err
	}

	if _, err := s.getRule(ctx, ID, userID, ErrUpdateRule); err != nil {
		return err
	}

	err := s.Store.UpdateRuleActive(ctx, ID, userID, active)
	if err != nil {
		log.Printf("error updating rule active: %s", err.Error())
		return ErrUpdateRule
	}

	return nil
}

// ReorderRules sets the evaluation order to the order of ids, which must be exactly the user's
// rule IDs.
func (s *Service) ReorderRules(ctx context.Context, userID string, ids []string) error {
	if err := s.checkRuleSubscription(ctx, userID, ErrReorderRulesInactiveSub, ErrReorderRules); err != nil {
		return err
	}

	err := s.Store.ReorderRules(ctx, userID, ids)
	if errors.Is(err, model.ErrRuleOrderMismatch) {
		return ErrReorderRulesStale
	}
	if err != nil {
		log.Printf("error reordering rules: %s", err.Error())
		return ErrReorderRules
	}

	return nil
}

func (s *Service) DeleteRule(ctx context.Context, ID string, userID string) error {
	if err := s.checkRuleSubscription(ctx, userID, ErrDeleteRuleInactiveSub, ErrDeleteRule); err != nil {
		return err
	}

	if _, err := s.getRule(ctx, ID, userID, ErrDeleteRule); err != nil {
		return err
	}

	err := s.Store.DeleteRule(ctx, ID, userID)
	if err != nil {
		log.Printf("error deleting rule: %s", err.Error())
		return ErrDeleteRule
	}

	return nil
}

func (s *Service) DeleteRulesByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteRulesByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting rules by user ID: %s", err.Error())
		return ErrDeleteRulesByUserID
	}

	return nil
}

// checkRuleSubscription rejects changes to rules while the subscription is inactive. Existing
// rules keep being applied, and toggling one counts as a change.
func (s *Service) checkRuleSubscription(ctx context.Context, userID string, errInactive error, errFallback error) error {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		log.Printf("error fetching subscription: %s", err.Error())
		return errFallback
	}

	if !sub.ActiveStatus() {
		return errInactive
	}

	return nil
}

// getRule returns model.ErrRuleNotFound for a rule that doesn't exist or belongs to another user.
func (s *Service) getRule(ctx context.Context, ID string, userID string, errFallback error) (model.Rule, error) {
	rule, err := s.Store.GetRule(ctx, ID, userID)
	if errors.Is(err, model.ErrRuleNotFound) {
		return model.Rule{}, model.ErrRuleNotFound
	}
	if err != nil {
		log.Printf("error fetching rule: %s", err.Error())
		return model.Rule{}, errFallback
	}

	return rule, nil
}

// prepareRule normalizes and validates a rule from the user, including that its forward_to
// recipients are the user's own verified ones.
func (s *Service) prepareRule(ctx context.Context, rule model.Rule, errFallback error) (model.Rule, error) {
	rule.Normalize()
	if err := rule.Validate(); err != nil {
		return model.Rule{}, err
	}

	ids := forwardToIDs(rule)
	if len(ids) == 0 {
		return rule, nil
	}

	rcps, err := s.Store.GetRecipients(ctx, rule.UserID)
	if err != nil {
		log.Printf("error fetching recipients for rule: %s", err.Error())
		return model.Rule{}, errFallback
	}

	active := activeRecipientIDs(rcps)
	for _, id := range ids {
		if !active[id] {
			return model.Rule{}, ErrRuleRecipient
		}
	}

	return rule, nil
}

func hasForwardTo(rules []model.Rule) bool {
	for _, rule := range rules {
		if len(forwardToIDs(rule)) > 0 {
			return true
		}
	}

	return false
}

func forwardToIDs(rule model.Rule) []string {
	var ids []string
	for _, a := range rule.Actions {
		if a.Type == model.RuleActionForwardTo {
			ids = append(ids, a.RecipientIDs()...)
		}
	}

	return ids
}

// activeRecipientIDs returns the IDs of the verified recipients, the only ones mail is sent to.
func activeRecipientIDs(rcps []model.Recipient) map[string]bool {
	active := make(map[string]bool, len(rcps))
	for _, rcp := range rcps {
		if rcp.IsActive {
			active[rcp.ID] = true
		}
	}

	return active
}

// ApplyRules evaluates the user's active rules for the email type against the original email and
// counts each matching rule once, in a single query. A failed lookup applies no rules, so the
// email is delivered the way it would be without them. destination is the address a reply or send
// goes to, as FindRecipients resolved it.
func (s *Service) ApplyRules(ctx context.Context, alias model.Alias, msg model.Msg, relayType model.MessageType, destination string, data []byte) model.RuleOutcome {
	if alias.UserID == "" {
		return model.RuleOutcome{}
	}

	rules, err := s.Store.GetActiveRules(ctx, alias.UserID, relayType)
	if err != nil {
		log.Printf("error fetching active rules: %s", err.Error())
		return model.RuleOutcome{}
	}
	if len(rules) == 0 {
		return model.RuleOutcome{}
	}

	in := model.RuleInput{
		Sender:           msg.From,
		SenderName:       msg.FromName,
		Subject:          msg.Subject,
		Alias:            alias.Name,
		AliasDescription: alias.Description,
		AliasWildcard:    alias.Wildcard,
		MailingList:      msg.IsMailingList(),
		// Parsed only if a matching rule gets as far as an attachment condition
		HasAttachments: sync.OnceValue(func() bool { return model.HasAttachments(data) }),
	}
	if relayType != model.Forward {
		in.Destination = destination
	}

	outcome := model.EvaluateRules(rules, in)
	if outcome.Matched() {
		if err := s.Store.IncrementRulesApplied(ctx, outcome.MatchedIDs); err != nil {
			log.Printf("error incrementing rules applied count: %s", err.Error())
		}
	}

	return outcome
}

// ApplyRuleSideEffects runs the blocklist and alias actions of the matching rules. They never
// stop the email: failures are logged and delivery carries on.
func (s *Service) ApplyRuleSideEffects(ctx context.Context, alias model.Alias, from string, outcome model.RuleOutcome) {
	email := model.NormalizeBlocklistValue(model.BlocklistEmail, from)

	if outcome.BlocklistSender && blocklistValidator.Var(email, model.BlocklistValueRules(model.BlocklistEmail)) == nil {
		s.addRuleBlocklistItem(ctx, alias.UserID, model.BlocklistEmail, email)
	}

	if at := strings.LastIndex(email, "@"); outcome.BlocklistDomain && at >= 0 {
		domain := model.NormalizeBlocklistValue(model.BlocklistDomain, email[at+1:])
		if blocklistValidator.Var(domain, model.BlocklistValueRules(model.BlocklistDomain)) == nil {
			s.addRuleBlocklistItem(ctx, alias.UserID, model.BlocklistDomain, domain)
		}
	}

	// An on-the-fly alias that couldn't be saved has nothing to change
	if alias.ID == "" {
		return
	}

	if outcome.AliasDescription != nil && *outcome.AliasDescription != alias.Description {
		if err := s.Store.UpdateAliasDescription(ctx, alias.ID, alias.UserID, *outcome.AliasDescription); err != nil {
			log.Printf("error setting alias description by rule: %s", err.Error())
		}
	}

	if outcome.DeactivateAlias && alias.Enabled {
		if err := s.Store.BulkUpdateAliasEnabled(ctx, []string{alias.ID}, alias.UserID, false); err != nil {
			log.Printf("error deactivating alias by rule: %s", err.Error())
		}
	}
}

// addRuleBlocklistItem adds a sender to the blocklist, ignoring an entry that's already there
// and a full blocklist.
func (s *Service) addRuleBlocklistItem(ctx context.Context, userID string, itemType model.BlocklistType, value string) {
	_, err := s.PostBlocklistItem(ctx, model.BlocklistItem{UserID: userID, Type: itemType, Value: value})
	if err != nil && !errors.Is(err, ErrDuplicateBlocklistItem) && !errors.Is(err, ErrPostBlocklistItemLimit) {
		log.Printf("error adding %s to blocklist by rule: %s", itemType, err.Error())
	}
}

// BlockMessageByRule discards an email a rule blocked. A blocked forward is recorded in the
// alias's stats like a blocklist hit; with diagnostics enabled, any blocked email is logged.
func (s *Service) BlockMessageByRule(alias model.Alias, from string, to string, relayType model.MessageType, outcome model.RuleOutcome) {
	ctx := context.Background()

	if relayType == model.Forward {
		if err := s.SaveMessage(ctx, alias, model.Block); err != nil {
			log.Println("error saving message", err)
		}
	}

	settings, err := s.GetSettings(ctx, alias.UserID)
	if err != nil {
		log.Println("error getting settings", err)
		return
	}

	if settings.LogIssues {
		reason := "Blocked by rule \"" + outcome.BlockedBy + "\""
		if err := s.ProcessDiagnosticLog(alias, from, to, reason, model.BlockedByRule); err != nil {
			log.Println("error processing diagnostic log", err)
		}
	}
}

// RuleRecipients returns who a forward goes to: the forward_to recipients of the matching rules
// if there are any, otherwise the alias's own. Recipients that have since been deleted or
// deactivated are skipped, and if none are left the alias's recipients are used.
func (s *Service) RuleRecipients(ctx context.Context, alias model.Alias, recipients []model.Recipient, outcome model.RuleOutcome) []model.Recipient {
	if len(outcome.ForwardTo) == 0 {
		return recipients
	}

	rcps, err := s.Store.GetRecipients(ctx, alias.UserID)
	if err != nil {
		log.Printf("error fetching rule recipients, using the alias's recipients: %s", err.Error())
		return recipients
	}

	byID := make(map[string]model.Recipient, len(rcps))
	for _, rcp := range rcps {
		byID[rcp.ID] = rcp
	}

	var result []model.Recipient
	for _, id := range outcome.ForwardTo {
		if rcp, ok := byID[id]; ok && rcp.IsActive {
			result = append(result, rcp)
		}
	}

	if len(result) == 0 {
		log.Printf("rule recipients for alias %s are no longer available, using the alias's recipients", alias.ID)
		return recipients
	}

	return result
}

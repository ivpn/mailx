package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// ruleOrder is the evaluation order; created_at breaks ties left by concurrent creates.
const ruleOrder = "`order` ASC, created_at ASC"

func (d *Database) GetRules(ctx context.Context, userID string) ([]model.Rule, error) {
	rules := []model.Rule{}
	err := d.Client.WithContext(ctx).Where("user_id = ?", userID).Order(ruleOrder).Find(&rules).Error
	return rules, err
}

// GetRule returns model.ErrRuleNotFound for a rule that doesn't exist or belongs to another user.
func (d *Database) GetRule(ctx context.Context, ID string, userID string) (model.Rule, error) {
	var rule model.Rule
	err := d.Client.WithContext(ctx).Where("id = ? AND user_id = ?", ID, userID).First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Rule{}, model.ErrRuleNotFound
	}
	return rule, err
}

// GetActiveRules returns the user's active rules for an email type in evaluation order, served by
// the (user_id, order) index.
func (d *Database) GetActiveRules(ctx context.Context, userID string, relayType model.MessageType) ([]model.Rule, error) {
	column := ""
	switch relayType {
	case model.Forward:
		column = "forwards"
	case model.Reply:
		column = "replies"
	case model.Send:
		column = "sends"
	default:
		return []model.Rule{}, nil
	}

	rules := []model.Rule{}
	err := d.Client.WithContext(ctx).
		Where("user_id = ? AND active = ? AND "+column+" = ?", userID, true, true).
		Order(ruleOrder).
		Find(&rules).Error
	return rules, err
}

func (d *Database) GetRuleCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.WithContext(ctx).Model(&model.Rule{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

// PostRule appends the rule below the user's existing rules.
func (d *Database) PostRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	err := d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxOrder sql.NullInt64
		err := tx.Model(&model.Rule{}).Where("user_id = ?", rule.UserID).Select("MAX(`order`)").Scan(&maxOrder).Error
		if err != nil {
			return err
		}

		rule.Order = 0
		if maxOrder.Valid {
			rule.Order = int(maxOrder.Int64) + 1
		}

		return tx.Create(&rule).Error
	})
	return rule, err
}

// UpdateRule replaces a rule's definition, leaving its order, active flag and counters as they are.
func (d *Database) UpdateRule(ctx context.Context, rule model.Rule) error {
	return d.Client.WithContext(ctx).Model(&rule).Where("user_id = ?", rule.UserID).
		Select("name", "operator", "conditions", "actions", "forwards", "replies", "sends", "updated_at").
		Updates(&rule).Error
}

func (d *Database) UpdateRuleActive(ctx context.Context, ID string, userID string, active bool) error {
	return d.Client.WithContext(ctx).Model(&model.Rule{}).Where("id = ? AND user_id = ?", ID, userID).
		Update("active", active).Error
}

// ReorderRules sets each rule's order to its index in ids, which must be exactly the user's rule
// IDs, so a list that's gone stale (e.g. a rule created in another tab) is rejected as a whole.
func (d *Database) ReorderRules(ctx context.Context, userID string, ids []string) error {
	return d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing []string
		if err := tx.Model(&model.Rule{}).Where("user_id = ?", userID).Pluck("id", &existing).Error; err != nil {
			return err
		}

		if !sameIDs(existing, ids) {
			return model.ErrRuleOrderMismatch
		}

		for i, id := range ids {
			err := tx.Model(&model.Rule{}).Where("id = ? AND user_id = ?", id, userID).UpdateColumn("order", i).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// sameIDs reports whether ids holds every ID in existing exactly once, and nothing else.
func sameIDs(existing []string, ids []string) bool {
	if len(existing) != len(ids) {
		return false
	}

	set := make(map[string]bool, len(existing))
	for _, id := range existing {
		set[id] = true
	}
	for _, id := range ids {
		if !set[id] {
			return false
		}
		delete(set, id)
	}

	return true
}

// IncrementRulesApplied counts one application of each rule in a single UPDATE, so concurrent
// deliveries matching the same rule don't overwrite each other's increments.
func (d *Database) IncrementRulesApplied(ctx context.Context, ids []string) error {
	return d.Client.WithContext(ctx).Model(&model.Rule{}).Where("id IN ?", ids).UpdateColumns(map[string]any{
		"applied":         gorm.Expr("applied + ?", 1),
		"last_applied_at": time.Now(),
	}).Error
}

func (d *Database) DeleteRule(ctx context.Context, ID string, userID string) error {
	return d.Client.WithContext(ctx).Where("id = ? AND user_id = ?", ID, userID).Delete(&model.Rule{}).Error
}

func (d *Database) DeleteRulesByUserID(ctx context.Context, userID string) error {
	return d.Client.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.Rule{}).Error
}

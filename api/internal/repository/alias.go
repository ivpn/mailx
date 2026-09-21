package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetAlias(ctx context.Context, ID string, userID string) (model.Alias, error) {
	var alias model.Alias
	var aliasStats model.AliasStats
	err := d.Client.Where("id = ? AND user_id = ?", ID, userID).
		First(&alias).Error
	if err != nil {
		return alias, err
	}

	err = d.Client.Model(&model.Message{}).
		Select("SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as forwards, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as blocks, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as replies, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as sends",
			model.Forward, model.Block, model.Reply, model.Send).
		Where("alias_id = ?", ID).
		Scan(&aliasStats).Error
	if err != nil {
		return alias, err
	}

	alias.Stats = aliasStats

	return alias, nil
}

// sanitizeAliasSort independently re-checks sort inputs against the shared allow-list
// (model.AliasSortColumns/AliasSortOrders) so the query builder is safe even if a future
// caller skips the handler-level check, and returns the fully-qualified column to sort by.
// columnPrefix is always a fixed literal supplied by the call site, never request input.
func sanitizeAliasSort(columnPrefix string, sortBy string, sortOrder string) (string, string) {
	if !model.AliasSortColumns[sortBy] {
		sortBy = "created_at"
	}
	if !model.AliasSortOrders[sortOrder] {
		sortOrder = "DESC"
	}
	return columnPrefix + sortBy, sortOrder
}

// aliasSearchFilter builds the wildcard/search WHERE fragment using bound parameters instead
// of concatenating untrusted input into the query text. columnPrefix is always a fixed literal
// ("a." or "") supplied by the call site, never derived from request input. The wildcard flag
// is stored in the catch_all column (kept as-is to avoid a schema migration).
func aliasSearchFilter(columnPrefix string, wildcard string, search string) (string, []any) {
	var filter string
	var args []any

	switch wildcard {
	case "true":
		filter += " AND " + columnPrefix + "catch_all = ?"
		args = append(args, true)
	case "false":
		filter += " AND " + columnPrefix + "catch_all = ?"
		args = append(args, false)
	}

	if search != "" {
		filter += " AND (" + columnPrefix + "name LIKE ? OR " + columnPrefix + "description LIKE ?)"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	return filter, args
}

// aliasStatusFilter builds the deleted_at/enabled WHERE fragment (with bound parameters, so no
// untrusted input reaches the query text) for the given status value. unscoped reports whether
// the caller must bypass GORM's automatic soft-delete scope (only relevant to GetAliasCount's
// query-builder path; GetAliases runs raw SQL, which is never auto-scoped).
func aliasStatusFilter(columnPrefix string, status string) (filter string, args []any, unscoped bool) {
	switch status {
	case "deleted":
		return "AND " + columnPrefix + "deleted_at IS NOT NULL", nil, true
	case "all":
		return "", nil, true
	case "active":
		return "AND " + columnPrefix + "deleted_at IS NULL AND " + columnPrefix + "enabled = ?", []any{true}, false
	case "inactive":
		return "AND " + columnPrefix + "deleted_at IS NULL AND " + columnPrefix + "enabled = ?", []any{false}, false
	default: // "active_inactive" and any unrecognized value fall back to the broadest non-deleted view
		return "AND " + columnPrefix + "deleted_at IS NULL", nil, false
	}
}

// findAliases returns one page of aliases without touching the messages table, so its cost
// scales with the page size rather than with the user's message history.
func (d *Database) findAliases(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	sortBy, sortOrder = sanitizeAliasSort("", sortBy, sortOrder)

	statusFilter, statusArgs, unscoped := aliasStatusFilter("", status)
	filter, filterArgs := aliasSearchFilter("", wildcard, search)

	args := []any{userID}
	args = append(args, statusArgs...)
	args = append(args, filterArgs...)

	q := d.Client.WithContext(ctx).Model(&model.Alias{})
	if unscoped {
		q = q.Unscoped()
	}
	q = q.Where("user_id = ? "+statusFilter+filter, args...).
		Order("pinned DESC").
		Order(sortBy + " " + sortOrder)

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	aliases := []model.Alias{}
	err := q.Find(&aliases).Error
	return aliases, err
}

// attachAliasStats fills in per-alias message counts for a single page of aliases in one
// query served by the (alias_id, type) index, instead of aggregating the whole messages
// table for every alias the user owns.
func (d *Database) attachAliasStats(ctx context.Context, aliases []model.Alias) error {
	if len(aliases) == 0 {
		return nil
	}

	ids := make([]string, 0, len(aliases))
	for i := range aliases {
		ids = append(ids, aliases[i].ID)
	}

	var counts []struct {
		AliasID string
		Type    model.MessageType
		Total   int
	}
	err := d.Client.WithContext(ctx).Model(&model.Message{}).
		Select("alias_id, type, COUNT(*) AS total").
		Where("alias_id IN (?)", ids).
		Group("alias_id").Group("type").
		Scan(&counts).Error
	if err != nil {
		return err
	}

	stats := make(map[string]*model.AliasStats, len(aliases))
	for i := range aliases {
		stats[aliases[i].ID] = &aliases[i].Stats
	}

	for _, c := range counts {
		s, ok := stats[c.AliasID]
		if !ok {
			continue
		}
		switch c.Type {
		case model.Forward:
			s.Forwards = c.Total
		case model.Block:
			s.Blocks = c.Total
		case model.Reply:
			s.Replies = c.Total
		case model.Send:
			s.Sends = c.Total
		}
	}

	return nil
}

func (d *Database) GetAliases(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	aliases, err := d.findAliases(ctx, userID, limit, offset, sortBy, sortOrder, wildcard, search, status)
	if err != nil {
		return nil, err
	}

	if err := d.attachAliasStats(ctx, aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

// GetAliasesNoStats serves callers that only read alias rows, letting them skip the
// per-alias message aggregation entirely.
func (d *Database) GetAliasesNoStats(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	return d.findAliases(ctx, userID, limit, offset, sortBy, sortOrder, wildcard, search, status)
}

func (d *Database) GetAliasesByDomain(ctx context.Context, domain string, userId string) ([]model.Alias, error) {
	aliases := []model.Alias{}
	err := d.Client.Where("name LIKE ? AND user_id = ?", "%@"+domain, userId).Find(&aliases).Error
	return aliases, err
}

func (d *Database) GetAllAliases(ctx context.Context, userID string) ([]model.Alias, error) {
	aliases := []model.Alias{}
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&aliases).Error
	return aliases, err
}

func (d *Database) GetAliasCount(ctx context.Context, userID string, wildcard string, search string, status string) (int, error) {
	statusFilter, statusArgs, unscoped := aliasStatusFilter("", status)
	filter, filterArgs := aliasSearchFilter("", wildcard, search)

	args := []any{userID}
	args = append(args, statusArgs...)
	args = append(args, filterArgs...)

	q := d.Client.Model(&model.Alias{})
	if unscoped {
		q = q.Unscoped()
	}
	q = q.Where("user_id = ? "+statusFilter+filter, args...)

	var count int64
	err := q.Count(&count).Error
	return int(count), err
}

func (d *Database) GetAliasByName(name string) (model.Alias, error) {
	var alias model.Alias
	err := d.Client.Where("name = ?", name).First(&alias).Error
	return alias, err
}

func (d *Database) PostAlias(ctx context.Context, alias model.Alias, maxDaily int, maxInboundHourly int) (model.Alias, error) {
	err := d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedUser model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", alias.UserID).First(&lockedUser).Error; err != nil {
			return err
		}

		// Inbound alias hourly limit check
		if alias.Origin == model.Inbound {
			var hourly int64
			if err := tx.Unscoped().Model(&model.Alias{}).
				Where("user_id = ? AND origin = ? AND created_at > NOW() - INTERVAL 1 HOUR", alias.UserID, model.Inbound).
				Count(&hourly).Error; err != nil {
				return err
			}
			if int(hourly) >= maxInboundHourly {
				return model.ErrInboundHourlyLimit
			}
		}

		// Daily alias limit for non-imported aliases check
		if alias.Origin != model.Import {
			var daily int64
			if err := tx.Unscoped().Model(&model.Alias{}).
				Where("user_id = ? AND created_at > NOW() - INTERVAL 1 DAY", alias.UserID).
				Count(&daily).Error; err != nil {
				return err
			}
			if int(daily) >= maxDaily {
				return model.ErrDailyAliasLimit
			}
		}

		return tx.Create(&alias).Error
	})

	return alias, err
}

func (d *Database) UpdateAlias(ctx context.Context, alias model.Alias) error {
	return d.Client.Model(&alias).Where("user_id = ?", alias.UserID).Updates(map[string]any{
		"description": alias.Description,
		"enabled":     alias.Enabled,
		"recipients":  alias.Recipients,
		"from_name":   alias.FromName,
	}).Error
}

func (d *Database) DeleteAlias(ctx context.Context, ID string, userID string) error {
	return d.Client.Where("id = ? AND user_id = ?", ID, userID).Delete(&model.Alias{}).Error
}

// UpdateAliasPinned sets pinned in isolation (no recipients/description involved), so pinning
// never depends on the alias's recipients being resolvable/verified. GORM's automatic
// soft-delete scope already excludes deleted_at rows here, same as UpdateAlias/DeleteAlias.
func (d *Database) UpdateAliasPinned(ctx context.Context, ID string, userID string, pinned bool) error {
	return d.Client.Model(&model.Alias{}).Where("id = ? AND user_id = ?", ID, userID).Update("pinned", pinned).Error
}

func (d *Database) DeleteAliasByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Alias{}).Error
}

func (d *Database) DeleteAliasByDomain(ctx context.Context, domain string, userID string) error {
	return d.Client.Where("name LIKE ? AND user_id = ?", "%@"+domain, userID).Delete(&model.Alias{}).Error
}

func (d *Database) RestoreAlias(ctx context.Context, ID string, userID string) error {
	return d.Client.Model(&model.Alias{}).Unscoped().Where("id = ? AND user_id = ?", ID, userID).Update("deleted_at", nil).Error
}

// GetAliasUnscoped fetches an alias regardless of soft-delete state, used by ForgetAlias to
// verify ownership/domain before permanently removing the alias.
func (d *Database) GetAliasUnscoped(ctx context.Context, ID string, userID string) (model.Alias, error) {
	var alias model.Alias
	err := d.Client.Unscoped().Where("id = ? AND user_id = ?", ID, userID).First(&alias).Error
	return alias, err
}

// ForgetAlias permanently removes an alias row, regardless of its soft-delete state - the
// same effect as the automatic 90-day cleanup job, but triggered immediately by the user.
func (d *Database) ForgetAlias(ctx context.Context, ID string, userID string) error {
	return d.Client.Unscoped().Where("id = ? AND user_id = ?", ID, userID).Delete(&model.Alias{}).Error
}

// BulkUpdateAliasEnabled sets enabled for all given IDs unconditionally (no eligibility check -
// GORM's automatic soft-delete scope already excludes deleted_at rows, same as UpdateAlias).
func (d *Database) BulkUpdateAliasEnabled(ctx context.Context, ids []string, userID string, enabled bool) error {
	return d.Client.WithContext(ctx).Model(&model.Alias{}).Where("id IN (?) AND user_id = ?", ids, userID).Update("enabled", enabled).Error
}

// BulkUpdateAliasPinned sets pinned for all given IDs unconditionally, mirroring UpdateAliasPinned.
func (d *Database) BulkUpdateAliasPinned(ctx context.Context, ids []string, userID string, pinned bool) error {
	return d.Client.WithContext(ctx).Model(&model.Alias{}).Where("id IN (?) AND user_id = ?", ids, userID).Update("pinned", pinned).Error
}

// BulkDeleteAlias soft-deletes all given IDs, but only if every one of them is currently
// non-deleted - otherwise the whole batch is rolled back via ErrBulkAliasNotEligible.
func (d *Database) BulkDeleteAlias(ctx context.Context, ids []string, userID string) error {
	return d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.Alias{}).Where("id IN (?) AND user_id = ?", ids, userID).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return model.ErrBulkAliasNotEligible
		}

		return tx.Where("id IN (?) AND user_id = ?", ids, userID).Delete(&model.Alias{}).Error
	})
}

// BulkRestoreAlias restores all given IDs, but only if every one of them is currently
// soft-deleted - otherwise the whole batch is rolled back via ErrBulkAliasNotEligible.
func (d *Database) BulkRestoreAlias(ctx context.Context, ids []string, userID string) error {
	return d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Unscoped().Model(&model.Alias{}).Where("id IN (?) AND user_id = ? AND deleted_at IS NOT NULL", ids, userID).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return model.ErrBulkAliasNotEligible
		}

		return tx.Unscoped().Model(&model.Alias{}).Where("id IN (?) AND user_id = ?", ids, userID).Update("deleted_at", nil).Error
	})
}

// GetAliasesUnscopedByIDs fetches aliases regardless of soft-delete state, used by
// BulkForgetAlias to verify ownership/domain before permanently removing the aliases.
func (d *Database) GetAliasesUnscopedByIDs(ctx context.Context, ids []string, userID string) ([]model.Alias, error) {
	var aliases []model.Alias
	err := d.Client.WithContext(ctx).Unscoped().Where("id IN (?) AND user_id = ?", ids, userID).Find(&aliases).Error
	return aliases, err
}

// BulkForgetAlias permanently removes the given alias rows, regardless of soft-delete state.
// Eligibility (custom-domain-only) is validated by the service layer before this is called.
func (d *Database) BulkForgetAlias(ctx context.Context, ids []string, userID string) error {
	return d.Client.WithContext(ctx).Unscoped().Where("id IN (?) AND user_id = ?", ids, userID).Delete(&model.Alias{}).Error
}

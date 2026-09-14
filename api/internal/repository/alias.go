package repository

import (
	"context"
	"strconv"

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
func sanitizeAliasSort(sortBy string, sortOrder string) (string, string) {
	if !model.AliasSortColumns[sortBy] {
		sortBy = "created_at"
	}
	if !model.AliasSortOrders[sortOrder] {
		sortOrder = "DESC"
	}
	return "a." + sortBy, sortOrder
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

func (d *Database) GetAliases(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	sortBy, sortOrder = sanitizeAliasSort(sortBy, sortOrder)

	statusFilter, statusArgs, _ := aliasStatusFilter("a.", status)
	filter, filterArgs := aliasSearchFilter("a.", wildcard, search)

	aliases := []model.Alias{}
	query := `
		SELECT a.*,
			COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS forwards,
			COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS blocks,
			COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS replies,
			COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS sends
		FROM aliases a
		LEFT JOIN messages m
		ON a.id = m.alias_id
		WHERE a.user_id = ? ` + statusFilter + filter + `
		GROUP BY a.id
		ORDER BY ` + sortBy + " " + sortOrder

	if limit > 0 {
		query += "\nLIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		query += "\nOFFSET " + strconv.Itoa(offset)
	}

	args := []any{model.Forward, model.Block, model.Reply, model.Send, userID}
	args = append(args, statusArgs...)
	args = append(args, filterArgs...)

	rows, err := d.Client.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alias model.Alias
		var forwards, blocks, replies, sends int
		if err := rows.Scan(&alias.ID, &alias.CreatedAt, &alias.UpdatedAt, &alias.DeletedAt, &alias.Name, &alias.UserID, &alias.Enabled, &alias.Description, &alias.Recipients, &alias.FromName, &alias.Wildcard, &alias.Origin, &forwards, &blocks, &replies, &sends); err != nil {
			return nil, err
		}
		alias.Stats = model.AliasStats{
			Forwards: forwards,
			Blocks:   blocks,
			Replies:  replies,
			Sends:    sends,
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
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

func (d *Database) DeleteAliasByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Alias{}).Error
}

func (d *Database) DeleteAliasByDomain(ctx context.Context, domain string, userID string) error {
	return d.Client.Where("name LIKE ? AND user_id = ?", "%@"+domain, userID).Delete(&model.Alias{}).Error
}

func (d *Database) RestoreAlias(ctx context.Context, ID string, userID string) error {
	return d.Client.Model(&model.Alias{}).Unscoped().Where("id = ? AND user_id = ?", ID, userID).Update("deleted_at", nil).Error
}

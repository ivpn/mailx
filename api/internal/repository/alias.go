package repository

import (
	"context"
	"strconv"

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

func (d *Database) GetAliases(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, catchAll string, search string) ([]model.Alias, error) {
	sortBy = "a." + sortBy

	if catchAll == "true" {
		catchAll = "AND a.catch_all = true"
	} else if catchAll == "false" {
		catchAll = "AND a.catch_all = false"
	} else {
		catchAll = ""
	}

	if search != "" {
		search = "AND (a.name LIKE '%" + search + "%' OR a.description LIKE '%" + search + "%')"
	}

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
		WHERE a.user_id = ? AND a.deleted_at IS NULL ` + catchAll + " " + search + `
		GROUP BY a.id
		ORDER BY ` + sortBy + " " + sortOrder

	if limit > 0 {
		query += "\nLIMIT " + strconv.Itoa(limit)
	}

	if offset > 0 {
		query += "\nOFFSET " + strconv.Itoa(offset)
	}

	rows, err := d.Client.Raw(query, model.Forward, model.Block, model.Reply, model.Send, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alias model.Alias
		var forwards, blocks, replies, sends int
		if err := rows.Scan(&alias.ID, &alias.CreatedAt, &alias.UpdatedAt, &alias.DeletedAt, &alias.Name, &alias.UserID, &alias.Enabled, &alias.Description, &alias.Recipients, &alias.FromName, &alias.CatchAll, &forwards, &blocks, &replies, &sends); err != nil {
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

func (d *Database) GetAliasCount(ctx context.Context, userID string, catchAll string, search string) (int, error) {
	if catchAll == "true" {
		catchAll = " AND catch_all = true"
	} else if catchAll == "false" {
		catchAll = " AND catch_all = false"
	} else {
		catchAll = ""
	}

	if search != "" {
		search = " AND (name LIKE '%" + search + "%' OR description LIKE '%" + search + "%')"
	}

	var count int64
	err := d.Client.Model(&model.Alias{}).Where("user_id = ?"+catchAll+search, userID).Count(&count).Error
	return int(count), err
}

func (d *Database) GetAliasDailyCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.Model(&model.Alias{}).Where("user_id = ? AND created_at > NOW() - INTERVAL 1 DAY", userID).Count(&count).Error
	return int(count), err
}

func (d *Database) GetAliasByName(name string) (model.Alias, error) {
	var alias model.Alias
	err := d.Client.Where("name = ?", name).First(&alias).Error
	return alias, err
}

func (d *Database) PostAlias(ctx context.Context, alias model.Alias) (model.Alias, error) {
	return alias, d.Client.Create(&alias).Error
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

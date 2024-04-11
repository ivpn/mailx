package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetAlias(ctx context.Context, ID string) (model.Alias, error) {
	var alias model.Alias
	var aliasStats model.AliasStats
	err := d.Client.Where("id = ?", ID).
		First(&alias).Error
	if err != nil {
		return alias, err
	}

	err = d.Client.Model(&model.Message{}).
		Select("SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as forwards, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as blocks, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as replies, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as sends, "+
			"SUM(size) as bandwidth",
			model.Forward, model.Block, model.Reply, model.Send).
		Where("alias_id = ?", ID).
		Scan(&aliasStats).Error
	if err != nil {
		return alias, err
	}

	alias.Stats = aliasStats

	return alias, nil
}

func (d *Database) GetAliases(ctx context.Context, userID string) ([]model.Alias, error) {
	var aliases []model.Alias
	query := `
        SELECT a.*,
            COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS forwards,
            COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS blocks,
            COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS replies,
            COALESCE(SUM(CASE WHEN m.type = ? THEN 1 ELSE 0 END), 0) AS sends,
            COALESCE(SUM(m.size), 0) AS bandwidth
        FROM aliases a
        LEFT JOIN messages m
		ON a.id = m.alias_id
		AND a.user_id = ?
        GROUP BY a.id
    `
	rows, err := d.Client.Raw(query, model.Forward, model.Block, model.Reply, model.Send, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alias model.Alias
		var forwards, blocks, replies, sends, bandwidth int
		if err := rows.Scan(&alias.ID, &alias.CreatedAt, &alias.UpdatedAt, &alias.Name, &alias.UserID, &alias.Enabled, &alias.Description, &alias.Recipients, &alias.FromName, &forwards, &blocks, &replies, &sends, &bandwidth); err != nil {
			return nil, err
		}
		alias.Stats = model.AliasStats{
			Forwards:  forwards,
			Blocks:    blocks,
			Replies:   replies,
			Sends:     sends,
			Bandwidth: bandwidth,
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func (d *Database) GetAliasByName(name string) (model.Alias, error) {
	var alias model.Alias
	err := d.Client.Where("name = ?", name).First(&alias).Error
	return alias, err
}

func (d *Database) PostAlias(ctx context.Context, alias model.Alias) error {
	return d.Client.Create(&alias).Error
}

func (d *Database) UpdateAlias(ctx context.Context, alias model.Alias) error {
	return d.Client.Model(&alias).Updates(map[string]interface{}{
		"description": alias.Description,
		"enabled":     alias.Enabled,
		"recipients":  alias.Recipients,
		"from_name":   alias.FromName,
	}).Error
}

func (d *Database) DeleteAlias(ctx context.Context, ID string) error {
	return d.Client.Where("id = ?", ID).Delete(&model.Alias{}).Error
}

func (d *Database) DeleteAliasByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Alias{}).Error
}

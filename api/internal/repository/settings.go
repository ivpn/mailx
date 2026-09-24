package repository

import (
	"context"
	"fmt"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSettings(ctx context.Context, userID string) (model.Settings, error) {
	var settings model.Settings
	q := d.Client.Where("user_id = ?", userID).Find(&settings)
	if q.RowsAffected == 0 {
		return model.Settings{}, fmt.Errorf("could not get settings by user ID")
	}

	return settings, q.Error
}

func (d *Database) PostSettings(ctx context.Context, settings model.Settings) error {
	return d.Client.Create(&settings).Error
}

func (d *Database) UpdateSettings(ctx context.Context, settings model.Settings) error {
	return d.Client.Model(&settings).Where("user_id = ?", settings.UserID).Updates(map[string]any{
		"domain":        settings.Domain,
		"recipient":     settings.Recipient,
		"from_name":     settings.FromName,
		"alias_format":  settings.AliasFormat,
		"log_issues":    settings.LogIssues,
		"remove_header": settings.RemoveHeader,
	}).Error
}

func (d *Database) UpdateListUnsubscribeBehaviour(ctx context.Context, userID string, behaviour model.ListUnsubscribeBehaviour) error {
	q := d.Client.Model(&model.Settings{}).Where("user_id = ?", userID).Update("list_unsubscribe_behaviour", behaviour)
	if q.Error != nil {
		return q.Error
	}
	if q.RowsAffected == 0 {
		return fmt.Errorf("could not find settings by user ID")
	}

	return nil
}

func (d *Database) DeleteSettings(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Settings{}).Error
}

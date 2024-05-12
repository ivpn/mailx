package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSettings(ctx context.Context, userID string) (model.Settings, error) {
	var settings model.Settings
	err := d.Client.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

func (d *Database) PostSettings(ctx context.Context, settings model.Settings) error {
	return d.Client.Create(&settings).Error
}

func (d *Database) UpdateSettings(ctx context.Context, settings model.Settings) error {
	return d.Client.Where("user_id = ?", settings.UserID).Updates(settings).Error
}

func (d *Database) DeleteSettings(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Settings{}).Error
}

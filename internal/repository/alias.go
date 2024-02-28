package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) GetAlias(ctx context.Context, ID string) (model.Alias, error) {
	var alias model.Alias
	err := d.Client.Where("id = ?", ID).First(&alias).Error
	return alias, err
}

func (d *Database) GetAliases(ctx context.Context, userID string) ([]model.Alias, error) {
	var aliases []model.Alias
	err := d.Client.Where("user_id = ?", userID).Find(&aliases).Error
	return aliases, err
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
	return d.Client.Save(alias).Error
}

func (d *Database) DeleteAlias(ctx context.Context, ID string) error {
	return d.Client.Where("id = ?", ID).Delete(&model.Alias{}).Error
}

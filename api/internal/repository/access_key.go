package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetAccessKeys(ctx context.Context, userId string) ([]model.AccessKey, error) {
	var accessKeys []model.AccessKey
	err := d.Client.Where("user_id = ?", userId).Order("created_at desc").Find(&accessKeys).Error
	if err != nil {
		return nil, err
	}

	return accessKeys, nil
}

func (d *Database) GetAccessKeyByHash(ctx context.Context, tokenHash string) (model.AccessKey, error) {
	var accessKey model.AccessKey
	err := d.Client.Where("token_hash = ?", tokenHash).First(&accessKey).Error
	if err != nil {
		return model.AccessKey{}, err
	}

	return accessKey, nil
}

func (d *Database) PostAccessKey(ctx context.Context, accessKey model.AccessKey) (model.AccessKey, error) {
	err := d.Client.Create(&accessKey).Error
	return accessKey, err
}

func (d *Database) DeleteAccessKey(ctx context.Context, accessKeyID string, userId string) error {
	return d.Client.Where("id = ? AND user_id = ?", accessKeyID, userId).Delete(&model.AccessKey{}).Error
}

func (d *Database) DeleteAccessKeysByUserID(ctx context.Context, userId string) error {
	return d.Client.Where("user_id = ?", userId).Delete(&model.AccessKey{}).Error
}

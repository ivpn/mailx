package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetAccessKeys(ctx context.Context, userID string) ([]model.AccessKey, error) {
	var accessKeys []model.AccessKey
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&accessKeys).Error
	if err != nil {
		return nil, err
	}

	return accessKeys, nil
}

func (d *Database) PostAccessKey(ctx context.Context, accessKey model.AccessKey) error {
	return d.Client.Create(&accessKey).Error
}

func (d *Database) DeleteAccessKey(ctx context.Context, accessKeyID string, userID string) error {
	return d.Client.Where("id = ? AND user_id = ?", accessKeyID, userID).Delete(&model.AccessKey{}).Error
}

func (d *Database) DeleteAccessKeysByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.AccessKey{}).Error
}

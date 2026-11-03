package repository

import (
	"context"

	"gorm.io/gorm"
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

func (d *Database) GetAccessKey(ctx context.Context, id string) (model.AccessKey, error) {
	var accessKey model.AccessKey
	err := d.Client.Where("id = ?", id).First(&accessKey).Error
	if err != nil {
		return model.AccessKey{}, err
	}

	return accessKey, nil
}

func (d *Database) PostAccessKey(ctx context.Context, accessKey model.AccessKey) (model.AccessKey, error) {
	err := d.Client.Create(&accessKey).Error
	return accessKey, err
}

// DeleteAccessKey deletes the key and ends the sessions made from it. It also ends the user's api
// sessions saved before sessions recorded their access key, since any of them may be this key's.
// The key is deleted first, so its row stays locked while the sessions are deleted, and SaveSession
// can't add a session from it in between.
func (d *Database) DeleteAccessKey(ctx context.Context, accessKeyID string, userId string) error {
	return d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ? AND user_id = ?", accessKeyID, userId).Delete(&model.AccessKey{}).Error
		if err != nil {
			return err
		}

		return tx.Where("user_id = ? AND (access_key_id = ? OR (scope = ? AND (access_key_id IS NULL OR access_key_id = '')))",
			userId, accessKeyID, model.SessionScopeAPI).Delete(&model.Session{}).Error
	})
}

// DeleteAccessKeysByUserID deletes the user's keys and ends every api session of the user, including
// ones saved before sessions recorded their access key.
func (d *Database) DeleteAccessKeysByUserID(ctx context.Context, userId string) error {
	return d.Client.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ? AND scope = ?", userId, model.SessionScopeAPI).Delete(&model.Session{}).Error
		if err != nil {
			return err
		}

		return tx.Where("user_id = ?", userId).Delete(&model.AccessKey{}).Error
	})
}

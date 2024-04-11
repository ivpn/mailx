package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetMessagesByUser(ctx context.Context, userID string) ([]model.Message, error) {
	var messages []model.Message
	err := d.Client.Where("user_id = ?", userID).Find(&messages).Error
	return messages, err
}

func (d *Database) GetMessagesByAlias(ctx context.Context, aliasID string) ([]model.Message, error) {
	var messages []model.Message
	err := d.Client.Where("alias_id = ?", aliasID).Find(&messages).Error
	return messages, err
}

func (d *Database) PostMessage(ctx context.Context, message model.Message) error {
	return d.Client.Create(&message).Error
}

func (d *Database) DeleteMessageByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Message{}).Error
}

package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) GetMessages(ctx context.Context, aliasID string) ([]model.Message, error) {
	var messages []model.Message
	err := d.Client.Where("alias_id = ?", aliasID).Find(&messages).Error
	return messages, err
}

func (d *Database) PostMessage(ctx context.Context, message model.Message) error {
	return d.Client.Create(&message).Error
}

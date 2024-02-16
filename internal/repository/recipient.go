package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) GetRecipient(ctx context.Context, ID string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("id = ?", ID).First(&recipient).Error
	return recipient, err
}

func (d *Database) GetRecipients(ctx context.Context, userID string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := d.Client.Where("user_id = ?", userID).Find(&recipients).Error
	return recipients, err
}

func (d *Database) PostRecipient(ctx context.Context, recipient model.Recipient) error {
	return d.Client.Create(&recipient).Error
}

func (d *Database) UpdateRecipient(ctx context.Context, recipient model.Recipient) error {
	return d.Client.Save(recipient).Error
}

func (d *Database) DeleteRecipient(ctx context.Context, ID string) error {
	return d.Client.Where("id = ?", ID).Delete(&model.Recipient{}).Error
}

func (d *Database) VerifyRecipient(ctx context.Context, ID string, verification string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("id = ? AND verification = ?", ID, verification).First(&recipient).Error
	if err != nil {
		return recipient, err
	}

	recipient.Verified = true
	recipient.Verification = ""

	return recipient, d.Client.Save(&recipient).Error
}

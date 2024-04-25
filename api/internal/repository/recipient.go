package repository

import (
	"context"
	"strings"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetRecipient(ctx context.Context, ID string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("id = ?", ID).First(&recipient).Error
	return recipient, err
}

func (d *Database) GetRecipientByEmail(ctx context.Context, email string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("email = ?", email).First(&recipient).Error
	return recipient, err
}

func (d *Database) GetRecipients(ctx context.Context, userID string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := d.Client.Where("user_id = ?", userID).Order("created_at").Find(&recipients).Error
	return recipients, err
}

func (d *Database) GetVerifiedRecipients(ctx context.Context, recipientEmails string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := d.Client.Where("email IN (?) AND is_active = true", strings.Split(recipientEmails, ",")).Find(&recipients).Error
	return recipients, err
}

func (d *Database) PostRecipient(ctx context.Context, recipient model.Recipient) (model.Recipient, error) {
	err := d.Client.Create(&recipient).Error
	return recipient, err
}

func (d *Database) UpdateRecipient(ctx context.Context, recipient model.Recipient) error {
	return d.Client.Updates(recipient).Error
}

func (d *Database) DeleteRecipient(ctx context.Context, ID string) error {
	return d.Client.Where("id = ?", ID).Delete(&model.Recipient{}).Error
}

func (d *Database) ActivateRecipient(ctx context.Context, ID string) error {
	return d.Client.Model(&model.Recipient{}).Where("id = ?", ID).Update("is_active", true).Error
}

func (d *Database) DeleteRecipientByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Recipient{}).Error
}

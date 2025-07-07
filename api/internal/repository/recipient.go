package repository

import (
	"context"
	"strings"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetRecipient(ctx context.Context, ID string, userID string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("id = ? AND user_id = ?", ID, userID).First(&recipient).Error
	return recipient, err
}

func (d *Database) GetRecipientByEmail(ctx context.Context, email string, userID string) (model.Recipient, error) {
	var recipient model.Recipient
	err := d.Client.Where("email = ? AND user_id = ?", email, userID).First(&recipient).Error
	return recipient, err
}

func (d *Database) CheckDuplicateRecipient(ctx context.Context, email string) (bool, error) {
	var count int64
	err := d.Client.Model(&model.Recipient{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Database) GetRecipients(ctx context.Context, userID string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&recipients).Error
	return recipients, err
}

func (d *Database) GetRecipientsCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.Model(&model.Recipient{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

func (d *Database) GetVerifiedRecipients(ctx context.Context, recipientEmails string, userID string) ([]model.Recipient, error) {
	var recipients []model.Recipient
	err := d.Client.Where("email IN (?) AND is_active = true AND user_id = ?", strings.Split(recipientEmails, ","), userID).Find(&recipients).Error
	return recipients, err
}

func (d *Database) PostRecipient(ctx context.Context, recipient model.Recipient) (model.Recipient, error) {
	err := d.Client.Create(&recipient).Error
	return recipient, err
}

func (d *Database) UpdateRecipient(ctx context.Context, recipient model.Recipient) error {
	return d.Client.Model(&recipient).Where("user_id = ?", recipient.UserID).Updates(map[string]any{
		"pgp_key":     recipient.PGPKey,
		"pgp_enabled": recipient.PGPEnabled,
		"pgp_inline":  recipient.PGPInline,
	}).Error
}

func (d *Database) DeleteRecipient(ctx context.Context, ID string, userID string) error {
	return d.Client.Where("id = ? AND user_id = ?", ID, userID).Delete(&model.Recipient{}).Error
}

func (d *Database) ActivateRecipient(ctx context.Context, ID string, userID string) error {
	return d.Client.Model(&model.Recipient{}).Where("id = ? AND user_id = ?", ID, userID).Update("is_active", true).Error
}

func (d *Database) DeleteRecipientByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Recipient{}).Error
}

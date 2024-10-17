package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetCredentials(ctx context.Context, userID string) ([]model.Credential, error) {
	var credentials []model.Credential
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&credentials).Error
	if err != nil {
		return nil, err
	}

	for i := range credentials {
		err = credentials[i].Unmarshal()
		if err != nil {
			return nil, err
		}
	}

	return credentials, nil
}

func (d *Database) GetCredentialsCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.Model(&model.Credential{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

func (d *Database) SaveCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	data, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	return d.Client.Create(&model.Credential{
		UserID: userID,
		Cred:   credential,
		Data:   data,
	}).Error
}

func (d *Database) UpdateCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	data, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	return d.Client.Model(&model.Credential{}).Where("user_id = ?", userID).Update("data", data).Error
}

func (d *Database) DeleteCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	// Get all credentials by user id
	var credentials []model.Credential
	err := d.Client.Where("user_id = ?", userID).Find(&credentials).Error
	if err != nil {
		return err
	}

	// Find the credential to delete
	for _, c := range credentials {
		if bytes.Equal(c.Cred.ID, credential.ID) {
			return d.Client.Where("id = ? AND user_id = ?", c.ID, userID).Delete(&model.Alias{}).Error
		}
	}

	return nil
}

func (d *Database) DeleteCredentialByID(ctx context.Context, ID string, userID string) error {
	return d.Client.Where("id = ? AND user_id = ?", ID, userID).Delete(&model.Credential{}).Error
}

func (d *Database) DeleteCredentialByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Credential{}).Error
}

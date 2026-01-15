package repository

import (
	"context"
	"errors"
	"log"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetUser(ctx context.Context, ID string) (model.User, error) {
	var user model.User
	err := d.Client.Where("id = ?", ID).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	var credentials []model.Credential
	err = d.Client.Where("user_id = ?", user.ID).Find(&credentials).Error
	if err != nil {
		return model.User{}, err
	}

	var creds = []webauthn.Credential{}
	for _, c := range credentials {
		err := c.Unmarshal()
		if err != nil {
			log.Println("error unmarshalling credential:", err)
		}
		creds = append(creds, c.Cred)
	}

	user.Creds = creds

	user.TotpEnabled = user.TotpSecret != ""
	return user, nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := d.Client.Where("email = ?", email).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	var credentials []model.Credential
	err = d.Client.Where("user_id = ?", user.ID).Find(&credentials).Error
	if err != nil {
		return model.User{}, err
	}

	var creds = []webauthn.Credential{}
	for _, c := range credentials {
		err := c.Unmarshal()
		if err != nil {
			log.Println("error unmarshalling credential:", err)
		}
		creds = append(creds, c.Cred)
	}

	user.Creds = creds

	user.TotpEnabled = user.TotpSecret != ""
	return user, nil
}

func (d *Database) GetUserByEmailUnfinishedSignup(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := d.Client.Where("email = ? AND password_hash = ''", email).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	var creds []model.Credential
	err = d.Client.Where("user_id = ?", user.ID).Find(&creds).Error
	if err != nil {
		return model.User{}, err
	}

	if len(creds) > 0 {
		return model.User{}, errors.New("user has already registered webauthn credentials")
	}

	return user, nil
}

func (d *Database) PostUser(ctx context.Context, user model.User) (model.User, error) {
	err := d.Client.Create(&user).Error
	return user, err
}

func (d *Database) SaveUser(ctx context.Context, user model.User) error {
	return d.Client.Save(&user).Error
}

func (d *Database) ActivateUser(ctx context.Context, ID string) error {
	return d.Client.Model(&model.User{}).Where("id = ?", ID).Update("is_active", true).Error
}

func (d *Database) DeleteUser(ctx context.Context, ID string) error {
	return d.Client.Where("id = ?", ID).Delete(&model.User{}).Error
}

func (d *Database) GetUserStats(ctx context.Context, ID string) (model.UserStats, error) {
	var userStats model.UserStats

	err := d.Client.Model(&model.Alias{}).
		Where("user_id = ?", ID).
		Count(&userStats.Aliases).Error
	if err != nil {
		return model.UserStats{}, err
	}

	err = d.Client.Model(&model.Message{}).
		Select("SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as forwards, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as blocks, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as replies, "+
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as sends",
			model.Forward, model.Block, model.Reply, model.Send).
		Where("user_id = ?", ID).
		Where("created_at > NOW() - INTERVAL 90 DAY").
		Scan(&userStats).Error
	if err != nil {
		return model.UserStats{}, err
	}

	var messages = []model.Message{}
	err = d.Client.Where("user_id = ? AND created_at > NOW() - INTERVAL ? DAY", ID, 7).Find(&messages).Error
	if err != nil {
		return model.UserStats{}, err
	}
	userStats.Messages = make([]any, len(messages))
	for i, msg := range messages {
		userStats.Messages[i] = msg
	}

	return userStats, nil
}

func (d *Database) TotpEnable(ctx context.Context, ID string, secret string, backupCodes string) error {
	return d.Client.Model(&model.User{}).Where("id = ?", ID).Updates(map[string]any{
		"totp_secret": secret,
		"totp_backup": backupCodes,
	}).Error
}

func (d *Database) TotpDisable(ctx context.Context, ID string) error {
	return d.Client.Model(&model.User{}).Where("id = ?", ID).Updates(map[string]any{
		"totp_secret":      nil,
		"totp_backup":      nil,
		"totp_backup_used": nil,
	}).Error
}

func (d *Database) TotpGetBackup(ctx context.Context, ID string) (string, string, error) {
	var user model.User
	err := d.Client.Select("totp_backup,totp_backup_used").Where("id = ?", ID).First(&user).Error
	return user.TotpBackup, user.TotpBackupUsed, err
}

func (d *Database) TotpSetUsedBackup(ctx context.Context, ID string, backupCodes string) error {
	return d.Client.Model(&model.User{}).Where("id = ?", ID).Update("totp_backup_used", backupCodes).Error
}

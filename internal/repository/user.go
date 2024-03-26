package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) GetUser(ctx context.Context, ID string) (model.User, error) {
	var user model.User
	err := d.Client.Where("id = ?", ID).First(&user).Error
	return user, err
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := d.Client.Where("email = ?", email).First(&user).Error
	return user, err
}

func (d *Database) PostUser(ctx context.Context, user model.User) (model.User, error) {
	err := d.Client.Create(&user).Error
	return user, err
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
			"SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as sends, "+
			"SUM(size) as bandwidth",
			model.Forward, model.Block, model.Reply, model.Send).
		Where("user_id = ?", ID).
		Where("created_at > NOW() - INTERVAL 90 DAY").
		Scan(&userStats).Error
	if err != nil {
		return model.UserStats{}, err
	}

	return userStats, nil
}

package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

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

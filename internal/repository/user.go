package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) PostUser(ctx context.Context, user model.User) error {
	return d.Client.Create(&user).Error
}

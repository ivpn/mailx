package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetBouncesByUser(ctx context.Context, userID string) ([]model.Bounce, error) {
	var bounces []model.Bounce
	err := d.Client.Where("user_id = ?", userID).Find(&bounces).Error
	return bounces, err
}

func (d *Database) PostBounce(ctx context.Context, bounce model.Bounce) error {
	return d.Client.Create(&bounce).Error
}

func (d *Database) DeleteBounceByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Bounce{}).Error
}

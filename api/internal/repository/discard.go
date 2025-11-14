package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetDiscards(ctx context.Context, userID string) ([]model.Discard, error) {
	var discards []model.Discard
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&discards).Error
	return discards, err
}

func (d *Database) PostDiscard(ctx context.Context, discard model.Discard) error {
	return d.Client.Create(&discard).Error
}

func (d *Database) DeleteDiscards(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Discard{}).Error
}

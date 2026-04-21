package repository

import (
	"context"
	"fmt"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSubscription(ctx context.Context, userID string) (model.Subscription, error) {
	var subscription model.Subscription
	q := d.Client.Where("user_id = ?", userID).Find(&subscription)
	if q.RowsAffected == 0 {
		return model.Subscription{}, fmt.Errorf("could not get subscription by user ID")
	}

	return subscription, q.Error
}

func (d *Database) PostSubscription(ctx context.Context, sub model.Subscription) error {
	return d.Client.Create(&sub).Error
}

func (d *Database) UpdateSubscription(ctx context.Context, sub model.Subscription) error {
	return d.Client.Select("*").Updates(&sub).Error
}

func (d *Database) DeleteSubscription(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Subscription{}).Error
}

package repository

import (
	"context"

	"ivpn.net/email-service/internal/model"
)

func (d *Database) GetSubscription(ctx context.Context, userID string) (model.Subscription, error) {
	var subscription model.Subscription
	err := d.Client.First("user_id = ?", userID).Find(&subscription).Error
	return subscription, err
}

func (d *Database) PostSubscription(ctx context.Context, subscription model.Subscription) error {
	return d.Client.Create(&subscription).Error
}

func (d *Database) UpdateSubscription(ctx context.Context, subscription model.Subscription) error {
	return d.Client.Save(subscription).Error
}

func (d *Database) DeleteSubscription(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Subscription{}).Error
}

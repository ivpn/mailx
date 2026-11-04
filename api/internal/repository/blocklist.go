package repository

import (
	"context"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetBlocklist(ctx context.Context, userID string, limit int, offset int) ([]model.BlocklistItem, error) {
	items := []model.BlocklistItem{}
	q := d.Client.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&items).Error
	return items, err
}

func (d *Database) GetBlocklistCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.WithContext(ctx).Model(&model.BlocklistItem{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

// FindBlocklistMatches returns the user's items matching either the sender address or one of the
// sender's domains, served by the (user_id, type, value) unique index.
func (d *Database) FindBlocklistMatches(ctx context.Context, userID string, email string, domains []string) ([]model.BlocklistItem, error) {
	items := []model.BlocklistItem{}
	q := d.Client.WithContext(ctx).Where("user_id = ?", userID)
	if len(domains) > 0 {
		q = q.Where("((type = ? AND value = ?) OR (type = ? AND value IN (?)))", model.BlocklistEmail, email, model.BlocklistDomain, domains)
	} else {
		q = q.Where("type = ? AND value = ?", model.BlocklistEmail, email)
	}
	err := q.Find(&items).Error
	return items, err
}

func (d *Database) PostBlocklistItem(ctx context.Context, item model.BlocklistItem) (model.BlocklistItem, error) {
	err := d.Client.WithContext(ctx).Create(&item).Error
	return item, err
}

// IncrementBlocklistBlocked bumps the counter in a single UPDATE so concurrent deliveries
// blocked by the same item don't overwrite each other's increments.
func (d *Database) IncrementBlocklistBlocked(ctx context.Context, ID string) error {
	return d.Client.WithContext(ctx).Model(&model.BlocklistItem{}).Where("id = ?", ID).
		UpdateColumn("blocked", gorm.Expr("blocked + ?", 1)).Error
}

// DeleteBlocklistItem returns ErrBlocklistItemNotFound when no row matched, so a missing ID and
// another user's ID are reported the same way, without revealing whether the ID exists.
func (d *Database) DeleteBlocklistItem(ctx context.Context, ID string, userID string) error {
	res := d.Client.WithContext(ctx).Where("id = ? AND user_id = ?", ID, userID).Delete(&model.BlocklistItem{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return model.ErrBlocklistItemNotFound
	}

	return nil
}

func (d *Database) DeleteBlocklistByUserID(ctx context.Context, userID string) error {
	return d.Client.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.BlocklistItem{}).Error
}

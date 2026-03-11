package repository

import (
	"context"

	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	var domains []model.Domain
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&domains).Error
	return domains, err
}

func (d *Database) GetVerifiedDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	var domains []model.Domain
	err := d.Client.Where("user_id = ? AND owner_verified_at IS NOT NULL AND mx_verified_at IS NOT NULL AND send_verified_at IS NOT NULL", userID).Order("created_at desc").Find(&domains).Error
	return domains, err
}

func (d *Database) GetDomain(ctx context.Context, domainID string, userID string) (model.Domain, error) {
	var domain model.Domain
	err := d.Client.Where("id = ? AND user_id = ?", domainID, userID).First(&domain).Error
	return domain, err
}

func (d *Database) GetDomainsCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := d.Client.Model(&model.Domain{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (d *Database) PostDomain(ctx context.Context, domain model.Domain) (model.Domain, error) {
	err := d.Client.Create(&domain).Error
	return domain, err
}

func (d *Database) UpdateDomain(ctx context.Context, domain model.Domain) error {
	return d.Client.Model(&domain).Where("user_id = ?", domain.UserID).Updates(map[string]any{
		"name":              domain.Name,
		"description":       domain.Description,
		"recipient":         domain.Recipient,
		"from_name":         domain.FromName,
		"enabled":           domain.Enabled,
		"owner_verified_at": domain.OwnerVerifiedAt,
		"mx_verified_at":    domain.MXVerifiedAt,
		"send_verified_at":  domain.SendVerifiedAt,
	}).Error
}

func (d *Database) DeleteDomain(ctx context.Context, domainID string, userID string) error {
	return d.Client.Where("id = ? AND user_id = ?", domainID, userID).Delete(&model.Domain{}).Error
}

func (d *Database) DeleteDomainsByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Domain{}).Error
}

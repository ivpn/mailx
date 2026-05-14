package repository

import "context"

func (d *Database) GetManagedUserEmails(ctx context.Context) ([]string, error) {
	var results []struct{ Email string }
	err := d.Client.WithContext(ctx).
		Table("subscriptions").
		Select("users.email").
		Joins("JOIN users ON subscriptions.user_id = users.id").
		Where("subscriptions.type = ?", "Managed").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	emails := make([]string, 0, len(results))
	for _, r := range results {
		emails = append(emails, r.Email)
	}
	return emails, nil
}

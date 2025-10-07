package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

func NotifyExpiringSubscriptionsJob(cfg config.Config, db *gorm.DB) {
	// Reset `notified` for active subscriptions
	UpdateActiveSubscriptions(db)

	// Get expiring subscriptions
	subs, err := GetExpiringSubscriptions(db)
	if err != nil {
		log.Println("Error getting expiring subscriptions:", err)
		return
	}

	if len(subs) == 0 {
		return
	}

	// Send notifications
	log.Printf("Notifying %d expiring subscriptions...", len(subs))
	NotifyExpiringSubscriptions(cfg, db, subs)

	// Mark as notified
	MarkSubscriptionsNotified(db, subs)
}

// Set `notified` to false for all subscriptions that are active
func UpdateActiveSubscriptions(db *gorm.DB) {
	err := db.Model(&model.Subscription{}).
		Where("active_until >= NOW()").
		Update("notified", false).Error
	if err != nil {
		log.Println("Error resetting notified flag for active subscriptions:", err)
	}
}

// Find subscriptions with `notified` false and `active_until` expired 1 day ago
func GetExpiringSubscriptions(db *gorm.DB) ([]model.Subscription, error) {
	subs := []model.Subscription{}
	err := db.Where("notified = false AND active_until < NOW() - INTERVAL 1 DAY").Find(&subs).Error
	if err != nil {
		log.Println("Error fetching expiring subscriptions:", err)
		return nil, err
	}

	return subs, nil
}

// Send email notifications for expiring subscriptions
func NotifyExpiringSubscriptions(cfg config.Config, db *gorm.DB, subs []model.Subscription) {
	for _, sub := range subs {
		// Send email notification
		err := sendSubscriptionExpiryEmail(cfg, db, sub)
		if err != nil {
			log.Println("Error sending subscription expiry email:", err)
			continue
		}

	}
}

// Mark expiring subscriptions as notified
func MarkSubscriptionsNotified(db *gorm.DB, subs []model.Subscription) {
	ids := make([]string, 0, len(subs))
	for _, sub := range subs {
		ids = append(ids, sub.ID)
	}

	err := db.Model(&model.Subscription{}).
		Where("id IN ?", ids).
		Update("notified", true).Error
	if err != nil {
		log.Println("Error marking subscriptions as notified:", err)
	}
}

// Send subscription expiry email
func sendSubscriptionExpiryEmail(cfg config.Config, db *gorm.DB, sub model.Subscription) error {
	user := model.User{}
	err := db.Where("id = ?", sub.UserID).First(&user).Error
	if err != nil {
		return err
	}

	utils.Background(func() {
		data := map[string]any{
			"from": cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(cfg.SMTPClient)
		mailer.Sender = cfg.SMTPClient.Sender
		mailer.SenderName = cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "Limited Access Mode", "expiring_sub.tmpl", data)
		if err != nil {
			log.Printf("error sending expiring subscription email: %s", err.Error())
		}
	})

	return nil
}

package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// Delete unverified users older than 7 days
func DeleteUnverifiedUsers(db *gorm.DB) {
	log.Println("Deleting unverified users older than 7 days")

	users := []model.User{}
	err := db.Where("is_active = ? AND created_at < NOW() - INTERVAL ? DAY", false, 7).Find(&users).Error
	if err != nil {
		log.Println("Error deleting unverified users:", err)
		return
	}

	for _, user := range users {
		userID := user.ID
		// Delete aliases of the user
		err = db.Where("user_id = ?", userID).Delete(&model.Alias{}).Error
		if err != nil {
			log.Println("Error deleting aliases of user:", err)
			return
		}

		// Delete recipients of the user
		err = db.Where("user_id = ?", userID).Delete(&model.Recipient{}).Error
		if err != nil {
			log.Println("Error deleting recipients of user:", err)
			return
		}

		// Delete messages of the user
		err = db.Where("user_id = ?", userID).Delete(&model.Message{}).Error
		if err != nil {
			log.Println("Error deleting messages of user:", err)
			return
		}

		// Delete subscriptions of the user
		err = db.Where("user_id = ?", userID).Delete(&model.Subscription{}).Error
		if err != nil {
			log.Println("Error deleting subscriptions of user:", err)
			return
		}

		// Delete settings of the user
		err = db.Where("user_id = ?", userID).Delete(&model.Settings{}).Error
		if err != nil {
			log.Println("Error deleting settings of user:", err)
			return
		}

		// Delete the user
		err = db.Where("id = ?", userID).Delete(&model.User{}).Error
		if err != nil {
			log.Println("Error deleting user:", err)
			return
		}
	}

	log.Println("Unverified users deleted successfully")
}

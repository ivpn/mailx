package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// Delete unverified users older than 7 days
func DeleteUnverifiedUsers(db *gorm.DB) {
	users := []model.User{}
	err := db.Where("is_active = ? AND created_at < NOW() - INTERVAL ? DAY", false, 7).Find(&users).Error
	if err != nil {
		log.Println("Error deleting unverified users:", err)
		return
	}

	deleteUsers(db, users)
}

func deleteUsers(db *gorm.DB, users []model.User) {
	for _, user := range users {
		ID := user.ID
		// Delete aliases of the user
		err := db.Where("user_id = ?", ID).Delete(&model.Alias{}).Error
		if err != nil {
			log.Println("Error deleting aliases of user:", err)
			return
		}

		// Delete recipients of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Recipient{}).Error
		if err != nil {
			log.Println("Error deleting recipients of user:", err)
			return
		}

		// Delete messages of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Message{}).Error
		if err != nil {
			log.Println("Error deleting messages of user:", err)
			return
		}

		// Delete settings of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Settings{}).Error
		if err != nil {
			log.Println("Error deleting settings of user:", err)
			return
		}

		// Delete subscriptions of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Subscription{}).Error
		if err != nil {
			log.Println("Error deleting subscriptions of user:", err)
			return
		}

		// Delete logs of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Log{}).Error
		if err != nil {
			log.Println("Error deleting logs of user:", err)
			return
		}

		// Delete access keys of the user
		err = db.Where("user_id = ?", ID).Delete(&model.AccessKey{}).Error
		if err != nil {
			log.Println("Error deleting access keys of user:", err)
			return
		}

		// Delete domains of the user
		err = db.Where("user_id = ?", ID).Delete(&model.Domain{}).Error
		if err != nil {
			log.Println("Error deleting domains of user:", err)
			return
		}

		// Delete the user
		err = db.Where("id = ?", ID).Delete(&model.User{}).Error
		if err != nil {
			log.Println("Error deleting user:", err)
			return
		}
	}
}

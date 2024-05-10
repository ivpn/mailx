package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// Delete unverified recipients older than 7 days
func DeleteUnverifiedRecipients(db *gorm.DB) {
	log.Println("Deleting unverified recipients older than 7 days")

	err := db.Where("is_active = ? AND created_at < NOW() - INTERVAL ? DAY", false, 7).Delete(&model.Recipient{}).Error
	if err != nil {
		log.Println("Error deleting unverified recipients:", err)
		return
	}

	log.Println("Unverified recipients deleted successfully")
}

package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// Delete messages older than 90 days
func DeleteOldMessages(db *gorm.DB) {
	err := db.Where("created_at < NOW() - INTERVAL ? DAY", 90).Delete(&model.Message{}).Error
	if err != nil {
		log.Println("Error deleting old messages:", err)
		return
	}
}

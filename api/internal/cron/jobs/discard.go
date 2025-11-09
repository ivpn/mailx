package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

const (
	DiscardExpDays = 7
)

// Delete discards older than 7 days
func DeleteOldDiscards(db *gorm.DB) {
	err := db.Where("created_at < NOW() - INTERVAL ? DAY", DiscardExpDays).Delete(&model.Discard{}).Error
	if err != nil {
		log.Println("Error deleting old discards:", err)
	}
}

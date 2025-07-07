package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

// Cleanup deleted aliases older than 90 days
func CleanupDeletedAliases(db *gorm.DB) {
	err := db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < NOW() - INTERVAL ? DAY", 90).Delete(&model.Alias{}).Error
	if err != nil {
		log.Println("Error cleaning up deleted aliases:", err)
		return
	}
}

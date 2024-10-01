package jobs

import (
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

// Delete expired sessions
func DeleteExpiredSessions(db *gorm.DB, cfg config.APIConfig) {
	exp := cfg.TokenExpiration.Hours()
	err := db.Where("created_at < NOW() - INTERVAL ? HOUR", exp).Delete(&model.Session{}).Error
	if err != nil {
		log.Println("Error deleting expired sessions:", err)
		return
	}
}

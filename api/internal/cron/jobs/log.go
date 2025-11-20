package jobs

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

const (
	EmlLogBaseDir = "/var/log/eml"
	LogExpDays    = 7
)

// Delete logs older than 7 days
func DeleteOldLogs(db *gorm.DB) {
	err := cleanupOldLogFiles(LogExpDays * 24 * time.Hour)
	if err != nil {
		log.Println("Error cleaning up old log files:", err)
	}

	err = db.Where("created_at < NOW() - INTERVAL ? DAY", LogExpDays).Delete(&model.Log{}).Error
	if err != nil {
		log.Println("Error deleting old logs:", err)
	}
}

func cleanupOldLogFiles(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)

	err := filepath.Walk(EmlLogBaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			log.Println("Deleting old file:", path, "(modified", info.ModTime(), ")")
			if err := os.Remove(path); err != nil {
				log.Println("Error deleting file:", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

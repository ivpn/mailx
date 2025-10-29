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
	baseDir = "/var/log/bounce"
	expDays = 14
)

// Delete bounces older than 14 days
func DeleteOldBounces(db *gorm.DB) {
	err := cleanupOldBounceFiles(expDays * 24 * time.Hour)
	if err != nil {
		log.Println("Error cleaning up old bounce files:", err)
	}

	err = db.Where("created_at < NOW() - INTERVAL ? DAY", expDays).Delete(&model.Bounce{}).Error
	if err != nil {
		log.Println("Error deleting old bounces:", err)
	}
}

func cleanupOldBounceFiles(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
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

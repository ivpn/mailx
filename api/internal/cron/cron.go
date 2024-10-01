package cron

import (
	"log"

	"github.com/jasonlvhit/gocron"
	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/cron/jobs"
)

func New(db *gorm.DB) {
	cfg, err := config.New()
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	gocron.Every(1).Hour().Do(jobs.DeleteOldMessages, db)
	gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedRecipients, db)
	gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedUsers, db)
	gocron.Every(1).Hour().Do(jobs.CleanupDeletedAliases, db)
	gocron.Every(1).Hour().Do(jobs.DeleteExpiredSessions, db, cfg.API)
	gocron.Start()
	log.Println("Cron jobs started")
}

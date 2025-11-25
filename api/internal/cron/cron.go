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

	err = gocron.Every(1).Hour().Do(jobs.DeleteOldMessages, db)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedRecipients, db)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedUsers, db)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.CleanupDeletedAliases, db)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.DeleteExpiredSessions, db, cfg.API)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.DeleteExpiredUsers, db, cfg.Service)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	err = gocron.Every(1).Hour().Do(jobs.DeleteOldLogs, db)
	if err != nil {
		log.Println("Error scheduling job:", err)
		return
	}

	gocron.Start()

	log.Println("Cron jobs started")
}

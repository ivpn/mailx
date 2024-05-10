package cron

import (
	"log"

	"github.com/jasonlvhit/gocron"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/cron/jobs"
)

func New(db *gorm.DB) {
	gocron.Every(1).Hour().Do(jobs.DeleteOldMessages, db)
	gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedRecipients, db)
	gocron.Every(1).Hour().Do(jobs.DeleteUnverifiedUsers, db)
	gocron.Start()
	log.Println("Cron jobs started")
}

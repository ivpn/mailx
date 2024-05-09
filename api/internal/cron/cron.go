package cron

import (
	"github.com/jasonlvhit/gocron"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/cron/jobs"
)

func New(db *gorm.DB) {
	gocron.Every(1).Hour().Do(jobs.DeleteOldMessages, db)
	gocron.Every(1).Day().Do(jobs.DeleteUnverifiedRecipients, db)
	gocron.Start()
}

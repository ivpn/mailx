package cron

import (
	"github.com/jasonlvhit/gocron"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/cron/jobs"
)

func New(db *gorm.DB) {
	gocron.Every(1).Day().Do(jobs.DeleteOldMessages, db)
	gocron.Start()
}

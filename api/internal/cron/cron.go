package cron

import (
	"log"
	"time"

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

	every(time.Hour, "DeleteOldMessages", func() { jobs.DeleteOldMessages(db) }, nil)
	every(time.Hour, "DeleteUnverifiedRecipients", func() { jobs.DeleteUnverifiedRecipients(db) }, nil)
	// every(time.Hour, "DeleteUnverifiedUsers", func() { jobs.DeleteUnverifiedUsers(db) }, nil)
	every(time.Hour, "DeleteTerminatedUsers", func() { jobs.DeleteTerminatedUsers(db) }, nil)
	every(time.Hour, "CleanupDeletedAliases", func() { jobs.CleanupDeletedAliases(db) }, nil)
	every(time.Hour, "DeleteExpiredSessions", func() { jobs.DeleteExpiredSessions(db, cfg.API) }, nil)
	every(time.Hour, "BackfillSessionScopes", func() { jobs.BackfillSessionScopes(db, cfg.API) }, nil)
	every(time.Hour, "NotifyExpiringSubscriptions", func() { jobs.NotifyExpiringSubscriptionsJob(cfg, db) }, nil)
	every(time.Hour, "DeleteOldLogs", func() { jobs.DeleteOldLogs(db) }, nil)
	every(12*time.Hour, "VerifyDomains", func() { jobs.VerifyDomainsJob(cfg, db) }, nil)

	log.Println("Cron jobs started")
}

// every runs fn once per interval, first after one interval has elapsed, until
// stop is closed (nil runs forever). Runs of the same job never overlap, and a
// panic in fn is logged instead of crashing the process.
func every(interval time.Duration, name string, fn func(), stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				runJob(name, fn)
			}
		}
	}()
}

func runJob(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in cron job %s: %v", name, r)
		}
	}()

	fn()
}

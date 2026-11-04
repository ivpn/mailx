package jobs

import (
	"context"
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/repository"
)

// BackfillSessionScopes sets the scope of sessions saved without one: before scopes existed, or by an
// instance still on the old version during a rolling deploy. Until it runs, such sessions are refused
// by both auth middlewares.
func BackfillSessionScopes(db *gorm.DB, cfg config.APIConfig) {
	if cfg.ApiTokenExpiration <= cfg.TokenExpiration {
		log.Println("BackfillSessionScopes: API_TOKEN_EXPIRATION isn't longer than TOKEN_EXPIRATION, so access-key sessions can't be told apart and are set to web")
	}

	repo := &repository.Database{Client: db}
	n, err := repo.BackfillSessionScopes(context.Background(), cfg.TokenExpiration)
	if err != nil {
		log.Println("Error backfilling session scopes:", err)
		return
	}
	if n > 0 {
		log.Printf("BackfillSessionScopes: set the scope of %d sessions", n)
	}
}

// Delete expired sessions
func DeleteExpiredSessions(db *gorm.DB, cfg config.APIConfig) {
	err := db.Where("expires_at < NOW()").Delete(&model.Session{}).Error
	if err != nil {
		log.Println("Error deleting expired sessions:", err)
		return
	}
}

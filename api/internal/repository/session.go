package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ivpn.net/email/api/internal/model"
)

var ErrAccessKeyNotFound = errors.New("access key not found")

func (d *Database) GetSession(ctx context.Context, token string) (model.Session, bool, error) {
	var session model.Session
	// Expired rows are only purged hourly by cron, so expiry is enforced here as well
	q := d.Client.Where("token = ? AND expires_at > ?", token, time.Now()).Find(&session)
	if q.RowsAffected == 0 {
		return model.Session{}, false, fmt.Errorf("could not get session by token")
	}

	err := session.UnmarshalSessionData()
	if err != nil {
		log.Println("error unmarshalling session data:", err)
	}

	return session, true, q.Error
}

func (d *Database) GetSessionCount(ctx context.Context, userID string) (int, error) {
	var count int64
	err := d.Client.Model(&model.Session{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

// SaveSession refuses a session without a valid scope, so every caller decides what it may reach.
// A session made from an access key is saved only while the key exists. The key row is locked until
// the insert, so a DeleteAccessKey running at the same time either ends this session or makes it
// fail, and never leaves it outliving the key.
func (d *Database) SaveSession(ctx context.Context, session model.Session) error {
	if !session.Scope.Valid() {
		return model.ErrInvalidSessionScope
	}

	data, err := json.Marshal(session.SessionData)
	if err != nil {
		return err
	}

	return d.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if session.AccessKeyID != "" {
			var key model.AccessKey
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
				Where("id = ? AND user_id = ?", session.AccessKeyID, session.UserID).First(&key).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAccessKeyNotFound
			}
			if err != nil {
				return err
			}
		}

		return tx.Create(&model.Session{
			UserID:      session.UserID,
			Token:       session.Token,
			Data:        data,
			ExpiresAt:   session.ExpiresAt,
			Remember:    session.Remember,
			Scope:       session.Scope,
			AccessKeyID: session.AccessKeyID,
		}).Error
	})
}

// BackfillSessionScopes sets the scope of sessions saved before sessions had one, and returns how
// many it set. A session without "remember" that was made to last longer than tokenExpiration (plus
// a minute of slack for the time between computing the expiry and saving the row) came from an
// access key, so it's api. Every other session is web.
//
// An access-key session made when its key had tokenExpiration or less left can't be told apart from
// a web session, and is set to web. It expires within tokenExpiration of being made, so the window
// closes at most tokenExpiration after the deploy.
func (d *Database) BackfillSessionScopes(ctx context.Context, tokenExpiration time.Duration) (int, error) {
	const batchSize = 1000
	maxWebLifetime := tokenExpiration + time.Minute
	total := 0

	for {
		var sessions []model.Session
		err := d.Client.WithContext(ctx).
			Select("id", "created_at", "expires_at", "remember").
			Where("scope IS NULL OR scope = ?", "").
			Limit(batchSize).
			Find(&sessions).Error
		if err != nil {
			return total, err
		}
		if len(sessions) == 0 {
			return total, nil
		}

		var webIDs, apiIDs []string
		for _, s := range sessions {
			if !s.Remember && s.ExpiresAt.Sub(s.CreatedAt) > maxWebLifetime {
				apiIDs = append(apiIDs, s.ID)
			} else {
				webIDs = append(webIDs, s.ID)
			}
		}

		for scope, ids := range map[model.SessionScope][]string{model.SessionScopeWeb: webIDs, model.SessionScopeAPI: apiIDs} {
			if len(ids) == 0 {
				continue
			}
			err := d.Client.WithContext(ctx).Model(&model.Session{}).Where("id IN ?", ids).UpdateColumn("scope", scope).Error
			if err != nil {
				return total, err
			}
		}

		total += len(sessions)
	}
}

func (d *Database) DeleteSession(ctx context.Context, token string) error {
	return d.Client.Where("token = ?", token).Delete(&model.Session{}).Error
}

func (d *Database) DeleteSessionByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

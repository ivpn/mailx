package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSession(ctx context.Context, token string) (model.Session, bool, error) {
	var session model.Session
	q := d.Client.Where("token = ?", token).Find(&session)
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

func (d *Database) SaveSession(ctx context.Context, sessionData webauthn.SessionData, token string, userID string, exp time.Time) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return d.Client.Create(&model.Session{
		UserID:    userID,
		Token:     token,
		Data:      data,
		ExpiresAt: exp,
	}).Error
}

func (d *Database) DeleteSession(ctx context.Context, token string) error {
	return d.Client.Where("token = ?", token).Delete(&model.Session{}).Error
}

func (d *Database) DeleteSessionByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

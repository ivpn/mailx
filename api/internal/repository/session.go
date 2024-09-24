package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSession(ctx context.Context, token string) (webauthn.SessionData, bool, error) {
	var session model.Session
	q := d.Client.Where("token = ?", token).Find(&session)
	if q.RowsAffected == 0 {
		return webauthn.SessionData{}, false, fmt.Errorf("could not get session by token")
	}

	var sessionData webauthn.SessionData
	err := json.Unmarshal(session.SessionData, &sessionData)
	if err != nil {
		return webauthn.SessionData{}, false, err
	}

	return sessionData, true, q.Error
}

func (d *Database) SaveSession(ctx context.Context, session webauthn.SessionData, token string, userID string) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return d.Client.Create(&model.Session{
		UserID:      userID,
		Token:       token,
		SessionData: sessionData,
	}).Error
}

func (d *Database) DeleteSession(ctx context.Context, token string) error {
	return d.Client.Where("token = ?", token).Delete(&model.Session{}).Error
}

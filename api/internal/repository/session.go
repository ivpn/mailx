package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

func (d *Database) GetSession(ctx context.Context, token string) (model.Session, bool, error) {
	var session model.Session
	q := d.Client.Where("token = ?", token).Find(&session)
	if q.RowsAffected == 0 {
		return model.Session{}, false, fmt.Errorf("could not get session by token")
	}

	session.UnmarshalSessionData()

	return session, true, q.Error
}

func (d *Database) SaveSession(ctx context.Context, sessionData webauthn.SessionData, token string, userID string) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return d.Client.Create(&model.Session{
		UserID: userID,
		Token:  token,
		Data:   data,
	}).Error
}

func (d *Database) DeleteSession(ctx context.Context, token string) error {
	return d.Client.Where("token = ?", token).Delete(&model.Session{}).Error
}

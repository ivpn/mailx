package model

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
)

type Session struct {
	BaseModel
	UserID      string               `json:"-"`
	Token       string               `gorm:"unique" json:"token"`
	Data        []byte               `gorm:"type:blob" json:"-"`
	SessionData webauthn.SessionData `gorm:"-" json:"-"`
}

func GenSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *Session) UnmarshalSessionData() error {
	var data webauthn.SessionData
	if err := json.Unmarshal(s.Data, &data); err != nil {
		return err
	}

	s.SessionData = data
	return nil
}

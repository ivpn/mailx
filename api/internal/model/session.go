package model

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"

	"github.com/go-webauthn/webauthn/webauthn"
)

type Session struct {
	BaseModel
	UserID      string               `json:"-"`
	Token       string               `gorm:"unique" json:"token"`
	Data        []byte               `gorm:"type:blob" json:"-"`
	SessionData webauthn.SessionData `gorm:"-" json:"-"`
}

func GenSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *Session) UnmarshalSessionData() error {
	var data webauthn.SessionData
	if err := json.Unmarshal(s.Data, &data); err != nil {
		return err
	}

	s.SessionData = data
	return nil
}

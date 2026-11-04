package model

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

var ErrInvalidSessionScope = errors.New("session scope is invalid")

// SessionScope is what a session may reach. It's set when the session is saved, and checked by the
// auth middleware of each route group.
type SessionScope string

const (
	// SessionScopeWeb is a sign-in on the web app or the iOS app: every /v1 route
	SessionScopeWeb SessionScope = "web"
	// SessionScopeAPI is a session made from an access key: the alias API (/v1/api) only
	SessionScopeAPI SessionScope = "api"
)

func (s SessionScope) Valid() bool {
	return s == SessionScopeWeb || s == SessionScopeAPI
}

type Session struct {
	BaseModel
	UserID      string               `json:"-"`
	Token       string               `gorm:"unique" json:"token"`
	Data        []byte               `gorm:"type:blob" json:"-"`
	SessionData webauthn.SessionData `gorm:"-" json:"-"`
	ExpiresAt   time.Time            `json:"expires_at"`
	Remember    bool                 `json:"-"`
	Scope       SessionScope         `gorm:"size:8" json:"-"`
	AccessKeyID string               `gorm:"size:36;index" json:"-"` // the access key an api session was made from
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

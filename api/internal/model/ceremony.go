package model

import "github.com/go-webauthn/webauthn/webauthn"

// Ceremony is the server-side state of an in-progress WebAuthn ceremony (begin -> finish).
// It is kept apart from Session so a ceremony handle can never be presented as a login session.
type Ceremony struct {
	UserID      string               `json:"user_id"`
	SessionData webauthn.SessionData `json:"session_data"`
	Remember    bool                 `json:"remember"`
}

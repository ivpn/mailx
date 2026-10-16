package model

import "github.com/go-webauthn/webauthn/webauthn"

// CeremonyType records which begin handler created a ceremony, so a finish handler only accepts
// state from its own begin. The library's SessionData can't tell a login challenge from a
// registration one, and the client data type is written by the client.
type CeremonyType string

const (
	CeremonyRegistration CeremonyType = "registration"
	CeremonyLogin        CeremonyType = "login"
	CeremonyPasskeyLogin CeremonyType = "passkey_login"
	CeremonyAddPasskey   CeremonyType = "add_passkey"
	CeremonyStepUp       CeremonyType = "step_up"
)

// Ceremony is the server-side state of an in-progress WebAuthn ceremony (begin -> finish).
// It is kept apart from Session so a ceremony handle can never be presented as a login session.
type Ceremony struct {
	Type        CeremonyType         `json:"type"`
	UserID      string               `json:"user_id"`
	SessionData webauthn.SessionData `json:"session_data"`
	Remember    bool                 `json:"remember"`
}

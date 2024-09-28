package model

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrHashFailed     = errors.New("password hash failed")
)

type User struct {
	BaseModel
	Email          string                `gorm:"unique" json:"email"`
	PasswordHash   string                `json:"-"`
	PasswordPlain  *string               `gorm:"-" json:"-"`
	IsActive       bool                  `json:"is_active"`
	TotpSecret     string                `json:"-"`
	TotpBackup     string                `json:"-"`
	TotpBackupUsed string                `json:"-"`
	TotpEnabled    bool                  `gorm:"-" json:"totp_enabled"`
	creds          []webauthn.Credential `gorm:"-" json:"-"`
	Credentials    []byte                `gorm:"type:blob" json:"-"`
}

type UserStats struct {
	Forwards  int           `json:"forwards"`
	Blocks    int           `json:"blocks"`
	Replies   int           `json:"replies"`
	Sends     int           `json:"sends"`
	Bandwidth int           `json:"bandwidth"`
	Aliases   int64         `json:"aliases"`
	Messages  []interface{} `json:"messages" gorm:"type:text"`
}

func (u *User) SetPassword(passwordPlain string) error {
	hash, err := utils.Hash(passwordPlain)
	if err != nil {
		return ErrHashFailed
	}

	u.PasswordHash = hash
	u.PasswordPlain = nil

	return nil
}

func (u *User) Matches(passwordPlain string) bool {
	return utils.HashMatches(passwordPlain, u.PasswordHash)
}

func (u *User) IsTotpEnabled() bool {
	return u.TotpSecret != ""
}

func (u *User) VerifyTotp(otp string) (bool, error) {
	code, err := strconv.Atoi(otp)
	if err != nil {
		return false, err
	}

	return utils.Check(u.TotpSecret, code)
}

// WebAuthnCredentials implements webauthn.User
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.creds
}

// WebAuthnDisplayName implements webauthn.User
func (u User) WebAuthnDisplayName() string {
	return u.Email
}

// WebAuthnID implements webauthn.User
func (u User) WebAuthnID() []byte {
	return []byte(u.ID)
}

// WebAuthnIcon implements webauthn.User
func (u User) WebAuthnIcon() string {
	return ""
}

// WebAuthnName implements webauthn.User
func (u User) WebAuthnName() string {
	return u.Email
}

func (u *User) AddCredential(credential *webauthn.Credential) {
	u.creds = append(u.creds, *credential)

	credsJSON, err := json.Marshal(u.creds)
	if err == nil {
		u.Credentials = credsJSON
	}
}

func (u *User) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range u.creds {
		if string(c.ID) == string(credential.ID) {
			u.creds[i] = *credential
		}
	}

	credsJSON, err := json.Marshal(u.creds)
	if err == nil {
		u.Credentials = credsJSON
	}
}

func (u *User) UnmarshalCreds() error {
	if u.Credentials == nil {
		return nil
	}

	var creds []webauthn.Credential
	if err := json.Unmarshal(u.Credentials, &creds); err != nil {
		return err
	}

	u.creds = creds
	return nil
}

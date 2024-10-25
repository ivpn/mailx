package model

import (
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
)

type Credential struct {
	BaseModel
	UserID string              `json:"user_id"`
	CredID []byte              `gorm:"-" json:"-"`
	Cred   webauthn.Credential `gorm:"-" json:"-"`
	Data   []byte              `gorm:"type:blob" json:"-"`
}

func (c *Credential) Unmarshal() error {
	if c.Data == nil {
		return nil
	}

	var cred webauthn.Credential
	if err := json.Unmarshal(c.Data, &cred); err != nil {
		return err
	}

	c.Cred = cred
	c.CredID = cred.ID

	return nil
}

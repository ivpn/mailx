package model

import (
	"encoding/base64"
	"math/rand"
)

type Session struct {
	BaseModel
	UserID      string `json:"-"`
	Token       string `gorm:"unique" json:"token"`
	SessionData []byte `gorm:"type:blob" json:"-"`
}

func GenSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

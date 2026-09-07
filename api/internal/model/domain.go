package model

import (
	"errors"
	"time"
)

var (
	ErrDuplicateDomain = errors.New("domain already exists")
)

type Domain struct {
	BaseModel
	UserID          string     `json:"-"`
	Name            string     `gorm:"unique" json:"name"`
	Description     string     `gorm:"default:''" json:"description"`
	Recipient       string     `gorm:"default:''" json:"recipient"`
	FromName        string     `gorm:"default:''" json:"from_name"`
	Enabled         bool       `json:"enabled"`
	OwnerVerifiedAt *time.Time `json:"owner_verified_at"` // nullable
	MXVerifiedAt    *time.Time `json:"mx_verified_at"`    // nullable
	SendVerifiedAt  *time.Time `json:"send_verified_at"`  // nullable
	CatchAll        bool       `gorm:"default:false" json:"catch_all"`
	CreateAlias     bool       `gorm:"default:false" json:"create_alias"`
}

type DNSConfig struct {
	Verify string   `json:"verify"`
	Domain string   `json:"domain"`
	DKIM   []string `json:"dkim_selectors"`
	Hosts  []string `json:"mx_hosts"`
}

type RecordCheck struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Error  string `json:"error,omitempty"`
}

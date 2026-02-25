package model

import "time"

type Domain struct {
	BaseModel
	UserID             string     `json:"-"`
	Name               string     `gorm:"unique" json:"name"`
	Description        string     `gorm:"default:''" json:"description"`
	Recipient          string     `gorm:"default:''" json:"recipient"`
	FromName           string     `gorm:"default:''" json:"from_name"`
	Enabled            bool       `json:"enabled"`
	OwnerVerifiedAt    *time.Time `json:"owner_verified_at"`    // nullable
	InboundVerifiedAt  *time.Time `json:"inbound_verified_at"`  // nullable
	OutboundVerifiedAt *time.Time `json:"outbound_verified_at"` // nullable
}

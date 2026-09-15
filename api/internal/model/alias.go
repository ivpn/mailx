package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrDuplicateAlias       = errors.New("alias already exists")
	ErrDuplicateAliasDomain = errors.New("wildcard aliases limit reached for this domain")
	ErrDailyAliasLimit      = errors.New("daily alias limit reached")
	ErrInboundHourlyLimit   = errors.New("hourly inbound alias limit reached")
)

type AliasOrigin int

const (
	Manual  AliasOrigin = 0
	Inbound AliasOrigin = 1
	Import  AliasOrigin = 2
)

// Scan handles NULL origin values from rows predating the column addition.
func (a *AliasOrigin) Scan(src any) error {
	if src == nil {
		*a = Manual
		return nil
	}
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("AliasOrigin: unsupported scan type %T", src)
	}
	*a = AliasOrigin(v)
	return nil
}

type Alias struct {
	BaseModel
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name             string         `gorm:"unique" json:"name"`
	UserID           string         `json:"-"`
	Enabled          bool           `json:"enabled"`
	Description      string         `gorm:"default:''" json:"description"`
	Recipients       string         `gorm:"default:''" json:"recipients"`
	FromName         string         `gorm:"default:''" json:"from_name"`
	CatchAll         bool           `json:"catch_all"`
	Origin           AliasOrigin    `json:"origin"`
	Stats            AliasStats     `gorm:"-" json:"stats"`
	IsCustomDomain   bool           `gorm:"-" json:"is_custom_domain"`
	IsDomainVerified *bool          `gorm:"-" json:"is_domain_verified"`
	IsDomainEnabled  bool           `gorm:"-" json:"is_domain_enabled"`
}

type AliasStats struct {
	Forwards int `json:"forwards"`
	Blocks   int `json:"blocks"`
	Replies  int `json:"replies"`
	Sends    int `json:"sends"`
}

type AliasList struct {
	Aliases []Alias `json:"aliases"`
	Total   int     `json:"total"`
}

type AliasImportReq struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Recipients  string `json:"recipients" validate:"required"`
	FromName    string `json:"from_name"`
	Format      string `json:"format"`
	Domain      string `json:"domain" validate:"required"`
	LocalPart   string `json:"local_part" validate:"omitempty,min=6,max=24"`
}

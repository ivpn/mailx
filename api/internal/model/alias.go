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
	ErrBulkAliasNotEligible = errors.New("one or more selected aliases are not eligible for this action")
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
	DeletedAt        gorm.DeletedAt `gorm:"index;index:idx_aliases_user_id_deleted_at,priority:2" json:"deleted_at"`
	Name             string         `gorm:"unique" json:"name"`
	UserID           string         `json:"-" gorm:"index:idx_aliases_user_id_deleted_at,priority:1"`
	Enabled          bool           `json:"enabled"`
	Description      string         `gorm:"default:''" json:"description"`
	Recipients       string         `gorm:"default:''" json:"recipients"`
	FromName         string         `gorm:"default:''" json:"from_name"`
	Wildcard         bool           `gorm:"column:catch_all" json:"wildcard"`
	Origin           AliasOrigin    `json:"origin"`
	Pinned           bool           `gorm:"default:false" json:"pinned"`
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

// WildcardDomainInfo describes a user's existing Wildcard Aliases for a single domain, used
// by the frontend to know whether/which delimiters can still be used for that domain.
type WildcardDomainInfo struct {
	Count          int      `json:"count"`
	Limit          int      `json:"limit"`
	DelimitersUsed []string `json:"delimiters_used"`
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

// AliasSortColumns is the single source of truth for GetAliases sort_by values, enforced
// independently by both the API handler and the repository query builder.
var AliasSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
}

// AliasSortOrders is the single source of truth for GetAliases sort_order values, enforced
// independently by both the API handler and the repository query builder.
var AliasSortOrders = map[string]bool{
	"ASC":  true,
	"DESC": true,
}

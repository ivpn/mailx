package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrDuplicateAlias       = errors.New("alias already exists")
	ErrDuplicateAliasDomain = errors.New("wildcard aliases limit reached for this domain")
)

type Alias struct {
	BaseModel
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"unique" json:"name"`
	UserID      string         `json:"-"`
	Enabled     bool           `json:"enabled"`
	Description string         `gorm:"default:''" json:"description"`
	Recipients  string         `gorm:"default:''" json:"recipients"`
	FromName    string         `gorm:"default:''" json:"from_name"`
	CatchAll    bool           `json:"catch_all"`
	Stats       AliasStats     `gorm:"-" json:"stats"`
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
	LocalPart   string `json:"local_part" validate:"omitempty,alphanum,min=6,max=24"`
}

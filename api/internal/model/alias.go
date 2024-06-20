package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrDuplicateAlias = errors.New("alias already exists")
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
	Stats       AliasStats     `gorm:"-" json:"stats"`
}

type AliasStats struct {
	Forwards  int `json:"forwards"`
	Blocks    int `json:"blocks"`
	Replies   int `json:"replies"`
	Sends     int `json:"sends"`
	Bandwidth int `json:"bandwidth"`
}

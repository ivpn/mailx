package model

import "errors"

var (
	ErrDuplicateAlias = errors.New("alias already exists")
)

type Alias struct {
	BaseModel
	UserID      string     `json:"-"`
	Name        string     `gorm:"unique" json:"name"`
	Enabled     bool       `json:"enabled"`
	Description string     `json:"description"`
	Recipients  string     `json:"recipients"`
	Stats       AliasStats `gorm:"-" json:"stats"`
}

type AliasStats struct {
	Forwards  int `json:"forwards"`
	Blocks    int `json:"blocks"`
	Replies   int `json:"replies"`
	Sends     int `json:"sends"`
	Bandwidth int `json:"bandwidth"`
}

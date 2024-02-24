package model

import "errors"

var (
	ErrDuplicateAlias = errors.New("duplicate alias")
)

type Alias struct {
	BaseModel
	RecipientID string `json:"recipient_id"`
	Name        string `gorm:"unique" json:"name"`
	Descripion  string `json:"description"`
}

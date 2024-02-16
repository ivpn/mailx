package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDBaseModel struct {
	PK        uint      `gorm:"primaryKey"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (base *UUIDBaseModel) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

package model

import "time"

type MessageType int

const (
	Forward MessageType = 0
	Block   MessageType = 1
	Reply   MessageType = 2
	Send    MessageType = 3
)

type Message struct {
	ID        uint        `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time   `json:"created_at"`
	UserID    string      `json:"-"`
	AliasID   string      `json:"-"`
	Type      MessageType `json:"type"`
	Size      int         `json:"-"`
}

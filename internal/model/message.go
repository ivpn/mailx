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
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UserID    string
	AliasID   string
	Type      MessageType
	Size      int
}

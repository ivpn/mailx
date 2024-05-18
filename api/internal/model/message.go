package model

import (
	"strings"
	"time"
)

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

func ParseReplyTo(email string) (string, string) {
	alias := email

	// Get respond to email between "+" and "@"
	rcp := email[strings.Index(email, "+")+1 : strings.Index(email, "@")]

	// Check if respond to email is not empty and contains "="
	if rcp != "" && strings.Contains(rcp, "=") {
		// Replace "=" with "@" to get valid respond to email
		rcp = strings.Replace(rcp, "=", "@", 1)

		// Get alias name up to "+" and domain after "@"
		alias = email[:strings.Index(email, "+")] + email[strings.Index(email, "@"):]
	} else {
		rcp = ""
	}

	return alias, rcp
}

func GenerateReplyTo(alias string, to string) string {
	replaced := strings.Replace(to, "@", "=", 1)
	email := strings.Replace(alias, "@", "+"+replaced+"@", 1)
	return email
}

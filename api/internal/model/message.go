package model

import (
	"html"
	"strings"
	"time"
)

type MessageType int

const (
	Forward    MessageType = 0
	Block      MessageType = 1
	Reply      MessageType = 2
	Send       MessageType = 3
	FailBounce MessageType = 4
)

type Message struct {
	ID        uint        `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time   `json:"created_at"`
	UserID    string      `json:"-"`
	AliasID   string      `json:"-"`
	Type      MessageType `json:"type"`
}

func ParseReplyTo(email string) (string, string) {
	alias := email

	// Check if "+" exists in the email
	plusIndex := strings.Index(email, "+")
	if plusIndex == -1 {
		return alias, ""
	}

	// Get respond to email between "+" and "@"
	if plusIndex == -1 {
		return alias, ""
	}

	// Get respond to email between "+" and "@"
	rcp := email[plusIndex+1 : strings.Index(email, "@")]

	// Check if respond to email is not empty and contains "=" and "+"
	if rcp != "" && strings.Contains(rcp, "=") && strings.Contains(rcp, "+") {
		// Get rcp after "+" and replace "=" with "@" to get valid respond to email
		rcp = rcp[strings.Index(rcp, "+")+1:]
		rcp = strings.Replace(rcp, "=", "@", 1)

		// Get the string between the two "+" and the domain after "@"
		alias = "*" + email[strings.Index(email, "+"):strings.LastIndex(email, "+")] + email[strings.Index(email, "@"):]
		return alias, rcp
	}

	// Check if respond to email is not empty and contains "="
	if rcp != "" && strings.Contains(rcp, "=") {
		// Replace "=" with "@" to get valid respond to email
		rcp = strings.Replace(rcp, "=", "@", 1)

		// Get alias name up to "+" and domain after "@"
		alias = email[:strings.Index(email, "+")] + email[strings.Index(email, "@"):]
		return alias, rcp
	}

	// If there is "+" in the email, convert alias to catch-all format
	if rcp != "" && strings.Contains(email, "+") {
		alias = "*" + email[strings.Index(email, "+"):]
	}

	return alias, ""
}

func GenerateReplyTo(alias string, to string) string {
	replaced := strings.Replace(to, "@", "=", 1)
	email := strings.Replace(alias, "@", "+"+replaced+"@", 1)
	return email
}

func PlainTextToHTML(text string) string {
	escaped := html.EscapeString(text)
	html := strings.ReplaceAll(escaped, "\n", "<br>")
	return html
}

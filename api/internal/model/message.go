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
	AliasID   string      `json:"-" gorm:"index:idx_messages_alias_id_type,priority:1"`
	Type      MessageType `json:"type" gorm:"index:idx_messages_alias_id_type,priority:2"`
}

// ReplyTo is one way to read a reply-encoded address: the alias name and the address it replies to.
type ReplyTo struct {
	Alias string
	To    string
}

// ReplyToCandidates returns each way a reply-encoded address ("alias+to=example.com@domain.com")
// splits into an alias name and the address it replies to, longest alias name first. Both may
// contain "+", so "a+b+c=example.com@domain.com" reads as "a+b@domain.com" replying to
// "c@example.com" or as "a@domain.com" replying to "b+c@example.com". It returns nil for an
// address that isn't reply-encoded.
func ReplyToCandidates(email string) []ReplyTo {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return nil
	}

	local, domain := email[:atIndex], email[atIndex:]
	var candidates []ReplyTo
	for i := strings.LastIndex(local, "+"); i > 0; i = strings.LastIndex(local[:i], "+") {
		to := local[i+1:]
		if to == "" || !strings.Contains(to, "=") {
			continue
		}
		candidates = append(candidates, ReplyTo{Alias: local[:i] + domain, To: strings.Replace(to, "=", "@", 1)})
	}

	return candidates
}

// ParseReplyTo returns the alias name an address is for and, for a reply-encoded address, the
// address it replies to. A reply-encoded address splits at the first "+", so the reply target
// may contain "+" ("a+b+c=example.com@domain.com" -> "a@domain.com", "b+c@example.com"), except
// a Wildcard Alias reply ("*+shop+c=example.com@domain.com" -> "*+shop@domain.com"). It can't
// tell an alias with "+" in its name apart, so FindRecipients tries every ReplyToCandidates split.
func ParseReplyTo(email string) (string, string) {
	if candidates := ReplyToCandidates(email); len(candidates) > 0 {
		reply := candidates[len(candidates)-1]
		// "*@domain.com" is never an alias, only the start of a Wildcard Alias name
		if strings.HasPrefix(reply.Alias, "*@") && len(candidates) > 1 {
			reply = candidates[len(candidates)-2]
		}
		return reply.Alias, reply.To
	}

	// Plain sub-address tag (e.g. "alias+tag@domain.com"): resolve against the base alias.
	atIndex := strings.Index(email, "@")
	plusIndex := strings.Index(email, "+")
	if plusIndex > 0 && atIndex > plusIndex+1 {
		return email[:plusIndex] + email[atIndex:], ""
	}

	return email, ""
}

// WildcardAliasForDelimiter returns the wildcard-suffix form of a tagged address for the
// given delimiter (e.g. "anything+suffix@domain.com" with "+" -> "*+suffix@domain.com"),
// used to match Wildcard Aliases when no exact alias exists for the tagged address. ok is
// false when email has no plain tag for that delimiter (none present, before "@", or one
// that looks reply-encoded).
func WildcardAliasForDelimiter(email string, delimiter string) (string, bool) {
	atIndex := strings.Index(email, "@")
	delimIndex := strings.Index(email, delimiter)
	if delimIndex == -1 || atIndex == -1 || delimIndex > atIndex {
		return "", false
	}

	suffix := email[delimIndex+len(delimiter) : atIndex]
	if suffix == "" || strings.Contains(suffix, "=") {
		return "", false
	}

	return "*" + email[delimIndex:], true
}

// WildcardAlias returns the wildcard-suffix form of a plus-tagged address
// (e.g. "anything+suffix@domain.com" -> "*+suffix@domain.com"), used to
// match Wildcard Aliases when no exact alias exists for the tagged address.
// ok is false when email has no plain "+" tag (none, or a reply-encoded one).
func WildcardAlias(email string) (string, bool) {
	return WildcardAliasForDelimiter(email, WildcardDelimiterPlus)
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

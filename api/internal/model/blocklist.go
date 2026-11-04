package model

import (
	"errors"
	"strings"
)

var (
	ErrBlocklistItemNotFound = errors.New("blocklist item not found")
)

type BlocklistType string

const (
	BlocklistEmail  BlocklistType = "email"
	BlocklistDomain BlocklistType = "domain"
)

// MaxBlocklistItems is the per-user cap on blocklist entries.
const MaxBlocklistItems = 1000

// BlocklistItem is a sender address or domain whose inbound mail is discarded for the user's
// aliases. The (user_id, type, value) unique index doubles as the lookup index used while
// processing inbound mail; its columns are sized explicitly because GORM only sizes string
// columns for plain index/unique tags, and MySQL can't index the longtext it falls back to.
type BlocklistItem struct {
	BaseModel
	UserID  string        `json:"-" gorm:"size:191;uniqueIndex:idx_blocklist_items_user_id_type_value,priority:1"`
	Type    BlocklistType `json:"type" gorm:"size:16;uniqueIndex:idx_blocklist_items_user_id_type_value,priority:2"`
	Value   string        `json:"value" gorm:"size:255;uniqueIndex:idx_blocklist_items_user_id_type_value,priority:3"`
	Blocked int           `json:"blocked" gorm:"default:0"`
}

type BlocklistList struct {
	Items []BlocklistItem `json:"items"`
	Total int             `json:"total"`
}

// NormalizeBlocklistValue brings a user-entered value into the form it is stored and matched in:
// trimmed, lowercased and, for domains, without a leading "@" (e.g. "@Example.com" -> "example.com").
func NormalizeBlocklistValue(itemType BlocklistType, value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if itemType == BlocklistDomain {
		value = strings.TrimPrefix(value, "@")
	}

	return value
}

// BlocklistValueRules returns the validator rules a normalized value of the given type must pass.
func BlocklistValueRules(itemType BlocklistType) string {
	if itemType == BlocklistEmail {
		return "required,max=254,emailx"
	}

	return "required,max=253,fqdn"
}

// BlocklistCandidates returns the values a sender address can match: the address itself for
// email items, and its domain plus every parent domain for domain items, so blocking
// "example.com" also blocks "mail.example.com".
func BlocklistCandidates(from string) (string, []string) {
	email := strings.ToLower(strings.TrimSpace(from))

	at := strings.LastIndex(email, "@")
	if at == -1 || at == len(email)-1 {
		return email, nil
	}

	domain := email[at+1:]
	domains := []string{domain}
	for {
		dot := strings.Index(domain, ".")
		if dot == -1 {
			break
		}
		domain = domain[dot+1:]
		if !strings.Contains(domain, ".") {
			break // stop above the TLD, which can't be stored as a blocklist domain
		}
		domains = append(domains, domain)
	}

	return email, domains
}

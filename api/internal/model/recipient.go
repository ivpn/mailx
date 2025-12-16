package model

import (
	"errors"
	"strings"
)

var (
	ErrDuplicateRecipient = errors.New("email already exists")
)

type Recipient struct {
	BaseModel
	UserID     string `json:"-"`
	Email      string `gorm:"unique" json:"email"`
	IsActive   bool   `json:"is_active"`
	PGPKey     string `json:"pgp_key"`
	PGPEnabled bool   `json:"pgp_enabled"`
	PGPInline  bool   `json:"pgp_inline"`
}

func GetEmails(rcps []Recipient) string {
	emails := []string{}
	for _, r := range rcps {
		emails = append(emails, r.Email)
	}

	return strings.Join(emails, ",")
}

func MergeCommaSeparatedEmails(a, b string) string {
	set := make(map[string]bool)

	for s := range strings.SplitSeq(a, ",") {
		if s != "" {
			set[s] = true
		}
	}
	for s := range strings.SplitSeq(b, ",") {
		if s != "" {
			set[s] = true
		}
	}

	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}

	return strings.Join(result, ",")
}

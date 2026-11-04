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

// RemoveCommaSeparatedEmail removes email from a comma-separated list, matching whole
// addresses case-insensitively, and reports whether it was present.
func RemoveCommaSeparatedEmail(list string, email string) (string, bool) {
	found := false
	result := []string{}
	for s := range strings.SplitSeq(list, ",") {
		if strings.EqualFold(s, email) {
			found = true
			continue
		}
		if s != "" {
			result = append(result, s)
		}
	}

	return strings.Join(result, ","), found
}

// ContainsCommaSeparatedEmail reports whether a comma-separated list contains email as a whole
// address, case-insensitively.
func ContainsCommaSeparatedEmail(list string, email string) bool {
	_, found := RemoveCommaSeparatedEmail(list, email)
	return found
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

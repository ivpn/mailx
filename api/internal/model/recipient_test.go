package model

import (
	"strings"
	"testing"
)

func TestGetEmails(t *testing.T) {
	tests := []struct {
		name       string
		recipients []Recipient
		expected   string
	}{
		{
			name:       "No recipients",
			recipients: []Recipient{},
			expected:   "",
		},
		{
			name: "Single recipient",
			recipients: []Recipient{
				{Email: "test1@example.com"},
			},
			expected: "test1@example.com",
		},
		{
			name: "Multiple recipients",
			recipients: []Recipient{
				{Email: "test1@example.com"},
				{Email: "test2@example.com"},
				{Email: "test3@example.com"},
			},
			expected: "test1@example.com,test2@example.com,test3@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEmails(tt.recipients)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMergeCommaSeparatedEmails(t *testing.T) {
	type tc struct {
		name        string
		a, b        string
		expected    []string
		expectExact bool
		exactStr    string
	}
	tests := []tc{
		{
			name:        "Both empty",
			a:           "",
			b:           "",
			expected:    []string{},
			expectExact: true,
			exactStr:    "",
		},
		{
			name:     "A empty B single",
			a:        "",
			b:        "one@example.com",
			expected: []string{"one@example.com"},
		},
		{
			name:     "B empty A single",
			a:        "one@example.com",
			b:        "",
			expected: []string{"one@example.com"},
		},
		{
			name:     "Duplicates in both",
			a:        "dup@example.com,dup@example.com",
			b:        "dup@example.com",
			expected: []string{"dup@example.com"},
		},
		{
			name:     "Overlap sets",
			a:        "a@example.com,b@example.com",
			b:        "b@example.com,c@example.com",
			expected: []string{"a@example.com", "b@example.com", "c@example.com"},
		},
		{
			name:     "Trailing and leading commas",
			a:        "a@example.com,",
			b:        ",b@example.com",
			expected: []string{"a@example.com", "b@example.com"},
		},
		{
			name:     "Whitespace treated as distinct",
			a:        "user@example.com",
			b:        " user@example.com",
			expected: []string{"user@example.com", " user@example.com"},
		},
		{
			name:     "Empty tokens ignored",
			a:        ",,x@example.com,,",
			b:        ",y@example.com,,",
			expected: []string{"x@example.com", "y@example.com"},
		},
		{
			name:     "Single duplicate across both",
			a:        "solo@example.com",
			b:        "solo@example.com",
			expected: []string{"solo@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeCommaSeparatedEmails(tt.a, tt.b)
			if tt.expectExact {
				if got != tt.exactStr {
					t.Fatalf("expected exact %q got %q", tt.exactStr, got)
				}
				return
			}
			gotSet := splitNonEmpty(got)
			if !equalSet(gotSet, tt.expected) {
				t.Fatalf("expected set %v got %v (raw: %q)", tt.expected, gotSet, got)
			}
		})
	}
}

func splitNonEmpty(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func equalSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		if m[v] == 0 {
			return false
		}
		m[v]--
		if m[v] < 0 {
			return false
		}
	}
	return true
}

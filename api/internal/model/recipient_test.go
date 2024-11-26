package model

import (
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

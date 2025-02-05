package model

import (
	"testing"
)

func TestGenerateReplyTo(t *testing.T) {
	tests := []struct {
		alias         string
		to            string
		expectedEmail string
	}{
		{
			alias:         "user@domain.com",
			to:            "reply@example.com",
			expectedEmail: "user+reply=example.com@domain.com",
		},
		{
			alias:         "user@domain.com",
			to:            "contact@another.com",
			expectedEmail: "user+contact=another.com@domain.com",
		},
		{
			alias:         "info@service.com",
			to:            "support@help.com",
			expectedEmail: "info+support=help.com@service.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.alias+"_"+tt.to, func(t *testing.T) {
			email := GenerateReplyTo(tt.alias, tt.to)
			if email != tt.expectedEmail {
				t.Errorf("expected email %s, got %s", tt.expectedEmail, email)
			}
		})
	}
}
func TestPlainTextToHTML(t *testing.T) {
	tests := []struct {
		plainText    string
		expectedHTML string
	}{
		{
			plainText:    "Hello, World!",
			expectedHTML: "Hello, World!",
		},
		{
			plainText:    "Line1\nLine2",
			expectedHTML: "Line1<br>Line2",
		},
		{
			plainText:    "<script>alert('XSS')</script>",
			expectedHTML: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
		{
			plainText:    "Hello\nWorld\n!",
			expectedHTML: "Hello<br>World<br>!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.plainText, func(t *testing.T) {
			html := PlainTextToHTML(tt.plainText)
			if html != tt.expectedHTML {
				t.Errorf("expected HTML %s, got %s", tt.expectedHTML, html)
			}
		})
	}
}

func TestParseReplyTo(t *testing.T) {
	tests := []struct {
		email         string
		expectedAlias string
		expectedRcp   string
	}{
		{
			email:         "user+reply=example.com@domain.com",
			expectedAlias: "user@domain.com",
			expectedRcp:   "reply@example.com",
		},
		{
			email:         "user@domain.com",
			expectedAlias: "user@domain.com",
			expectedRcp:   "",
		},
		{
			email:         "user+reply@domain.com",
			expectedAlias: "*+reply@domain.com",
			expectedRcp:   "",
		},
		{
			email:         "user+catchall+reply=example.com@domain.com",
			expectedAlias: "*+catchall@domain.com",
			expectedRcp:   "reply@example.com",
		},
		{
			email:         "user+catchall@domain.com",
			expectedAlias: "*+catchall@domain.com",
			expectedRcp:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			alias, rcp := ParseReplyTo(tt.email)
			if alias != tt.expectedAlias {
				t.Errorf("expected alias %s, got %s", tt.expectedAlias, alias)
			}
			if rcp != tt.expectedRcp {
				t.Errorf("expected rcp %s, got %s", tt.expectedRcp, rcp)
			}
		})
	}
}

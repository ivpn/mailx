package model

import (
	"net/mail"
	"strings"
	"testing"
)

func TestIsReply(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "reply message with In-Reply-To",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nIn-Reply-To: <message-id>\r\n\r\nThis is the body of the reply email.",
			want: true,
		},
		{
			name: "reply message with References",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nReferences: <message-id>\r\n\r\nThis is the body of the reply email.",
			want: true,
		},
		{
			name: "non-reply message",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email.",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := mail.ReadMessage(strings.NewReader(tt.data))
			if err != nil {
				t.Fatalf("failed to read message: %v", err)
			}
			if got := isReply(msg); got != tt.want {
				t.Errorf("isReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractDKIMDomains(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []string
		wantErr bool
	}{
		{
			name: "valid DKIM signature",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=test;\r\n" +
				"From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nMessage Body",
			want:    []string{"example.com"},
			wantErr: false,
		},
		{
			name: "multiple DKIM signatures",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=first.com; s=test;\r\n" +
				"DKIM-Signature: v=1; a=rsa-sha256; d=second.com; s=test;\r\n" +
				"From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nMessage Body",
			want:    []string{"first.com", "second.com"},
			wantErr: false,
		},
		{
			name:    "no DKIM signature",
			data:    "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: No DKIM\r\n\r\nMessage Body",
			want:    []string{},
			wantErr: false,
		},
		{
			name: "malformed DKIM signature - no domain",
			data: "DKIM-Signature: v=1; a=rsa-sha256; s=test;\r\n" +
				"From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nMessage Body",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid email format",
			data:    "This is not an email",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDKIMDomains([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("extractDKIMDomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("extractDKIMDomains() got = %v, want %v", got, tt.want)
					return
				}

				for i := range tt.want {
					if i < len(got) && got[i] != tt.want[i] {
						t.Errorf("extractDKIMDomains() got[%d] = %v, want[%d] = %v", i, got[i], i, tt.want[i])
					}
				}
			}
		})
	}
}

func TestParseMsgSuccess(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Msg
		wantErr bool
	}{
		{
			name: "successful parsing - simple message",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=test;\r\n" +
				"From: Sender Name <sender@example.com>\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: Test Subject\r\n\r\n" +
				"This is the body of the email.",
			want: Msg{
				From:     "sender@example.com",
				FromName: "Sender Name",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				Body:     "This is the body of the email.",
				Type:     Send,
			},
			wantErr: false,
		},
		{
			name: "successful parsing - reply message",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=test;\r\n" +
				"From: Sender Name <sender@example.com>\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: Re: Test Subject\r\n" +
				"In-Reply-To: <message-id>\r\n\r\n" +
				"This is a reply email.",
			want: Msg{
				From:     "sender@example.com",
				FromName: "Sender Name",
				To:       []string{"recipient@example.com"},
				Subject:  "Re: Test Subject",
				Body:     "This is a reply email.",
				Type:     Reply,
			},
			wantErr: false,
		},
		{
			name: "successful parsing - multiple recipients",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=test;\r\n" +
				"From: sender@example.com\r\n" +
				"To: recipient1@example.com, recipient2@example.com\r\n" +
				"Subject: Multi Recipients\r\n\r\n" +
				"This is an email with multiple recipients.",
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient1@example.com", "recipient2@example.com"},
				Subject:  "Multi Recipients",
				Body:     "This is an email with multiple recipients.",
				Type:     Send,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMsg([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !compareMessages(got, tt.want) {
					t.Errorf("ParseMsg() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestParseMsgErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "invalid To header",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=test;\r\n" +
				"From: sender@example.com\r\n" +
				"To: invalid@email@example.com\r\n" +
				"Subject: Invalid To\r\n\r\n" +
				"Message with invalid To header.",
			wantErr: true,
		},
		{
			name: "no valid DKIM signature",
			data: "From: sender@example.com\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: No DKIM\r\n\r\n" +
				"Email without DKIM signature.",
			wantErr: true,
		},
		{
			name: "DKIM domain mismatch",
			data: "DKIM-Signature: v=1; a=rsa-sha256; d=different.com; s=test;\r\n" +
				"From: sender@example.com\r\n" +
				"To: recipient@example.com\r\n" +
				"Subject: DKIM Mismatch\r\n\r\n" +
				"Email with DKIM domain mismatch.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMsg([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func compareMessages(a, b Msg) bool {
	if a.From != b.From || a.FromName != b.FromName || a.Subject != b.Subject || a.Body != b.Body || a.Type != b.Type {
		return false
	}
	if len(a.To) != len(b.To) {
		return false
	}
	for i := range a.To {
		if a.To[i] != b.To[i] {
			return false
		}
	}
	return true
}

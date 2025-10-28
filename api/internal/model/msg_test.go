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

func TestParseMessageError(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Msg
		wantErr bool
	}{
		{
			name: "valid message",
			data: []byte("From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				Body:     "This is the body of the email.",
				Type:     Send,
			},
			wantErr: true, // email authentication fails
		},
		{
			name: "valid reply message",
			data: []byte("From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nIn-Reply-To: <message-id>\r\n\r\nThis is the body of the reply email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient@example.com"},
				Subject:  "Re: Test Subject",
				Body:     "This is the body of the reply email.",
				Type:     Reply,
			},
			wantErr: true, // email authentication fails
		},
		{
			name:    "invalid message",
			data:    []byte("Invalid email data"),
			want:    Msg{},
			wantErr: true,
		},
		{
			name: "valid message with multiple recipients",
			data: []byte("From: sender@example.com\r\nTo: recipient1@example.com, recipient2@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient1@example.com", "recipient2@example.com"},
				Subject:  "Test Subject",
				Body:     "This is the body of the email.",
				Type:     Send,
			},
			wantErr: true, // email authentication fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMsg(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareMessages(got, tt.want) {
				t.Errorf("parseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBounce(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "bounce message with empty Return-Path",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Delivery Status Notification\r\nReturn-Path: <>\r\n\r\nThis is a bounce message.",
			want: true,
		},
		{
			name: "bounce message with multipart/report and delivery-status",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Delivery Status Notification\r\nContent-Type: multipart/report; report-type=delivery-status\r\n\r\nThis is a bounce message.",
			want: true,
		},
		{
			name: "bounce message with auto-replied",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Auto Reply\r\nAuto-Submitted: auto-replied\r\n\r\nThis is an auto-replied message.",
			want: true,
		},
		{
			name: "bounce message with case insensitive auto-replied",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Auto Reply\r\nAuto-Submitted: AUTO-REPLIED\r\n\r\nThis is an auto-replied message.",
			want: true,
		},
		{
			name: "bounce message with case insensitive multipart/report",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Delivery Status Notification\r\nContent-Type: MULTIPART/REPORT; REPORT-TYPE=DELIVERY-STATUS\r\n\r\nThis is a bounce message.",
			want: true,
		},
		{
			name: "non-bounce message with normal Return-Path",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\nReturn-Path: <sender@example.com>\r\n\r\nThis is a normal email.",
			want: false,
		},
		{
			name: "non-bounce message with different Content-Type",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\nContent-Type: text/plain\r\n\r\nThis is a normal email.",
			want: false,
		},
		{
			name: "non-bounce message with multipart/report but wrong report-type",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\nContent-Type: multipart/report; report-type=disposition-notification\r\n\r\nThis is not a bounce message.",
			want: false,
		},
		{
			name: "non-bounce message with different Auto-Submitted",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\nAuto-Submitted: auto-generated\r\n\r\nThis is not a bounce message.",
			want: false,
		},
		{
			name: "normal message without bounce headers",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\n\r\nThis is a normal email.",
			want: false,
		},
		{
			name: "message with invalid Content-Type",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\nContent-Type: invalid/content/type\r\n\r\nThis is a message with invalid content type.",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := mail.ReadMessage(strings.NewReader(tt.data))
			if err != nil {
				t.Fatalf("failed to read message: %v", err)
			}
			if got := isBounce(msg); got != tt.want {
				t.Errorf("isBounce() = %v, want %v", got, tt.want)
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

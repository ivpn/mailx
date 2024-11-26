package mailer

import (
	"testing"

	"gopkg.in/gomail.v2"
	"ivpn.net/email/api/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.SMTPClientConfig
		expected Mailer
	}{
		{
			name: "No user and password",
			cfg: config.SMTPClientConfig{
				Host:     "smtp.example.com",
				Port:     "587",
				Sender:   "sender@example.com",
				User:     "",
				Password: "",
			},
			expected: Mailer{
				dialer: &gomail.Dialer{Host: "smtp.example.com", Port: 587},
				Sender: "sender@example.com",
			},
		},
		{
			name: "With user and password",
			cfg: config.SMTPClientConfig{
				Host:     "smtp.example.com",
				Port:     "587",
				Sender:   "sender@example.com",
				User:     "user",
				Password: "password",
			},
			expected: Mailer{
				dialer: gomail.NewDialer("smtp.example.com", 587, "user", "password"),
				Sender: "sender@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := New(tt.cfg)
			if mailer.Sender != tt.expected.Sender {
				t.Errorf("expected sender %s, got %s", tt.expected.Sender, mailer.Sender)
			}
			if mailer.dialer.Host != tt.expected.dialer.Host {
				t.Errorf("expected host %s, got %s", tt.expected.dialer.Host, mailer.dialer.Host)
			}
			if mailer.dialer.Port != tt.expected.dialer.Port {
				t.Errorf("expected port %d, got %d", tt.expected.dialer.Port, mailer.dialer.Port)
			}
			if mailer.dialer.Username != tt.expected.dialer.Username {
				t.Errorf("expected username %s, got %s", tt.expected.dialer.Username, mailer.dialer.Username)
			}
			if mailer.dialer.Password != tt.expected.dialer.Password {
				t.Errorf("expected password %s, got %s", tt.expected.dialer.Password, mailer.dialer.Password)
			}
		})
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name    string
		mailer  Mailer
		to      string
		subject string
		body    string
		wantErr bool
	}{
		{
			name: "Send email with error",
			mailer: Mailer{
				dialer: &gomail.Dialer{Host: "invalid-host", Port: 587},
				Sender: "sender@example.com",
			},
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "Test Body",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mailer.Send(tt.to, tt.subject, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

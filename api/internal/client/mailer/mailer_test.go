package mailer

import (
	"testing"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/utils/gomail.v2"
)

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
func TestNew_InvalidPort(t *testing.T) {
	cfg := config.SMTPClientConfig{
		Host:     "smtp.example.com",
		Port:     "invalid",
		Sender:   "sender@example.com",
		User:     "",
		Password: "",
	}
	mailer := New(cfg)
	if mailer.dialer != nil {
		t.Errorf("expected dialer to be nil, got %v", mailer.dialer)
	}
	if mailer.Sender != cfg.Sender {
		t.Errorf("expected sender %s, got %s", cfg.Sender, mailer.Sender)
	}
}

func TestNew_MultipleHosts_Failure(t *testing.T) {
	cfg := config.SMTPClientConfig{
		Host:     "invalid-host1, invalid-host2",
		Port:     "587",
		Sender:   "sender@example.com",
		User:     "",
		Password: "",
	}
	mailer := New(cfg)
	if mailer.Sender != cfg.Sender {
		t.Errorf("expected sender %s, got %s", cfg.Sender, mailer.Sender)
	}
}

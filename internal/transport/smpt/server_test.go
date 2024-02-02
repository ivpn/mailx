package smpt

import (
	"testing"

	"ivpn.net/email-service/config"
)

var (
	cfg = config.Config{
		SMTPHost: "127.0.0.1",
		SMTPPort: "8025",
	}
)

func TestStartServer(t *testing.T) {
	err := Start(cfg)
	if err != nil {
		t.Errorf("Error starting server: %v", err)
	}
}

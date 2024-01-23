package server

import (
	"testing"

	"ivpn.net/email-service/services/inbound/config"
)

var (
	cfg = config.Config{
		ServerHost: "127.0.0.1",
		ServerPort: "8025",
	}
)

func TestStartServer(t *testing.T) {
	err := Start(cfg)
	if err != nil {
		t.Errorf("Error starting server: %v", err)
	}
}

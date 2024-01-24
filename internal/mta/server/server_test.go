package server

import (
	"testing"

	"ivpn.net/email-service/config"
)

var (
	cfg = config.Config{
		MTAHost: "127.0.0.1",
		MTAPort: "8025",
	}
)

func TestStartServer(t *testing.T) {
	err := Start(cfg)
	if err != nil {
		t.Errorf("Error starting server: %v", err)
	}
}

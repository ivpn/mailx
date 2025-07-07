package utils

import (
	"log"
	"os"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
	"ivpn.net/email/api/config"
)

func TestNewLogger(t *testing.T) {
	// Setup
	logFile := "test.log"
	cfg := config.APIConfig{
		LogFile: logFile,
	}

	// Call the function
	NewLogger(cfg)

	// Check if the log output is set correctly
	_, ok := log.Writer().(*lumberjack.Logger)
	if !ok {
		t.Fatalf("expected log writer to be *lumberjack.Logger, got %T", log.Writer())
	}

	// Cleanup
	os.Remove(logFile)
}

package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ivpn.net/email/api/config"
)

func TestSignupWebhookErrorIncludesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"subscription already used"}`))
	}))
	defer srv.Close()

	h := Http{Cfg: config.APIConfig{SignupWebhookURL: srv.URL}}

	err := h.SignupWebhook("sub-1")
	if err == nil {
		t.Fatal("Expected an error for a non-200 webhook response")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "subscription already used") {
		t.Errorf("Expected status and response body in error, got %q", err.Error())
	}
}

func TestSignupWebhookErrorTruncatesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(strings.Repeat("x", 10*maxLoggedResponse)))
	}))
	defer srv.Close()

	h := Http{Cfg: config.APIConfig{SignupWebhookURL: srv.URL}}

	err := h.SignupWebhook("sub-1")
	if err == nil {
		t.Fatal("Expected an error for a non-200 webhook response")
	}
	if strings.Count(err.Error(), "x") != maxLoggedResponse {
		t.Errorf("Expected response truncated to %d bytes, got %d", maxLoggedResponse, strings.Count(err.Error(), "x"))
	}
}

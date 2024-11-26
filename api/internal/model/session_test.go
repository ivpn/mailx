package model

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
)

func TestGenSessionToken(t *testing.T) {
	token, err := GenSessionToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(token) == 0 {
		t.Fatalf("expected a token, got an empty string")
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("expected a valid base64 token, got error: %v", err)
	}

	if len(decoded) != 32 {
		t.Fatalf("expected token length of 32 bytes, got %d bytes", len(decoded))
	}
}

func TestUnmarshalSessionData(t *testing.T) {
	sessionData := webauthn.SessionData{
		Challenge:        "challenge",
		UserID:           []byte("userID"),
		UserVerification: "required",
	}
	data, err := json.Marshal(sessionData)
	if err != nil {
		t.Fatalf("failed to marshal session data: %v", err)
	}

	session := &Session{
		Data: data,
	}

	err = session.UnmarshalSessionData()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(session.SessionData.Challenge) != "challenge" {
		t.Fatalf("expected challenge to be 'challenge', got %s", session.SessionData.Challenge)
	}

	if string(session.SessionData.UserID) != "userID" {
		t.Fatalf("expected userID to be 'userID', got %s", session.SessionData.UserID)
	}

	if session.SessionData.UserVerification != "required" {
		t.Fatalf("expected userVerification to be 'required', got %s", session.SessionData.UserVerification)
	}
}

func TestUnmarshalSessionData_InvalidData(t *testing.T) {
	session := &Session{
		Data: []byte("invalid data"),
	}

	err := session.UnmarshalSessionData()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

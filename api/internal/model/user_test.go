package model

import (
	"testing"
)

func TestSetPassword(t *testing.T) {
	u := User{}
	err := u.SetPassword("password")
	if err != nil {
		t.Errorf("SetPassword() failed: %v", err)
	}

	if u.PasswordHash == "" {
		t.Errorf("SetPassword() failed: password hash is empty")
	}

	if u.PasswordPlain != nil {
		t.Errorf("SetPassword() failed: password plain is not nil")
	}

	if u.PasswordHash == "password" {
		t.Errorf("SetPassword() failed: password hash is not hashed")
	}
}

func TestMatches(t *testing.T) {
	u := User{}
	err := u.SetPassword("password")
	if err != nil {
		t.Errorf("SetPassword() failed: %v", err)
	}

	if !u.Matches("password") {
		t.Errorf("Matches() failed: expected true, got false")
	}

	if u.Matches("wrongpassword") {
		t.Errorf("Matches() failed: expected false, got true")
	}
}

func TestIsTotpEnabled(t *testing.T) {
	u := User{}

	// Test when TOTP is not enabled
	u.TotpSecret = ""
	if u.IsTotpEnabled() {
		t.Errorf("IsTotpEnabled() failed: expected false, got true")
	}

	// Test when TOTP is enabled
	u.TotpSecret = "some-secret"
	if !u.IsTotpEnabled() {
		t.Errorf("IsTotpEnabled() failed: expected true, got false")
	}
}

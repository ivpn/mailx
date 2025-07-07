package utils

import (
	"testing"
)

func TestCreateOTP(t *testing.T) {
	otp, err := CreateOTP()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(otp.Secret) != 6 {
		t.Errorf("expected OTP secret length to be 6, got %d", len(otp.Secret))
	}

	if otp.Hash == "" {
		t.Errorf("expected OTP hash to be non-empty")
	}

	if !MatchOTP(otp.Secret, otp.Hash) {
		t.Errorf("expected OTP secret to match its hash")
	}
}

func TestCreateLongOTP(t *testing.T) {
	otp, err := CreateLongOTP()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(otp.Secret) != length {
		t.Errorf("expected OTP secret length to be %d, got %d", length, len(otp.Secret))
	}

	if otp.Hash == "" {
		t.Errorf("expected OTP hash to be non-empty")
	}

	if !MatchOTP(otp.Secret, otp.Hash) {
		t.Errorf("expected OTP secret to match its hash")
	}
}

func TestMatchOTP(t *testing.T) {
	otp, err := CreateOTP()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !MatchOTP(otp.Secret, otp.Hash) {
		t.Errorf("expected OTP secret to match its hash")
	}

	// Test with incorrect hash
	if MatchOTP(otp.Secret, "incorrectHash") {
		t.Errorf("expected OTP secret not to match an incorrect hash")
	}

	// Test with incorrect secret
	if MatchOTP("incorrectSecret", otp.Hash) {
		t.Errorf("expected incorrect secret not to match the hash")
	}
}

package utils

import (
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP" // This is a base32 encoded string

	// Generate a valid code for the current time
	t0 := int(time.Now().UTC().Unix() / 30)
	code, err := computeCode(secret, int64(t0))
	if err != nil {
		t.Fatalf("Failed to compute code: %v", err)
	}

	// Test with the valid code
	valid, err := Check(secret, code)
	if err != nil {
		t.Fatalf("Check returned an error: %v", err)
	}
	if !valid {
		t.Errorf("Expected code to be valid")
	}

	// Test with an invalid code
	invalidCode := (code + 1) % 1000000
	valid, err = Check(secret, invalidCode)
	if err != nil {
		t.Fatalf("Check returned an error: %v", err)
	}
	if valid {
		t.Errorf("Expected code to be invalid")
	}
}

func TestGenerateURI(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	user := "testuser"
	issuer := "testissuer"

	expectedURI := "otpauth://totp/testissuer:testuser?issuer=testissuer&secret=JBSWY3DPEHPK3PXP"
	uri := GenerateURI(secret, user, issuer)
	if uri != expectedURI {
		t.Errorf("Expected URI to be %s, but got %s", expectedURI, uri)
	}

	// Test without issuer
	expectedURI = "otpauth://totp/testuser?secret=JBSWY3DPEHPK3PXP"
	uri = GenerateURI(secret, user, "")
	if uri != expectedURI {
		t.Errorf("Expected URI to be %s, but got %s", expectedURI, uri)
	}
}

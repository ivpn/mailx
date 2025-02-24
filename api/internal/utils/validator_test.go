package utils

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()

	if v.Validate == nil {
		t.Error("expected validator to be initialized, but it was nil")
	}

	err := v.RegisterValidation("password", passwordValidation)
	if err != nil {
		t.Errorf("expected no error when registering password validation, but got: %v", err)
	}

	// Test if the custom password validation is registered correctly
	err = v.Var("ValidPassword1!", "password")
	if err != nil {
		t.Errorf("expected password to be valid, but got error: %v", err)
	}

	err = v.Var("short1!", "password")
	if err == nil {
		t.Error("expected password to be invalid due to length, but got no error")
	}

	err = v.Var("NoSpecialChar1", "password")
	if err == nil {
		t.Error("expected password to be invalid due to missing special character, but got no error")
	}

	err = v.RegisterValidation("pgp", pgpKeyValidation)
	if err != nil {
		t.Errorf("expected no error when registering pgp key validation, but got: %v", err)
	}

	// Test if the custom PGP key validation is registered correctly
	err = v.Var("-----BEGIN PGP PUBLIC KEY BLOCK-----", "pgp")
	if err != nil {
		t.Errorf("expected PGP key to be valid, but got error: %v", err)
	}

	err = v.Var("invalid-key", "pgp")
	if err == nil {
		t.Error("expected PGP key to be invalid, but got no error")
	}

	err = v.Var("", "pgp")
	if err != nil {
		t.Errorf("expected empty PGP key to be valid, but got error: %v", err)
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"", true},
		{"another.test@domain.co", false},
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
		}
	}
}

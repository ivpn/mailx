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
	err = v.Var("-----BEGIN PGP PUBLIC KEY BLOCK----- ... -----END PGP PUBLIC KEY BLOCK-----", "pgp")
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
		email string
		valid bool
		desc  string
	}{
		{"test@example.com", true, "valid email"},
		{"invalid-email", false, "missing @ symbol"},
		{"", false, "empty email"},
		{"another.test@domain.co", true, "valid email with subdomain"},
		{"\"test<svg/onload=alert()/>\"@domain.net", false, "XSS attempt in email"},
		{"\"<img src=x onerror=alert(\"1\")\"@domain.net>", false, "XSS attempt in email"},
		{"\"<svg/onload=alert(\"1\")>\"@domain.net", false, "XSS attempt in email"},
		{"\"<iframe src=\"javascript:alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<body onload=alert('XSS')>\"@domain.net", false, "XSS attempt in email"},
		{"\"<a href=\"javascript:alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<img src=\"x\" onerror=\"alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<script src=\"http://example.com/xss.js\"\"@domain.net></script>", false, "XSS attempt in email"},
		{"\"<link rel=\"stylesheet\" href=\"http://example.com/xss.css\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<meta http-equiv=\"refresh\" content=\"0;url=http://example.com/xss.html\">\"@domain.net", false, "XSS attempt in email"},
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		isValid := err == nil
		if isValid != tt.valid {
			t.Errorf("ValidateEmail(%q): got %v, want %v (%s)", tt.email, isValid, tt.valid, tt.desc)
		}
	}
}
func TestSqlEmailValidation(t *testing.T) {
	v := NewValidator()

	// Register the emailx validator if not already registered
	err := v.RegisterValidation("emailx", sqlEmailValidation)
	if err != nil {
		t.Errorf("expected no error when registering emailx validation, but got: %v", err)
	}

	// Test cases
	tests := []struct {
		email string
		valid bool
		desc  string
	}{
		{"test@example.com", true, "valid email"},
		{"user.name@domain.co.uk", true, "valid email with dots and subdomain"},
		{"user-name@domain.com", true, "valid email with hyphen"},
		{"user+tag@domain.com", true, "valid email with plus"},
		{"123@domain.com", true, "valid email with numbers"},
		{"", true, "empty email should be valid because of omitempty check"},
		{"invalid-email", false, "missing @ symbol"},
		{"user@", false, "missing domain"},
		{"user@domain", false, "missing TLD"},
		{"user@.com", false, "missing domain name"},
		{"@domain.com", false, "missing local part"},
		{"user@domain.", false, "TLD can't be empty"},
		{"user@domain.c", false, "TLD too short"},
		{"user space@domain.com", false, "spaces not allowed"},
		{"user!@domain.com", false, "invalid character in local part"},
		{"test<svg/onload=alert(\"1\")/>;@domain.net", false, "SVG tag injection attempt"},
		{"test@domain.com\" OR \"1\"=\"1\";--", false, "SQL injection with quotes"},
		{"test@domain.com\" UNION SELECT username,password FROM users;--", false, "SQL UNION injection"},
		{"test\") DROP TABLE users;--@domain.com", false, "DROP TABLE injection"},
		{"test\")) DELETE FROM emails;--@domain.com", false, "DELETE injection"},
		{"admin\"--@example.com", false, "Comment injection"},
		{"\"<script>alert(\"1\")</script>\"@domain.net", false, "XSS attempt in email"},
		{"\"<img src=x onerror=alert(\"1\")>\"@domain.net", false, "XSS attempt in email"},
		{"\"<svg/onload=alert(\"1\")>\"@domain.net", false, "XSS attempt in email"},
		{"\"<iframe src=\"javascript:alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<body onload=alert('XSS')>\"@domain.net", false, "XSS attempt in email"},
		{"\"<a href=\"javascript:alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<img src=\"x\" onerror=\"alert('XSS')\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<script src=\"http://example.com/xss.js\"></script>\"@domain.net", false, "XSS attempt in email"},
		{"\"<link rel=\"stylesheet\" href=\"http://example.com/xss.css\">\"@domain.net", false, "XSS attempt in email"},
		{"\"<meta http-equiv=\"refresh\" content=\"0;url=http://example.com/xss.html\">\"@domain.net", false, "XSS attempt in email"},
	}

	for _, tt := range tests {
		err := v.Var(tt.email, "emailx")
		isValid := err == nil
		if isValid != tt.valid {
			t.Errorf("emailxValidation(%q): got %v, want %v (%s)", tt.email, isValid, tt.valid, tt.desc)
		}
	}
}

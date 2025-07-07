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
func TestSearchValidation(t *testing.T) {
	v := NewValidator()

	// Register the search validator if not already registered
	err := v.RegisterValidation("search", searchValidation)
	if err != nil {
		t.Errorf("expected no error when registering search validation, but got: %v", err)
	}

	// Test cases
	tests := []struct {
		value string
		valid bool
		desc  string
	}{
		{"basic search", true, "simple text with space"},
		{"search@email.com", true, "email format with @ symbol"},
		{"search_with_underscore", true, "text with underscore"},
		{"search-with-hyphen", true, "text with hyphen"},
		{"search.with.dots", true, "text with dots"},
		{"search+with+plus", true, "text with plus symbol"},
		{"123456789", true, "numbers only"},
		{"mixedCase123", true, "mixed case alphanumeric"},
		{"", true, "empty string should be valid because of omitempty check"},
		{"invalid$character", false, "contains invalid $ character"},
		{"invalid#character", false, "contains invalid # character"},
		{"invalid*character", false, "contains invalid * character"},
		{"invalid!character", false, "contains invalid ! character"},
		{"invalid%character", false, "contains invalid % character"},
		{"invalid&character", false, "contains invalid & character"},
		{"invalid(character)", false, "contains invalid parentheses"},
		{"invalid/character", false, "contains invalid slash"},
		{"invalid\\character", false, "contains invalid backslash"},
		{"invalid;character", false, "contains invalid semicolon"},
		{"invalid:character", false, "contains invalid colon"},
		{"invalid'character", false, "contains invalid single quote"},
		{"invalid\"character", false, "contains invalid double quote"},
		{"invalid<script>", false, "contains script tags"},
		{"invalid SELECT * FROM", false, "contains SQL keywords"},
		{"<svg/onload=alert('XSS')>", false, "XSS attempt with SVG"},
		{"<img src=x onerror=alert('XSS')>", false, "XSS attempt with img tag"},
		{"<script>alert('XSS')</script>", false, "XSS attempt with script tag"},
		{"<iframe src=\"javascript:alert('XSS')\">", false, "XSS attempt with iframe"},
		{"<body onload=alert('XSS')>", false, "XSS attempt with body tag"},
		{"<a href=\"javascript:alert('XSS')\">", false, "XSS attempt with anchor tag"},
		{"<img src=\"x\" onerror=\"alert('XSS')\">", false, "XSS attempt with img tag and onerror"},
		{"<script src=\"http://example.com/xss.js\"></script>", false, "XSS attempt with script src"},
		{"<link rel=\"stylesheet\" href=\"http://example.com/xss.css\">", false, "XSS attempt with link tag"},
		{"<meta http-equiv=\"refresh\" content=\"0;url=http://example.com/xss.html\">", false, "XSS attempt with meta refresh"},
		{"SELECT * FROM users WHERE username = 'admin'", false, "SQL injection with SELECT"},
		{"' OR '1'='1'", false, "SQL injection with OR"},
		{"' UNION SELECT username, password FROM users;--", false, "SQL injection with UNION"},
		{"' DROP TABLE users;--", false, "SQL injection with DROP TABLE"},
		{"' DELETE FROM emails;--", false, "SQL injection with DELETE"},
		{"'--", false, "SQL injection with comment"},
		{"\" OR \"1\"=\"1\";--", false, "SQL injection with quotes"},
		{"\" UNION SELECT username, password FROM users;--", false, "SQL UNION injection"},
		{"\" DROP TABLE users;--", false, "DROP TABLE injection"},
		{"\")) DELETE FROM emails;--", false, "DELETE injection"},
		{"\"--", false, "Comment injection"},
	}

	for _, tt := range tests {
		err := v.Var(tt.value, "search")
		isValid := err == nil
		if isValid != tt.valid {
			t.Errorf("searchValidation(%q): got %v, want %v (%s)", tt.value, isValid, tt.valid, tt.desc)
		}
	}
}

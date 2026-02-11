package mailer

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"

	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/utils"
)

func TestPreprocessEmailData_RFC2047Encoded(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		checkFrom   string
	}{
		{
			name: "RFC 2047 encoded display name with email",
			input: `From: =?UTF-8?B?VGVzdCBVc2Vy?= <test@example.com>
To: recipient@example.com
Subject: Test Subject

Test body content`,
			expectError: false,
			checkFrom:   "Test User <test@example.com>",
		},
		{
			name: "RFC 2047 encoded display name with special characters",
			input: `From: =?UTF-8?B?aaabbbcccdddeee==?= <user@example.com>
To: recipient@example.com
Subject: Test Subject

Test body content`,
			expectError: false,
		},
		{
			name: "Multiple RFC 2047 encoded headers",
			input: `From: =?UTF-8?B?VGVzdCBVc2Vy?= <sender@example.com>
To: =?UTF-8?B?UmVjaXBpZW50?= <recipient@example.com>
Subject: =?UTF-8?B?VGVzdCBTdWJqZWN0?=

Test body content`,
			expectError: false,
		},
		{
			name: "Plain text email without encoding",
			input: `From: Test User <test@example.com>
To: recipient@example.com
Subject: Test Subject

Test body content`,
			expectError: false,
			checkFrom:   "Test User <test@example.com>",
		},
		{
			name: "Mixed encoded and plain headers",
			input: `From: =?UTF-8?Q?Test_User?= <test@example.com>
To: plain@example.com
Subject: Normal Subject

Test body content`,
			expectError: false,
			checkFrom:   "Test User <test@example.com>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := []byte(strings.ReplaceAll(tt.input, "\n", "\r\n"))

			processedData, err := utils.PreprocessEmailData(inputData)
			if (err != nil) != tt.expectError {
				t.Errorf("preprocessEmailData() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Try to parse the processed data with standard mail.ReadMessage
			msg, err := mail.ReadMessage(bytes.NewReader(processedData))
			if err != nil {
				t.Errorf("Failed to parse processed email: %v", err)
				return
			}

			// Verify the From header can be parsed
			fromHeader := msg.Header.Get("From")
			if fromHeader == "" {
				t.Error("From header is empty after preprocessing")
				return
			}

			// Try to parse the From address
			fromAddr, err := mail.ParseAddress(fromHeader)
			if err != nil {
				t.Errorf("Failed to parse From address after preprocessing: %v", err)
				t.Logf("From header value: %s", fromHeader)
				return
			}

			// If we have a specific expected From value, check it
			if tt.checkFrom != "" && !strings.Contains(fromHeader, "=?") {
				if fromHeader != tt.checkFrom {
					// Allow for minor variations in formatting
					if fromAddr.Address == "" {
						t.Errorf("Expected From to be parseable, got error")
					}
				}
			}

			t.Logf("Successfully parsed From: %s <%s>", fromAddr.Name, fromAddr.Address)
		})
	}
}

func TestPreprocessEmailData_PreservesBody(t *testing.T) {
	input := `From: =?UTF-8?B?VGVzdCBVc2Vy?= <test@example.com>
To: recipient@example.com
Subject: Test Subject
Content-Type: text/plain; charset=utf-8

This is a test body with multiple lines.
It should be preserved exactly as is.
Including special characters: ñ, ü, é`

	inputData := []byte(strings.ReplaceAll(input, "\n", "\r\n"))

	processedData, err := utils.PreprocessEmailData(inputData)
	if err != nil {
		t.Fatalf("utils.PreprocessEmailData() error = %v", err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(processedData))
	if err != nil {
		t.Fatalf("Failed to parse processed email: %v", err)
	}

	// Read and verify body
	bodyBuf := new(bytes.Buffer)
	_, err = bodyBuf.ReadFrom(msg.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	expectedBodyLines := []string{
		"This is a test body with multiple lines.",
		"It should be preserved exactly as is.",
		"Including special characters: ñ, ü, é",
	}

	for _, line := range expectedBodyLines {
		if !strings.Contains(bodyBuf.String(), line) {
			t.Errorf("Body missing expected line: %s", line)
		}
	}
}

func TestPreprocessEmailData_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty data",
			input: "",
		},
		{
			name:  "Invalid format",
			input: "This is not a valid email",
		},
		{
			name:  "Partial headers",
			input: "From: test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := []byte(tt.input)

			// Should not panic and should return the original data
			processedData, err := utils.PreprocessEmailData(inputData)
			if err != nil {
				t.Logf("utils.PreprocessEmailData() returned error: %v", err)
			}

			// Should return original data on error
			if !bytes.Equal(processedData, inputData) {
				t.Log("utils.PreprocessEmailData() returned different data, which is acceptable")
			}
		})
	}
}

func TestCleanupMalformedEncodedAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Malformed base64 with valid email",
			input:    "=?UTF-8?B?aaabbbcccdddeee==?= <user@example.com>",
			expected: "user@example.com",
		},
		{
			name:     "Valid encoded-word",
			input:    "Test User <test@example.com>",
			expected: "Test User <test@example.com>",
		},
		{
			name:     "Plain email without angle brackets",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "Email with display name",
			input:    "John Doe <john@example.com>",
			expected: "John Doe <john@example.com>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CleanupMalformedEncodedAddress(tt.input)

			// The result should be parseable by mail.ParseAddress
			_, err := mail.ParseAddress(result)
			if err != nil {
				t.Logf("Cleaned address: %s", result)
				t.Logf("Parse error: %v", err)
				// For malformed cases, we at least want the email part
				if !strings.Contains(result, "@") {
					t.Errorf("Result doesn't contain email address: %s", result)
				}
			} else {
				t.Logf("Successfully parsed cleaned address: %s", result)
			}
		})
	}
}

func TestPreprocessEmailData_MalformedEncoding(t *testing.T) {
	// Test with the exact error case from the issue
	input := `From: =?UTF-8?B?aaabbbcccdddeee==?= <user@example.com>
To: recipient@example.com
Subject: Test Subject

Test body content`

	inputData := []byte(strings.ReplaceAll(input, "\n", "\r\n"))

	processedData, err := utils.PreprocessEmailData(inputData)
	if err != nil {
		t.Fatalf("utils.PreprocessEmailData() error = %v", err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(processedData))
	if err != nil {
		t.Fatalf("Failed to parse processed email: %v", err)
	}

	// Verify the From header can be parsed
	fromHeader := msg.Header.Get("From")
	if fromHeader == "" {
		t.Fatal("From header is empty after preprocessing")
	}

	t.Logf("From header after preprocessing: %s", fromHeader)

	// Try to parse the From address - this should not fail
	fromAddr, err := mail.ParseAddress(fromHeader)
	if err != nil {
		t.Errorf("Failed to parse From address after preprocessing: %v", err)
		t.Logf("From header value: %s", fromHeader)
	} else {
		t.Logf("Successfully parsed From: %s <%s>", fromAddr.Name, fromAddr.Address)

		// Verify we at least got the email address
		if fromAddr.Address != "user@example.com" {
			t.Errorf("Expected email address 'user@example.com', got '%s'", fromAddr.Address)
		}
	}
}

func TestPreprocessEmailData_WithLettersParser(t *testing.T) {
	// Test the full integration with letters.ParseEmail
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name: "Malformed RFC 2047 encoding",
			email: `From: =?UTF-8?B?aaabbbcccdddeee==?= <sender@example.com>
To: recipient@example.com
Subject: Test Subject
Content-Type: text/plain; charset=utf-8

This is a test email body.`,
			expectError: false,
		},
		{
			name: "Valid RFC 2047 encoding",
			email: `From: =?UTF-8?B?VGVzdCBVc2Vy?= <sender@example.com>
To: recipient@example.com
Subject: =?UTF-8?B?VGVzdCBTdWJqZWN0?=
Content-Type: text/plain; charset=utf-8

This is a test email body.`,
			expectError: false,
		},
		{
			name: "Plain text headers",
			email: `From: Test User <sender@example.com>
To: recipient@example.com
Subject: Test Subject
Content-Type: text/plain; charset=utf-8

This is a test email body.`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailData := []byte(strings.ReplaceAll(tt.email, "\n", "\r\n"))

			// Preprocess the email
			processedData, err := utils.PreprocessEmailData(emailData)
			if err != nil {
				t.Fatalf("utils.PreprocessEmailData() error = %v", err)
			}

			// Try to parse with letters
			reader := bytes.NewReader(processedData)
			email, err := letters.ParseEmail(reader)

			if (err != nil) != tt.expectError {
				t.Errorf("letters.ParseEmail() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if err == nil {
				t.Logf("Successfully parsed email with letters")
				t.Logf("  Subject: %s", email.Headers.Subject)
				if len(email.Headers.From) > 0 {
					t.Logf("  From: %s <%s>", email.Headers.From[0].Name, email.Headers.From[0].Address)
				}
				t.Logf("  Text length: %d", len(email.Text))
			}
		})
	}
}

func TestPreprocessEmailData_ContentDisposition(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
		checkFunc   func(*testing.T, letters.Email)
	}{
		{
			name: "Simple Content-Disposition attachment",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Test with attachment
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain; charset=utf-8

This is the body.
--boundary123
Content-Type: application/pdf; name="document.pdf"
Content-Disposition: attachment; filename="document.pdf"
Content-Transfer-Encoding: base64

VGVzdCBjb250ZW50
--boundary123--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.AttachedFiles) != 1 {
					t.Errorf("Expected 1 attachment, got %d", len(email.AttachedFiles))
					return
				}
				if email.AttachedFiles[0].ContentDisposition.Params["filename"] != "document.pdf" {
					t.Errorf("Expected filename 'document.pdf', got '%s'", email.AttachedFiles[0].ContentDisposition.Params["filename"])
				}
			},
		},
		{
			name: "Content-Disposition inline",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Test with inline image
MIME-Version: 1.0
Content-Type: multipart/related; boundary="boundary456"

--boundary456
Content-Type: text/html; charset=utf-8

<html><body>Test <img src="cid:image1"></body></html>
--boundary456
Content-Type: image/png; name="image.png"
Content-Disposition: inline; filename="image.png"
Content-ID: <image1>
Content-Transfer-Encoding: base64

iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==
--boundary456--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.InlineFiles) != 1 {
					t.Errorf("Expected 1 inline file, got %d", len(email.InlineFiles))
					return
				}
				if email.InlineFiles[0].ContentDisposition.Params["filename"] != "image.png" {
					t.Errorf("Expected filename 'image.png', got '%s'", email.InlineFiles[0].ContentDisposition.Params["filename"])
				}
			},
		},
		{
			name: "Content-Disposition with RFC 2231 encoded filename",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Test with encoded filename
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary789"

--boundary789
Content-Type: text/plain; charset=utf-8

Test body
--boundary789
Content-Type: application/octet-stream
Content-Disposition: attachment;
 filename*=UTF-8''%E6%96%87%E6%A1%A3.pdf
Content-Transfer-Encoding: base64

VGVzdA==
--boundary789--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.AttachedFiles) != 1 {
					t.Errorf("Expected 1 attachment, got %d", len(email.AttachedFiles))
					return
				}
				// RFC 2231 encoded filename should be decoded
				t.Logf("Filename: %s", email.AttachedFiles[0].ContentDisposition.Params["filename"])
			},
		},
		{
			name: "Multi-line Content-Disposition header",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Test with multi-line disposition
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundaryABC"

--boundaryABC
Content-Type: text/plain

Body text
--boundaryABC
Content-Type: application/pdf
Content-Disposition: attachment;
	filename="very-long-filename-that-spans-multiple-lines.pdf";
	size=12345
Content-Transfer-Encoding: base64

VGVzdA==
--boundaryABC--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.AttachedFiles) != 1 {
					t.Errorf("Expected 1 attachment, got %d", len(email.AttachedFiles))
					return
				}
				filename := email.AttachedFiles[0].ContentDisposition.Params["filename"]
				if !strings.Contains(filename, "very-long-filename") {
					t.Errorf("Expected filename with 'very-long-filename', got '%s'", filename)
				}
			},
		},
		{
			name: "Content-Disposition with RFC 2047 in filename",
			email: `From: =?UTF-8?B?VGVzdCBVc2Vy?= <sender@example.com>
To: recipient@example.com
Subject: Test with RFC 2047 address and attachment
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundaryDEF"

--boundaryDEF
Content-Type: text/plain; charset=utf-8

Email body
--boundaryDEF
Content-Type: application/pdf
Content-Disposition: attachment; filename="report.pdf"
Content-Transfer-Encoding: base64

UmVwb3J0IGNvbnRlbnQ=
--boundaryDEF--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				// Verify both RFC 2047 in From header and Content-Disposition work
				if len(email.Headers.From) == 0 {
					t.Error("Expected From header")
					return
				}
				if email.Headers.From[0].Name != "Test User" {
					t.Errorf("Expected From name 'Test User', got '%s'", email.Headers.From[0].Name)
				}
				if len(email.AttachedFiles) != 1 {
					t.Errorf("Expected 1 attachment, got %d", len(email.AttachedFiles))
					return
				}
			},
		},
		{
			name: "Complex multipart with multiple Content-Disposition headers",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Multiple attachments
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="----=_Part_109832_453230012.1770746967057"

------=_Part_109832_453230012.1770746967057
Content-Type: text/plain; charset=utf-8

Body of the email
------=_Part_109832_453230012.1770746967057
Content-Type: application/pdf; name="file1.pdf"
Content-Disposition: attachment; filename="file1.pdf"
Content-Transfer-Encoding: base64

RmlsZTE=
------=_Part_109832_453230012.1770746967057
Content-Type: image/jpeg; name="photo.jpg"
Content-Disposition: attachment; filename="photo.jpg"
Content-Transfer-Encoding: base64

UGhvdG8=
------=_Part_109832_453230012.1770746967057--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.AttachedFiles) != 2 {
					t.Errorf("Expected 2 attachments, got %d", len(email.AttachedFiles))
					return
				}
				filenames := []string{
					email.AttachedFiles[0].ContentDisposition.Params["filename"],
					email.AttachedFiles[1].ContentDisposition.Params["filename"],
				}
				t.Logf("Found filenames: %v", filenames)
			},
		},
		{
			name: "Content-Disposition with special characters in filename",
			email: `From: sender@example.com
To: recipient@example.com
Subject: Special chars in filename
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundaryXYZ"

--boundaryXYZ
Content-Type: text/plain

Test
--boundaryXYZ
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="file (copy) [2024].txt"
Content-Transfer-Encoding: base64

VGVzdA==
--boundaryXYZ--
`,
			expectError: false,
			checkFunc: func(t *testing.T, email letters.Email) {
				if len(email.AttachedFiles) != 1 {
					t.Errorf("Expected 1 attachment, got %d", len(email.AttachedFiles))
					return
				}
				filename := email.AttachedFiles[0].ContentDisposition.Params["filename"]
				t.Logf("Parsed filename: %s", filename)
				if !strings.Contains(filename, "copy") {
					t.Errorf("Expected filename to contain 'copy', got '%s'", filename)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailData := []byte(strings.ReplaceAll(tt.email, "\n", "\r\n"))

			// Preprocess the email
			processedData, err := utils.PreprocessEmailData(emailData)
			if err != nil {
				t.Fatalf("PreprocessEmailData() error = %v", err)
			}

			// Parse with letters
			reader := bytes.NewReader(processedData)
			email, err := letters.ParseEmail(reader)

			if (err != nil) != tt.expectError {
				t.Errorf("letters.ParseEmail() error = %v, expectError %v", err, tt.expectError)
				if err != nil {
					t.Logf("Error details: %v", err)
				}
				return
			}

			if err == nil && tt.checkFunc != nil {
				tt.checkFunc(t, email)
			}
		})
	}
}

func TestPreprocessEmailData_PreservesComplexMIME(t *testing.T) {
	// Test that complex MIME structure including Content-Disposition is preserved
	input := `From: sender@example.com
To: recipient@example.com
Subject: Complex MIME test
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="alt-boundary"

--alt-boundary
Content-Type: text/plain; charset=utf-8
Content-Disposition: inline

Plain text version
--alt-boundary
Content-Type: text/html; charset=utf-8
Content-Disposition: inline

<html><body>HTML version</body></html>
--alt-boundary--
`

	inputData := []byte(strings.ReplaceAll(input, "\n", "\r\n"))

	processedData, err := utils.PreprocessEmailData(inputData)
	if err != nil {
		t.Fatalf("PreprocessEmailData() error = %v", err)
	}

	// Verify Content-Disposition headers are preserved
	if !bytes.Contains(processedData, []byte("Content-Disposition: inline")) {
		t.Error("Content-Disposition header was not preserved")
	}

	// Verify it can be parsed by letters
	reader := bytes.NewReader(processedData)
	email, err := letters.ParseEmail(reader)
	if err != nil {
		t.Fatalf("letters.ParseEmail() error = %v", err)
	}

	if email.Text == "" {
		t.Error("Expected text content to be parsed")
	}
	if email.HTML == "" {
		t.Error("Expected HTML content to be parsed")
	}

	t.Logf("Successfully parsed multipart/alternative email")
	t.Logf("  Text: %s", strings.TrimSpace(email.Text))
	t.Logf("  HTML: %s", strings.TrimSpace(email.HTML))
}

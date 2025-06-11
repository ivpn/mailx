package model

import (
	"log"
	"net/mail"
	"strings"
	"testing"
)

func TestIsReply(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "reply message with In-Reply-To",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nIn-Reply-To: <message-id>\r\n\r\nThis is the body of the reply email.",
			want: true,
		},
		{
			name: "reply message with References",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nReferences: <message-id>\r\n\r\nThis is the body of the reply email.",
			want: true,
		},
		{
			name: "non-reply message",
			data: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email.",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := mail.ReadMessage(strings.NewReader(tt.data))
			if err != nil {
				t.Fatalf("failed to read message: %v", err)
			}
			if got := isReply(msg); got != tt.want {
				t.Errorf("isReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMessageError(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Msg
		wantErr bool
	}{
		{
			name: "valid message",
			data: []byte("From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				Body:     "This is the body of the email.",
				Type:     Send,
			},
			wantErr: true, // email authentication fails
		},
		{
			name: "valid reply message",
			data: []byte("From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test Subject\r\nIn-Reply-To: <message-id>\r\n\r\nThis is the body of the reply email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient@example.com"},
				Subject:  "Re: Test Subject",
				Body:     "This is the body of the reply email.",
				Type:     Reply,
			},
			wantErr: true, // email authentication fails
		},
		{
			name:    "invalid message",
			data:    []byte("Invalid email data"),
			want:    Msg{},
			wantErr: true,
		},
		{
			name: "valid message with multiple recipients",
			data: []byte("From: sender@example.com\r\nTo: recipient1@example.com, recipient2@example.com\r\nSubject: Test Subject\r\n\r\nThis is the body of the email."),
			want: Msg{
				From:     "sender@example.com",
				FromName: "",
				To:       []string{"recipient1@example.com", "recipient2@example.com"},
				Subject:  "Test Subject",
				Body:     "This is the body of the email.",
				Type:     Send,
			},
			wantErr: true, // email authentication fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMsg(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareMessages(got, tt.want) {
				t.Errorf("parseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPGPSignatures(t *testing.T) {
	tests := []struct {
		name          string
		data          string
		wantNum       int
		wantFilenames []string
		wantErr       bool
	}{
		{
			name: "email with pgp signature",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Signed Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a signed email.

--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="signature.asc"

-----BEGIN PGP SIGNATURE-----
Version: Example
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8One
-----END PGP SIGNATURE-----
--boundary123--`,
			wantNum:       1,
			wantFilenames: []string{"signature.asc"},
			wantErr:       false,
		},
		{
			name: "email with multiple pgp signatures",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Multiple Signatures
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This email has multiple signatures.

--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="sig1.asc"

-----BEGIN PGP SIGNATURE-----
Version: Example
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8One
-----END PGP SIGNATURE-----
--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="sig2.asc"

-----BEGIN PGP SIGNATURE-----
Version: Example2
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8Two
-----END PGP SIGNATURE-----
--boundary123--`,
			wantNum:       2,
			wantFilenames: []string{"sig1.asc", "sig2.asc"},
			wantErr:       false,
		},
		{
			name: "email with signature but no filename",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Signed Email No Filename
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a signed email.

--boundary123
Content-Type: application/pgp-signature

-----BEGIN PGP SIGNATURE-----
Version: Example
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8One
-----END PGP SIGNATURE-----
--boundary123--`,
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name: "email without pgp signature",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Regular Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a regular email without signature.
--boundary123--`,
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name: "non-multipart email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Plain Email
Content-Type: text/plain

This is just a plain text email.`,
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name:          "invalid email data",
			data:          "Invalid email data",
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       true,
		},
		{
			name: "email with both asc and non-asc attachments",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Mixed Attachments
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Email with both .asc and non-.asc attachments.

--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="signature.asc"

-----BEGIN PGP SIGNATURE-----
Version: Example
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8One
-----END PGP SIGNATURE-----
--boundary123
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="document.pdf"

Sample PDF Content
--boundary123
Content-Type: application/pgp-signature 
Content-Disposition: attachment; filename="other.txt"

-----BEGIN PGP SIGNATURE-----
Version: Example
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCAAdFiEE+Y5JJsjFlnUSqMJJNnn76HnlCeEFAmVtZIUACgkQNnn76Hnl
CeHk9Qf9Eq4shrink7GFh75J7qbgbPHgbRhVuTrCGLeVIKbgDCURDjB2YJx5dA==
=s8One
-----END PGP SIGNATURE-----
--boundary123--`,
			wantNum:       1,
			wantFilenames: []string{"signature.asc"},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPGPSignatures([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPGPSignatures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("ExtractPGPSignatures() got %d signatures", len(got))
			for i, attachment := range got {
				log.Printf("Signature %d: Filename=%q, ContentType=%q", i, attachment.Filename, attachment.ContentType)
			}

			if !tt.wantErr {
				if len(got) != tt.wantNum {
					t.Errorf("ExtractPGPSignatures() returned %d signatures, want %d", len(got), tt.wantNum)
					return
				}

				for i, attachment := range got {
					if i < len(tt.wantFilenames) && attachment.Filename != tt.wantFilenames[i] {
						t.Errorf("ExtractPGPSignatures() signature %d filename = %q, want %q",
							i, attachment.Filename, tt.wantFilenames[i])
					}
					if !strings.HasPrefix(attachment.ContentType, "application/pgp-signature") {
						t.Errorf("ExtractPGPSignatures() signature %d ContentType = %q, want prefix %q",
							i, attachment.ContentType, "application/pgp-signature")
					}
					if attachment.Data == nil {
						t.Errorf("ExtractPGPSignatures() signature %d has nil Data", i)
					}
				}
			}
		})
	}
}

func TestExtractPGPKeys(t *testing.T) {
	tests := []struct {
		name          string
		data          string
		wantNum       int
		wantFilenames []string
		wantErr       bool
	}{
		{
			name: "email with pgp key",
			data: `From: sender@example.com
To: recipient@example.com
Subject: PGP Key
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here is my public key.

--boundary123
Content-Type: application/pgp-keys
Content-Disposition: attachment; filename="pubkey.asc"

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Example
Comment: GPGTools - https://gpgtools.org

mQINBGVtZIUBEADJ9Xdx5LJSgLMY7rFUGQR3YjRFR8PNY9F5DQ+pq7bn2FvKThBu
ExampleKeyContentExampleKeyContentExampleKeyContent
-----END PGP PUBLIC KEY BLOCK-----
--boundary123--`,
			wantNum:       1,
			wantFilenames: []string{"pubkey.asc"},
			wantErr:       false,
		},
		{
			name: "email with base64 encoded pgp key",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Base64 PGP Key
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here is my base64 encoded public key.

--boundary123
Content-Type: application/pgp-keys
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="pubkey.asc"

LS0tLS1CRUdJTiBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0tClZlcnNpb246IEV4YW1w
bGUKQ29tbWVudDogR1BHVG9vbHMgLSBodHRwczovL2dwZ3Rvb2xzLm9yZwoKbVFJTkJH
VnRaSVVCRUFESjlYZHg1TEpTZ0xNWTdyRlVHUVIzWWpSRlI4UE5ZOUY1RFErcXE3Ym4y
RnZLVGhCdQpFeGFtcGxlS2V5Q29udGVudEV4YW1wbGVLZXlDb250ZW50RXhhbXBsZUtl
eUNvbnRlbnQKLS0tLS1FTkQgUEdQIFBVQkxJQyBLRVkgQkxPQ0stLS0tLQ==
--boundary123--`,
			wantNum:       1,
			wantFilenames: []string{"pubkey.asc"},
			wantErr:       false,
		},
		{
			name: "email with multiple pgp keys",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Multiple PGP Keys
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here are my public keys.

--boundary123
Content-Type: application/pgp-keys
Content-Disposition: attachment; filename="key1.asc"

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Example1
Comment: GPGTools - https://gpgtools.org

mQINBGVtZIUBEADJ9Xdx5LJSgLMY7rFUGQR3YjRFR8PNY9F5DQ+pq7bn2FvKThBu
Key1ContentKey1ContentKey1Content
-----END PGP PUBLIC KEY BLOCK-----
--boundary123
Content-Type: application/pgp-keys
Content-Disposition: attachment; filename="key2.asc"

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Example2
Comment: GPGTools - https://gpgtools.org

mQINBGVtZIUBEADJ9Xdx5LJSgLMY7rFUGQR3YjRFR8PNY9F5DQ+pq7bn2FvKThBu
Key2ContentKey2ContentKey2Content
-----END PGP PUBLIC KEY BLOCK-----
--boundary123--`,
			wantNum:       2,
			wantFilenames: []string{"key1.asc", "key2.asc"},
			wantErr:       false,
		},
		{
			name: "pgp key without filename",
			data: `From: sender@example.com
To: recipient@example.com
Subject: PGP Key No Filename
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here is my public key.

--boundary123
Content-Type: application/pgp-keys

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Example
Comment: GPGTools - https://gpgtools.org

mQINBGVtZIUBEADJ9Xdx5LJSgLMY7rFUGQR3YjRFR8PNY9F5DQ+pq7bn2FvKThBu
ExampleKeyContentExampleKeyContentExampleKeyContent
-----END PGP PUBLIC KEY BLOCK-----
--boundary123--`,
			wantNum:       1,
			wantFilenames: []string{"publickey.asc"},
			wantErr:       false,
		},
		{
			name: "email without pgp keys",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Regular Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a regular email without keys.
--boundary123--`,
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name: "non-multipart email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Plain Email
Content-Type: text/plain

This is just a plain text email.`,
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name:          "invalid email data",
			data:          "Invalid email data",
			wantNum:       0,
			wantFilenames: []string{},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPGPKeys([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPGPKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != tt.wantNum {
					t.Errorf("ExtractPGPKeys() returned %d keys, want %d", len(got), tt.wantNum)
					return
				}

				for i, attachment := range got {
					if i < len(tt.wantFilenames) && attachment.Filename != tt.wantFilenames[i] {
						t.Errorf("ExtractPGPKeys() key %d filename = %q, want %q",
							i, attachment.Filename, tt.wantFilenames[i])
					}
					if !strings.HasPrefix(attachment.ContentType, "application/pgp-keys") {
						t.Errorf("ExtractPGPKeys() key %d ContentType = %q, want prefix %q",
							i, attachment.ContentType, "application/pgp-keys")
					}
					if attachment.Data == nil {
						t.Errorf("ExtractPGPKeys() key %d has nil Data", i)
					}
				}
			}
		})
	}
}

func TestExtractTextBody(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{
			name: "simple plain text email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Plain Text Email
Content-Type: text/plain

This is a plain text email body.`,
			want:    "This is a plain text email body.",
			wantErr: false,
		},
		{
			name: "multipart email with plain text part",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Multipart Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is the plain text part.
--boundary123
Content-Type: text/html

<html><body>This is the HTML part.</body></html>
--boundary123--`,
			want:    "This is the plain text part.",
			wantErr: false,
		},
		{
			name: "multipart email with encoded plain text part",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Encoded Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain
Content-Transfer-Encoding: base64

VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCB0ZXh0Lg==
--boundary123--`,
			want:    "This is base64 encoded text.",
			wantErr: false,
		},
		{
			name: "multipart email with quoted-printable encoding",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Quoted Printable Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain
Content-Transfer-Encoding: quoted-printable

This is quoted=3Dprintable text with special chars =C3=A4=C3=B6=C3=BC.
--boundary123--`,
			want:    "This is quoted=printable text with special chars äöü.",
			wantErr: false,
		},
		{
			name: "html only email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: HTML Only Email
Content-Type: text/html

<html><body>This is HTML content only.</body></html>`,
			want:    "",
			wantErr: true,
		},
		{
			name: "nested multipart email with plain text",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Nested Multipart Email
Content-Type: multipart/mixed; boundary="outer"

--outer
Content-Type: multipart/alternative; boundary="inner"

--inner
Content-Type: text/plain

This is the nested plain text.
--inner
Content-Type: text/html

<html><body>This is nested HTML.</body></html>
--inner--
--outer--`,
			want:    "This is the nested plain text.",
			wantErr: false,
		},
		{
			name:    "invalid email data",
			data:    "Invalid email data",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTextBody([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTextBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractTextBody() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractHTMLBody(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{
			name: "simple html email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: HTML Email
Content-Type: text/html

<html><body>This is an HTML email body.</body></html>`,
			want:    "<html><body>This is an HTML email body.</body></html>",
			wantErr: false,
		},
		{
			name: "multipart email with html part",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Multipart Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is the plain text part.
--boundary123
Content-Type: text/html

<html><body>This is the HTML part.</body></html>
--boundary123--`,
			want:    "<html><body>This is the HTML part.</body></html>",
			wantErr: false,
		},
		{
			name: "multipart email with encoded html part",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Encoded Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/html
Content-Transfer-Encoding: base64

PGh0bWw+PGJvZHk+VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCBIVE1MLjwvYm9keT48L2h0bWw+
--boundary123--`,
			want:    "<html><body>This is base64 encoded HTML.</body></html>",
			wantErr: false,
		},
		{
			name: "multipart email with quoted-printable html encoding",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Quoted Printable Email
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/html
Content-Transfer-Encoding: quoted-printable

<html><body>This is quoted=3Dprintable HTML with special chars =C3=A4=C3=B6=C3=BC.</body></html>
--boundary123--`,
			want:    "<html><body>This is quoted=printable HTML with special chars äöü.</body></html>",
			wantErr: false,
		},
		{
			name: "plain text only email",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Plain Text Only Email
Content-Type: text/plain

This is plain text content only.`,
			want:    "",
			wantErr: true,
		},
		{
			name: "nested multipart email with html",
			data: `From: sender@example.com
To: recipient@example.com
Subject: Nested Multipart Email
Content-Type: multipart/mixed; boundary="outer"

--outer
Content-Type: multipart/alternative; boundary="inner"

--inner
Content-Type: text/plain

This is the nested plain text.
--inner
Content-Type: text/html

<html><body>This is nested HTML content.</body></html>
--inner--
--outer--`,
			want:    "<html><body>This is nested HTML content.</body></html>",
			wantErr: false,
		},
		{
			name:    "invalid email data",
			data:    "Invalid email data",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractHTMLBody([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractHTMLBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractHTMLBody() = %q, want %q", got, tt.want)
			}
		})
	}
}

func compareMessages(a, b Msg) bool {
	if a.From != b.From || a.FromName != b.FromName || a.Subject != b.Subject || a.Body != b.Body || a.Type != b.Type {
		return false
	}
	if len(a.To) != len(b.To) {
		return false
	}
	for i := range a.To {
		if a.To[i] != b.To[i] {
			return false
		}
	}
	return true
}

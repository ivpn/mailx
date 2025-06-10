package model

import (
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
		emailData     string
		wantFilenames []string
		wantErr       bool
	}{
		{
			name: "email with pgp signature",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Signed Message
MIME-Version: 1.0
Content-Type: multipart/signed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a signed message.

--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="signature.asc"
Content-Transfer-Encoding: 7bit

-----BEGIN PGP SIGNATURE-----
Version: GnuPG v2.0.22

iQEcBAEBAgAGBQJTqB5EAAoJEJ96b4jpV5TbIUgH/3GUaTu4RaXw2Vf5MH1JQJ0u
qpB8mUEsAjGxL2M=
=qTm0
-----END PGP SIGNATURE-----
--boundary123--
`,
			wantFilenames: []string{"signature.asc"},
			wantErr:       false,
		},
		{
			name: "email with base64 encoded pgp signature",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Signed Message
MIME-Version: 1.0
Content-Type: multipart/signed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a signed message.

--boundary123
Content-Type: application/pgp-signature
Content-Transfer-Encoding: base64

LS0tLS1CRUdJTiBQR1AgU0lHTkFUVVJFLS0tLS0KVmVyc2lvbjogR251UEcgdjIuMC4y
MgoKaVFFY0JBRUJBbUFHQlFKVHFCNUVBQW9KRUo5NmI0anBWNVRiSVVnSC8zR1VhVHU0
UmFYdzJWZjVNSDFKUUowdQpxcEI4bVVFc0FqR3hMMk09Cj1xVG0wCi0tLS0tRU5EIFBH
UCBTSUdOQVRVUkUtLS0tLQo=
--boundary123--
`,
			wantFilenames: []string{"signature.asc"},
			wantErr:       false,
		},
		{
			name: "email with custom filename",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Signed Message
MIME-Version: 1.0
Content-Type: multipart/signed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a signed message.

--boundary123
Content-Type: application/pgp-signature
Content-Disposition: attachment; filename="custom.sig"

-----BEGIN PGP SIGNATURE-----
Version: GnuPG v2.0.22

iQEcBAEBAgAGBQJTqB5EAAoJEJ96b4jpV5TbIUgH/3GUaTu4RaXw2Vf5MH1JQJ0u
qpB8mUEsAjGxL2M=
=qTm0
-----END PGP SIGNATURE-----
--boundary123--
`,
			wantFilenames: []string{"custom.sig"},
			wantErr:       false,
		},
		{
			name: "email without pgp signature",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Regular Message
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a regular message.

--boundary123
Content-Type: application/pdf
Content-Disposition: attachment; filename="document.pdf"

PDF content here
--boundary123--
`,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name: "non-multipart email",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Plain Text
Content-Type: text/plain

This is just a plain text message.
`,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name:          "invalid email",
			emailData:     "Invalid email data",
			wantFilenames: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPGPSignatures([]byte(tt.emailData))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPGPSignatures() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if the number of attachments matches expected
				if len(got) != len(tt.wantFilenames) {
					t.Errorf("ExtractPGPSignatures() got %d attachments, want %d", len(got), len(tt.wantFilenames))
					return
				}

				// Check if filenames match
				for i, attachment := range got {
					if attachment.Filename != tt.wantFilenames[i] {
						t.Errorf("ExtractPGPSignatures() attachment[%d].Filename = %v, want %v",
							i, attachment.Filename, tt.wantFilenames[i])
					}

					// Verify content type contains application/pgp-signature
					if !strings.Contains(attachment.ContentType, "application/pgp-signature") {
						t.Errorf("ExtractPGPSignatures() attachment[%d].ContentType = %v, want to contain 'application/pgp-signature'",
							i, attachment.ContentType)
					}

					// Check that Data field is not nil
					if attachment.Data == nil {
						t.Errorf("ExtractPGPSignatures() attachment[%d].Data is nil", i)
					}
				}
			}
		})
	}
}

func TestExtractPGPKeys(t *testing.T) {
	tests := []struct {
		name          string
		emailData     string
		wantFilenames []string
		wantErr       bool
	}{
		{
			name: "email with pgp keys",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: My Public Key
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here's my public key.

--boundary123
Content-Type: application/pgp-keys
Content-Disposition: attachment; filename="pubkey.asc"
Content-Transfer-Encoding: 7bit

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2.0.22

mQINBFOoHkQBEACqM/x6GxQg6Is3zj2PQcw7iMzRqF+uHJYASbKmLdNdBIhkdwfa
KHDa7J65bDJX0jU=
=ZCUI
-----END PGP PUBLIC KEY BLOCK-----
--boundary123--
`,
			wantFilenames: []string{"pubkey.asc"},
			wantErr:       false,
		},
		{
			name: "email with base64 encoded pgp keys",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: My Public Key
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here's my encoded public key.

--boundary123
Content-Type: application/pgp-keys
Content-Transfer-Encoding: base64

LS0tLS1CRUdJTiBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0tClZlcnNpb246IEdudVBH
IHYyLjAuMjIKCm1RSU5CRk9vSGtRQkVBQ3FNL3g2R3hRZzZJczN6ajJQUWN3N2lNelJx
Rit1SEpZQVNiS21MZE5kQkl5VXh3ZmEKS2hEYTdKNjViREpYMGpVPQo9WkNVSQotLS0t
LUVORCBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0tCg==
--boundary123--
`,
			wantFilenames: []string{"publickey.asc"},
			wantErr:       false,
		},
		{
			name: "email with custom filename",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: My Public Key
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

Here's my public key with custom filename.

--boundary123
Content-Type: application/pgp-keys
Content-Disposition: attachment; filename="mykey.gpg"

-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2.0.22

mQINBFOoHkQBEACqM/x6GxQg6Is3zj2PQcw7iMzRqF+uHJYASbKmLdNdBIhkdwfa
KHDa7J65bDJX0jU=
=ZCUI
-----END PGP PUBLIC KEY BLOCK-----
--boundary123--
`,
			wantFilenames: []string{"mykey.gpg"},
			wantErr:       false,
		},
		{
			name: "email without pgp keys",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Regular Message
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary123"

--boundary123
Content-Type: text/plain

This is a regular message.

--boundary123
Content-Type: application/pdf
Content-Disposition: attachment; filename="document.pdf"

PDF content here
--boundary123--
`,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name: "non-multipart email",
			emailData: `From: sender@example.com
To: recipient@example.com
Subject: Plain Text
Content-Type: text/plain

This is just a plain text message.
`,
			wantFilenames: []string{},
			wantErr:       false,
		},
		{
			name:          "invalid email",
			emailData:     "Invalid email data",
			wantFilenames: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPGPKeys([]byte(tt.emailData))
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPGPKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if the number of attachments matches expected
				if len(got) != len(tt.wantFilenames) {
					t.Errorf("ExtractPGPKeys() got %d attachments, want %d", len(got), len(tt.wantFilenames))
					return
				}

				// Check if filenames match
				for i, attachment := range got {
					if attachment.Filename != tt.wantFilenames[i] {
						t.Errorf("ExtractPGPKeys() attachment[%d].Filename = %v, want %v",
							i, attachment.Filename, tt.wantFilenames[i])
					}

					// Verify content type contains application/pgp-keys
					if !strings.Contains(attachment.ContentType, "application/pgp-keys") {
						t.Errorf("ExtractPGPKeys() attachment[%d].ContentType = %v, want to contain 'application/pgp-keys'",
							i, attachment.ContentType)
					}

					// Check that Data field is not nil
					if attachment.Data == nil {
						t.Errorf("ExtractPGPKeys() attachment[%d].Data is nil", i)
					}
				}
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

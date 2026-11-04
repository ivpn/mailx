package mailer

import (
	"bytes"
	"mime"
	"net/mail"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"ivpn.net/email/api/internal/model"
)

const senderPGPMIMEInner = "Content-Type: text/plain; charset=utf-8\r\n\r\nMeet at noon.\r\n"

// senderPGPMIMEEmail returns an email the sender encrypted to key with PGP/MIME. It has LF
// line endings, as mail piped in by Postfix does.
func senderPGPMIMEEmail(t *testing.T, key *crypto.Key) string {
	t.Helper()

	enc, err := crypto.PGP().Encryption().Recipient(key).New()
	if err != nil {
		t.Fatalf("encryption handle: %v", err)
	}
	msg, err := enc.Encrypt([]byte(senderPGPMIMEInner))
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	armored, err := msg.ArmorBytes()
	if err != nil {
		t.Fatalf("ArmorBytes() error = %v", err)
	}

	return "From: Friend <friend@example.com>\n" +
		"To: myalias@mailx.net\n" +
		"Subject: Secret\n" +
		"Message-ID: <secret@example.com>\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: multipart/encrypted; protocol=\"application/pgp-encrypted\"; boundary=\"enc-boundary\"\n" +
		"\n" +
		"This is an OpenPGP/MIME encrypted message (RFC 4880 and 3156)\n" +
		"--enc-boundary\n" +
		"Content-Type: application/pgp-encrypted\n" +
		"Content-Description: PGP/MIME version identification\n" +
		"\n" +
		"Version: 1\n" +
		"\n" +
		"--enc-boundary\n" +
		"Content-Type: application/octet-stream; name=\"encrypted.asc\"\n" +
		"Content-Description: OpenPGP encrypted message\n" +
		"Content-Disposition: inline; filename=\"encrypted.asc\"\n" +
		"\n" +
		strings.TrimSuffix(string(armored), "\n") + "\n" +
		"\n" +
		"--enc-boundary--\n"
}

func TestForward_SenderPGPMIMEKeepsStructure(t *testing.T) {
	key := testKey(t)
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}
	raw := senderPGPMIMEEmail(t, key)

	tests := []struct {
		name    string
		rcp     model.Recipient
		outcome model.RuleOutcome
	}{
		{"PGP/MIME recipient", model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey}, model.RuleOutcome{}},
		{"PGP/Inline recipient", model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey, PGPInline: true}, model.RuleOutcome{}},
		{"unencrypted recipient", plainRecipient, model.RuleOutcome{}},
		{"rules applied", model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey},
			model.RuleOutcome{MatchedNames: []string{"Friends"}, SubjectValues: []string{"[Friends] {subject}"}, RemoveAttachments: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if _, err := buildRulesForwardMessage(t, raw, "Friend", tt.rcp, model.Settings{BlockTrackers: true}, tt.outcome).WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}
			if bytes.Contains(bytes.ReplaceAll(buf.Bytes(), []byte("\r\n"), nil), []byte("\n")) {
				t.Error("forward has bare LF line endings")
			}

			parsed, err := mail.ReadMessage(&buf)
			if err != nil {
				t.Fatalf("mail.ReadMessage() error = %v", err)
			}

			if cts := parsed.Header["Content-Type"]; len(cts) != 1 {
				t.Fatalf("Content-Type headers = %q, want one", cts)
			}
			mediaType, params, err := mime.ParseMediaType(parsed.Header.Get("Content-Type"))
			if err != nil || mediaType != "multipart/encrypted" || params["protocol"] != "application/pgp-encrypted" || params["boundary"] != "enc-boundary" {
				t.Fatalf("Content-Type = %q, want the sender's multipart/encrypted", parsed.Header.Get("Content-Type"))
			}
			if got := parsed.Header.Get("Mime-Version"); got != "1.0" {
				t.Errorf("Mime-Version = %q", got)
			}
			if got := parsed.Header.Get("X-Mailx-Original-To"); got != "real@domain.com" {
				t.Errorf("X-Mailx-Original-To = %q", got)
			}
			if got := parsed.Header.Get("Message-Id"); !strings.Contains(got, "secret@example.com") {
				t.Errorf("Message-ID = %q", got)
			}
			wantSubject := "Secret"
			if len(tt.outcome.MatchedNames) > 0 {
				wantSubject = "[Friends] Secret"
				if got := parsed.Header.Get("X-Mailx-Rules-Applied"); got != "Friends" {
					t.Errorf("X-Mailx-Rules-Applied = %q", got)
				}
			}
			if got := decodedHeader(t, parsed.Header, "Subject"); got != wantSubject {
				t.Errorf("Subject = %q, want %q", got, wantSubject)
			}

			parts := leafParts(t, parsed.Header.Get("Content-Type"), "", parsed.Body)
			if len(parts) != 2 {
				t.Fatalf("got %d parts, want the sender's 2", len(parts))
			}
			if ct := parts[0].header.Get("Content-Type"); ct != "application/pgp-encrypted" {
				t.Errorf("part 1 Content-Type = %q", ct)
			}
			if ct := parts[1].header.Get("Content-Type"); ct != `application/octet-stream; name="encrypted.asc"` {
				t.Errorf("part 2 Content-Type = %q", ct)
			}

			// The recipient decrypts the sender's message, not one MailX wrapped it in
			if got := decryptWith(t, key, parts[1].data, crypto.Armor); got != senderPGPMIMEInner {
				t.Errorf("decrypts to %q, want the sender's message", got)
			}
		})
	}
}

func TestSenderPGPMIME(t *testing.T) {
	const body = "--b\nContent-Type: application/pgp-encrypted\n\nVersion: 1\n--b--\n"

	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"PGP/MIME", `multipart/encrypted; protocol="application/pgp-encrypted"; boundary="b"`, true},
		{"case-insensitive", `Multipart/Encrypted; Protocol="Application/PGP-Encrypted"; boundary=b`, true},
		{"no protocol", `multipart/encrypted; boundary="b"`, false},
		{"other protocol", `multipart/encrypted; protocol="application/x-other"; boundary="b"`, false},
		{"no boundary", `multipart/encrypted; protocol="application/pgp-encrypted"`, false},
		{"multipart/mixed", `multipart/mixed; boundary="b"`, false},
		{"S/MIME", `application/pkcs7-mime; smime-type=enveloped-data`, false},
		{"invalid", `multipart/encrypted; protocol=`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := "Subject: x\nContent-Type: " + tt.contentType + "\n\n" + body
			ct, got, ok := senderPGPMIME([]byte(raw))
			if ok != tt.want {
				t.Fatalf("ok = %v, want %v", ok, tt.want)
			}
			if !ok {
				return
			}
			if ct != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.contentType)
			}
			if want := strings.ReplaceAll(body, "\n", "\r\n"); got != want {
				t.Errorf("body = %q, want %q", got, want)
			}
		})
	}
}

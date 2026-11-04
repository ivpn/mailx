package mailer

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"ivpn.net/email/api/internal/model"
)

const testListUnsubscribeRawEmail = "From: Sender <sender@example.com>\r\n" +
	"To: myalias@mailx.net\r\n" +
	"Subject: Newsletter\r\n" +
	"List-Unsubscribe: <mailto:unsub@example.com>\r\n" +
	"List-Unsubscribe-Post: List-Unsubscribe=One-Click\r\n" +
	"\r\n" +
	"Test body content\r\n"

var testListUnsubscribe = model.ListUnsubscribeHeaders{
	Value: "<https://api.mailx.net/v1/deactivate-one-click/alias-1?expires=1&signature=" + strings.Repeat("ab", 32) + ">",
	Post:  "List-Unsubscribe=One-Click",
}

func buildTestForward(t *testing.T, rcp model.Recipient, listUnsubscribe model.ListUnsubscribeHeaders) *mail.Message {
	t.Helper()

	m := Mailer{}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	templateData := map[string]any{"alias": alias.Name, "from": "sender@example.com"}

	msg, _, err := m.buildForwardMessage(
		"myalias+sender=example.com@mailx.net", "Sender", alias.Name, rcp,
		[]byte(testListUnsubscribeRawEmail), "header.tmpl", templateData, model.Settings{}, alias, listUnsubscribe, model.RuleOutcome{},
	)
	if err != nil {
		t.Fatalf("buildForwardMessage() error = %v", err)
	}

	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	parsed, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("mail.ReadMessage() error = %v", err)
	}

	return parsed
}

func assertListUnsubscribe(t *testing.T, header mail.Header, want model.ListUnsubscribeHeaders) {
	t.Helper()

	if got := header["List-Unsubscribe"]; len(got) != 1 || got[0] != want.Value {
		t.Errorf("List-Unsubscribe = %q, want exactly [%q]", got, want.Value)
	}
	if got := header["List-Unsubscribe-Post"]; len(got) != 1 || got[0] != want.Post {
		t.Errorf("List-Unsubscribe-Post = %q, want exactly [%q]", got, want.Post)
	}
}

func TestBuildForwardMessage_SetsListUnsubscribe(t *testing.T) {
	parsed := buildTestForward(t, model.Recipient{Email: "real@domain.com"}, testListUnsubscribe)

	// The long signed URL gets folded when written, and must unfold to the same value.
	assertListUnsubscribe(t, parsed.Header, testListUnsubscribe)
}

func TestBuildForwardMessage_NoListUnsubscribeDropsOriginal(t *testing.T) {
	parsed := buildTestForward(t, model.Recipient{Email: "real@domain.com"}, model.ListUnsubscribeHeaders{})

	if got := parsed.Header.Get("List-Unsubscribe"); got != "" {
		t.Errorf("List-Unsubscribe = %q, want the original header not copied", got)
	}
	if got := parsed.Header.Get("List-Unsubscribe-Post"); got != "" {
		t.Errorf("List-Unsubscribe-Post = %q, want the original header not copied", got)
	}
}

func TestBuildForwardMessage_PGPMIMEKeepsListUnsubscribe(t *testing.T) {
	key, err := crypto.PGP().KeyGeneration().AddUserId("Test", "real@domain.com").New().GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}

	rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey}
	parsed := buildTestForward(t, rcp, testListUnsubscribe)

	if ct := parsed.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/encrypted") {
		t.Fatalf("Content-Type = %q, want a PGP/MIME message", ct)
	}
	assertListUnsubscribe(t, parsed.Header, testListUnsubscribe)
}

package mailer

import (
	"testing"

	"ivpn.net/email/api/internal/model"
)

const testForwardRawEmail = "From: Sender <sender@example.com>\r\n" +
	"To: original-to@example.com\r\n" +
	"Subject: Hello\r\n" +
	"\r\n" +
	"Test body content\r\n"

func TestBuildForwardMessage_ToHeaderUsesOriginalAliasAddress(t *testing.T) {
	m := Mailer{}
	rcp := model.Recipient{Email: "real.recipient@domain.com"}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	settings := model.Settings{}
	to := "myalias@mailx.net"
	templateData := map[string]any{"alias": to, "from": "sender@example.com"}

	msg, _, err := m.buildForwardMessage(
		"sender@example.com", "Sender", to, rcp,
		[]byte(testForwardRawEmail), "header.tmpl", templateData, settings, alias,
	)
	if err != nil {
		t.Fatalf("buildForwardMessage() error = %v", err)
	}

	got := msg.GetHeader("To")
	if len(got) != 1 || got[0] != to {
		t.Errorf("To header = %v, want [%s]", got, to)
	}
	for _, v := range got {
		if v == rcp.Email {
			t.Errorf("To header must not contain the real recipient mailbox %q", rcp.Email)
		}
	}

	// The real mailbox is still preserved in the diagnostic header.
	originalTo := msg.GetHeader("X-Mailx-Original-To")
	if len(originalTo) != 1 || originalTo[0] != rcp.Email {
		t.Errorf("X-Mailx-Original-To = %v, want [%s]", originalTo, rcp.Email)
	}
}

func TestBuildForwardMessage_ToHeaderUsesLiteralTaggedAddress(t *testing.T) {
	m := Mailer{}
	rcp := model.Recipient{Email: "real.recipient@domain.com"}
	// alias.Name is the canonical, tag-stripped alias; it can even be a literal
	// wildcard pattern (e.g. "*+news@customdomain.com") for wildcard/catch-all
	// aliases, so the To: header must use the concrete address the sender used
	// instead ("to"), not alias.Name.
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	settings := model.Settings{}
	to := "myalias+shop@mailx.net"
	templateData := map[string]any{"alias": to, "from": "sender@example.com"}

	msg, _, err := m.buildForwardMessage(
		"sender@example.com", "Sender", to, rcp,
		[]byte(testForwardRawEmail), "header.tmpl", templateData, settings, alias,
	)
	if err != nil {
		t.Fatalf("buildForwardMessage() error = %v", err)
	}

	got := msg.GetHeader("To")
	if len(got) != 1 || got[0] != to {
		t.Errorf("To header = %v, want [%s]", got, to)
	}
}

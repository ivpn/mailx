package mailer

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"ivpn.net/email/api/internal/model"
)

// testThreadRawEmail is testRulesRawEmail sent as a reply in an existing conversation.
var testThreadRawEmail = strings.Replace(testRulesRawEmail,
	"Message-ID: <rules@shop.example>\r\n",
	"Message-ID: <rules@shop.example>\r\n"+
		"In-Reply-To: <parent@shop.example>\r\n"+
		"References: <root@shop.example>\r\n <parent@shop.example>\r\n",
	1)

func assertThreadHeaders(t *testing.T, h mail.Header, wantMessageID bool) {
	t.Helper()

	want := map[string]string{
		"In-Reply-To": "<parent@shop.example>",
		"References":  "<root@shop.example> <parent@shop.example>",
	}
	if wantMessageID {
		want["Message-Id"] = "<rules@shop.example>"
	}
	for k, v := range want {
		if got := h[k]; len(got) != 1 || got[0] != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
}

func assertContentID(t *testing.T, parts []mimePart) {
	t.Helper()

	for _, p := range parts {
		if p.contentType == "image/png" {
			if p.contentID != "<logo@shop.example>" {
				t.Errorf("inline image Content-ID = %q, want %q", p.contentID, "<logo@shop.example>")
			}
			return
		}
	}
	t.Errorf("parts = %+v, want the inline image", parts)
}

func TestForward_KeepsMessageIDBrackets(t *testing.T) {
	key := testKey(t)
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}

	t.Run("unencrypted", func(t *testing.T) {
		fwd := buildRulesForward(t, testThreadRawEmail, "Shop", plainRecipient, model.Settings{}, model.RuleOutcome{})
		assertThreadHeaders(t, fwd.header, true)
		assertContentID(t, fwd.parts)
	})

	t.Run("PGP/Inline", func(t *testing.T) {
		rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey, PGPInline: true}
		fwd := buildRulesForward(t, testThreadRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{})
		assertThreadHeaders(t, fwd.header, true)
	})

	t.Run("PGP/MIME", func(t *testing.T) {
		rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey}
		var buf bytes.Buffer
		if _, err := buildRulesForwardMessage(t, testThreadRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{}).WriteTo(&buf); err != nil {
			t.Fatalf("WriteTo() error = %v", err)
		}
		outer, err := mail.ReadMessage(&buf)
		if err != nil {
			t.Fatalf("mail.ReadMessage() error = %v", err)
		}
		// The outer message is what the recipient's mail server threads on
		assertThreadHeaders(t, outer.Header, true)

		parts := leafParts(t, outer.Header.Get("Content-Type"), "", outer.Body)
		if len(parts) != 2 {
			t.Fatalf("got %d parts, want 2", len(parts))
		}
		inner, err := mail.ReadMessage(strings.NewReader(decryptWith(t, key, parts[1].data, crypto.Armor)))
		if err != nil {
			t.Fatalf("read decrypted message: %v", err)
		}
		assertThreadHeaders(t, inner.Header, true)
		assertContentID(t, collectParts(t, inner.Header.Get("Content-Type"), inner.Body))
	})

	t.Run("sender PGP/MIME", func(t *testing.T) {
		raw := strings.Replace(senderPGPMIMEEmail(t, key), "Message-ID: <secret@example.com>\n", "Message-ID: <secret@example.com>\nIn-Reply-To: <parent@example.com>\n", 1)
		fwd := buildRulesForward(t, raw, "Friend", plainRecipient, model.Settings{}, model.RuleOutcome{})
		if got := fwd.header["Message-Id"]; len(got) != 1 || got[0] != "<secret@example.com>" {
			t.Errorf("Message-ID = %q, want %q", got, "<secret@example.com>")
		}
		if got := fwd.header["In-Reply-To"]; len(got) != 1 || got[0] != "<parent@example.com>" {
			t.Errorf("In-Reply-To = %q, want %q", got, "<parent@example.com>")
		}
	})
}

func TestReply_KeepsMessageIDBrackets(t *testing.T) {
	m := Mailer{}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}

	msg, _, err := m.buildReplyMessage(alias.Name, "Me", model.Recipient{Email: "news@shop.example"}, []byte(testThreadRawEmail), alias, model.RuleOutcome{})
	if err != nil {
		t.Fatalf("buildReplyMessage() error = %v", err)
	}
	reply := writeAndParse(t, msg)

	// A reply gets its own Message-ID from the mail server, not the one it answers
	assertThreadHeaders(t, reply.header, false)
	if got := reply.header.Get("Message-Id"); got == "<rules@shop.example>" {
		t.Errorf("Message-ID = %q, want a new one", got)
	}
	assertContentID(t, reply.parts)
}

func TestBracketID(t *testing.T) {
	tests := map[string]string{
		"":              "",
		"a@example.com": "<a@example.com>",
	}
	for in, want := range tests {
		if got := bracketID(in); got != want {
			t.Errorf("bracketID(%q) = %q, want %q", in, got, want)
		}
	}
}

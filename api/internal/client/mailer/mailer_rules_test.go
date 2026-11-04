package mailer

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"slices"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils/gomail.v2"
)

// testRulesRawEmail has a PDF attachment and an inline cid: image.
const testRulesRawEmail = "From: Shop <news@shop.example>\r\n" +
	"To: myalias@mailx.net\r\n" +
	"Subject: Weekly deals\r\n" +
	"Message-ID: <rules@shop.example>\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Content-Type: multipart/mixed; boundary=\"mixed\"\r\n" +
	"\r\n" +
	"--mixed\r\n" +
	"Content-Type: multipart/related; boundary=\"related\"\r\n" +
	"\r\n" +
	"--related\r\n" +
	"Content-Type: text/html; charset=utf-8\r\n" +
	"\r\n" +
	"<p>Deals</p><img src=\"cid:logo@shop.example\">\r\n" +
	"--related\r\n" +
	"Content-Type: image/png\r\n" +
	"Content-Transfer-Encoding: base64\r\n" +
	"Content-ID: <logo@shop.example>\r\n" +
	"Content-Disposition: inline; filename=\"logo.png\"\r\n" +
	"\r\n" +
	"iVBORw0KGgo=\r\n" +
	"--related--\r\n" +
	"--mixed\r\n" +
	"Content-Type: application/pdf; name=\"deals.pdf\"\r\n" +
	"Content-Transfer-Encoding: base64\r\n" +
	"Content-Disposition: attachment; filename=\"deals.pdf\"\r\n" +
	"\r\n" +
	"JVBERi0xLjQ=\r\n" +
	"--mixed--\r\n"

type builtMessage struct {
	header mail.Header
	parts  []mimePart
	html   string
}

type mimePart struct {
	contentType string
	disposition string
	filename    string
	contentID   string
}

func writeAndParse(t *testing.T, msg *gomail.Message) builtMessage {
	t.Helper()

	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	raw := buf.Bytes()

	parsed, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("mail.ReadMessage() error = %v", err)
	}
	built := builtMessage{header: parsed.Header}
	built.parts = collectParts(t, parsed.Header.Get("Content-Type"), parsed.Body)

	parsed, _ = mail.ReadMessage(bytes.NewReader(raw))
	built.html, _ = findHTMLPart(t, parsed.Header.Get("Content-Type"), parsed.Header.Get("Content-Transfer-Encoding"), parsed.Body)

	return built
}

func collectParts(t *testing.T, contentType string, body io.Reader) []mimePart {
	t.Helper()

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return nil
	}

	var parts []mimePart
	mr := multipart.NewReader(body, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err != nil {
			return parts
		}
		ct := part.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart/") {
			parts = append(parts, collectParts(t, ct, part)...)
			continue
		}
		disposition, dparams, _ := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		mediaType, _, _ := mime.ParseMediaType(ct)
		parts = append(parts, mimePart{mediaType, disposition, dparams["filename"], part.Header.Get("Content-Id")})
	}
}

func (b builtMessage) hasPart(match func(mimePart) bool) bool {
	return slices.ContainsFunc(b.parts, match)
}

func decodedHeader(t *testing.T, h mail.Header, key string) string {
	t.Helper()

	decoded, err := new(mime.WordDecoder).DecodeHeader(h.Get(key))
	if err != nil {
		t.Fatalf("DecodeHeader(%s) error = %v", key, err)
	}
	return decoded
}

func buildRulesForward(t *testing.T, raw string, name string, rcp model.Recipient, settings model.Settings, outcome model.RuleOutcome) builtMessage {
	t.Helper()

	return writeAndParse(t, buildRulesForwardMessage(t, raw, name, rcp, settings, outcome))
}

func buildRulesForwardMessage(t *testing.T, raw string, name string, rcp model.Recipient, settings model.Settings, outcome model.RuleOutcome) *gomail.Message {
	t.Helper()

	m := Mailer{}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	templateData := map[string]any{"alias": alias.Name, "from": "news@shop.example"}

	msg, _, err := m.buildForwardMessage(
		"myalias+news=shop.example@mailx.net", name, alias.Name, rcp,
		[]byte(raw), "header.tmpl", templateData, settings, alias, model.ListUnsubscribeHeaders{}, outcome,
	)
	if err != nil {
		t.Fatalf("buildForwardMessage() error = %v", err)
	}

	return msg
}

func testKey(t *testing.T) *crypto.Key {
	t.Helper()

	key, err := crypto.PGP().KeyGeneration().AddUserId("Test", "real@domain.com").New().GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	return key
}

func testPublicKey(t *testing.T) string {
	t.Helper()

	publicKey, err := testKey(t).GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}
	return publicKey
}

var plainRecipient = model.Recipient{Email: "real@domain.com"}

func TestForwardRules_NoOutcomeLeavesForwardUnchanged(t *testing.T) {
	fwd := buildRulesForward(t, testRulesRawEmail, "Shop", plainRecipient, model.Settings{}, model.RuleOutcome{})

	if got := decodedHeader(t, fwd.header, "Subject"); got != "Weekly deals" {
		t.Errorf("Subject = %q", got)
	}
	if _, ok := fwd.header["X-Mailx-Rules-Applied"]; ok {
		t.Error("X-Mailx-Rules-Applied set without matching rules")
	}
	if !fwd.hasPart(func(p mimePart) bool { return p.filename == "deals.pdf" }) {
		t.Errorf("parts = %+v, want the attachment", fwd.parts)
	}
}

func TestForwardRules_SubjectAndDisplayName(t *testing.T) {
	outcome := model.RuleOutcome{
		MatchedNames:  []string{"Deals", "Tag"},
		SubjectValues: []string{"[Deals] {subject}", "Ünïcode · {subject}"},
	}
	// The service passes the rule's sender name in place of the email's own
	fwd := buildRulesForward(t, testRulesRawEmail, "Renamed Shop", plainRecipient, model.Settings{}, outcome)

	if got := decodedHeader(t, fwd.header, "Subject"); got != "Ünïcode · [Deals] Weekly deals" {
		t.Errorf("Subject = %q", got)
	}
	from, err := mail.ParseAddress(fwd.header.Get("From"))
	if err != nil || from.Name != "Renamed Shop" {
		t.Errorf("From = %q (%v), want the rule's sender name", fwd.header.Get("From"), err)
	}
}

func TestForwardRules_RulesAppliedHeader(t *testing.T) {
	names := []string{"News", `Café "deals" <promo>, weekly`, strings.Repeat("ž", model.MaxRuleNameLength)}
	fwd := buildRulesForward(t, testRulesRawEmail, "Shop", plainRecipient, model.Settings{}, model.RuleOutcome{MatchedNames: names})

	raw := fwd.header.Get("X-Mailx-Rules-Applied")
	if strings.ContainsAny(raw, `<>"`) || strings.ContainsRune(raw, 'é') {
		t.Errorf("raw header %q has unencoded specials", raw)
	}

	values := strings.Split(raw, ", ")
	if len(values) != len(names) {
		t.Fatalf("header %q has %d values, want %d", raw, len(values), len(names))
	}
	for i, value := range values {
		decoded, err := new(mime.WordDecoder).DecodeHeader(value)
		if err != nil || decoded != names[i] {
			t.Errorf("value %d decodes to %q (%v), want %q", i, decoded, err, names[i])
		}
	}
	if values[0] != "News" {
		t.Errorf("plain name = %q, want it as is", values[0])
	}
}

func TestEncodeRuleName_WordsStayWithinLimit(t *testing.T) {
	name := strings.Repeat("ž", model.MaxRuleNameLength)
	encoded := encodeRuleName(name)

	for word := range strings.SplitSeq(encoded, " ") {
		if len(word) > 75 {
			t.Errorf("encoded-word %q is %d characters, want at most 75", word, len(word))
		}
	}
	if decoded, err := new(mime.WordDecoder).DecodeHeader(encoded); err != nil || decoded != name {
		t.Errorf("decoded = %q (%v)", decoded, err)
	}
}

func TestForwardRules_BannerOverride(t *testing.T) {
	const banner = "This email was sent to"
	hide, show := true, false

	fwd := buildRulesForward(t, testRulesRawEmail, "Shop", plainRecipient, model.RuleOutcome{RemoveHeader: &hide}.Settings(model.Settings{}), model.RuleOutcome{})
	if strings.Contains(fwd.html, banner) {
		t.Error("banner shown although a rule hid it")
	}

	fwd = buildRulesForward(t, testRulesRawEmail, "Shop", plainRecipient, model.RuleOutcome{RemoveHeader: &show}.Settings(model.Settings{RemoveHeader: true}), model.RuleOutcome{})
	if !strings.Contains(fwd.html, banner) {
		t.Error("banner hidden although a rule showed it")
	}
}

func TestForwardRules_RemoveAttachments(t *testing.T) {
	fwd := buildRulesForward(t, testRulesRawEmail, "Shop", plainRecipient, model.Settings{}, model.RuleOutcome{RemoveAttachments: true})

	if fwd.hasPart(func(p mimePart) bool { return p.filename == "deals.pdf" || p.contentType == "application/pdf" }) {
		t.Errorf("parts = %+v, want the attachment removed", fwd.parts)
	}
	if !fwd.hasPart(func(p mimePart) bool {
		return p.contentType == "image/png" && strings.Contains(p.contentID, "logo@shop.example")
	}) {
		t.Errorf("parts = %+v, want the inline cid: image kept", fwd.parts)
	}
	if !strings.Contains(fwd.html, "cid:logo@shop.example") {
		t.Error("HTML lost its cid: reference")
	}
}

func TestForwardRules_RemoveAttachmentsLeavesEncryptedEmail(t *testing.T) {
	raw := string(readFixture(t, "trackers_pgp_mime.eml"))
	before := buildRulesForward(t, raw, "Friend", plainRecipient, model.Settings{}, model.RuleOutcome{})
	after := buildRulesForward(t, raw, "Friend", plainRecipient, model.Settings{}, model.RuleOutcome{RemoveAttachments: true})

	if len(before.parts) == 0 || len(after.parts) != len(before.parts) {
		t.Errorf("parts before = %+v, after = %+v, want the encrypted email untouched", before.parts, after.parts)
	}
}

func TestForwardRules_DisableEncryption(t *testing.T) {
	publicKey := testPublicKey(t)

	for _, inline := range []bool{false, true} {
		rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey, PGPInline: inline}

		encrypted := buildRulesForward(t, testRulesRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{})
		if inline {
			if strings.Contains(encrypted.html, "Deals") {
				t.Errorf("inline: HTML forwarded although PGP/Inline is on")
			}
		} else if ct := encrypted.header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/encrypted") {
			t.Fatalf("Content-Type = %q, want PGP/MIME", ct)
		}

		plain := buildRulesForward(t, testRulesRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{DisableEncryption: true, MatchedNames: []string{"No PGP"}})
		if ct := plain.header.Get("Content-Type"); strings.HasPrefix(ct, "multipart/encrypted") {
			t.Errorf("inline=%v: Content-Type = %q, want no encryption", inline, ct)
		}
		if !strings.Contains(plain.html, "<p>Deals</p>") {
			t.Errorf("inline=%v: HTML = %q, want the plain body", inline, plain.html)
		}
	}
}

func TestForward_PGPInlineHasOneContentType(t *testing.T) {
	const plainText = "From: Friend <friend@example.com>\r\n" +
		"To: myalias@mailx.net\r\n" +
		"Subject: Hello\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"Hi there\r\n"
	rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: testPublicKey(t), PGPInline: true}

	fwd := buildRulesForward(t, plainText, "Friend", rcp, model.Settings{}, model.RuleOutcome{})
	if ct := fwd.header["Content-Type"]; len(ct) != 1 || !strings.HasPrefix(ct[0], "text/plain") {
		t.Errorf("Content-Type = %q, want one text/plain", ct)
	}

	fwd = buildRulesForward(t, testRulesRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{})
	if ct := fwd.header["Content-Type"]; len(ct) != 1 || !strings.HasPrefix(ct[0], "multipart/") {
		t.Errorf("with attachments: Content-Type = %q, want one multipart", ct)
	}
	if !fwd.hasPart(func(p mimePart) bool { return p.filename == "attachment-1.pdf.pgp" }) {
		t.Errorf("with attachments: parts = %+v, want the attachment", fwd.parts)
	}
}

// leafPart is a decoded non-multipart part of a built message.
type leafPart struct {
	header textproto.MIMEHeader
	data   []byte
}

func leafParts(t *testing.T, contentType string, encoding string, body io.Reader) []leafPart {
	t.Helper()

	mediaType, params, _ := mime.ParseMediaType(contentType)
	if !strings.HasPrefix(mediaType, "multipart/") {
		switch strings.ToLower(encoding) {
		case "base64":
			body = base64.NewDecoder(base64.StdEncoding, body)
		case "quoted-printable":
			body = quotedprintable.NewReader(body)
		}
		data, err := io.ReadAll(body)
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		return []leafPart{{textproto.MIMEHeader{"Content-Type": {contentType}}, data}}
	}

	var parts []leafPart
	mr := multipart.NewReader(body, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err != nil {
			return parts
		}
		for _, p := range leafParts(t, part.Header.Get("Content-Type"), part.Header.Get("Content-Transfer-Encoding"), part) {
			if !strings.HasPrefix(part.Header.Get("Content-Type"), "multipart/") {
				p.header = part.Header
			}
			parts = append(parts, p)
		}
	}
}

func decryptWith(t *testing.T, key *crypto.Key, data []byte, encoding int8) string {
	t.Helper()

	dec, err := crypto.PGP().Decryption().DecryptionKey(key).New()
	if err != nil {
		t.Fatalf("decryption handle: %v", err)
	}
	res, err := dec.Decrypt(data, encoding)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	return res.String()
}

func TestForward_PGPInlineEncryptsAttachments(t *testing.T) {
	key := testKey(t)
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}
	rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey, PGPInline: true}

	var buf bytes.Buffer
	if _, err := buildRulesForwardMessage(t, testRulesRawEmail, "Shop", rcp, model.Settings{}, model.RuleOutcome{}).WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	for _, name := range []string{"deals.pdf", "logo.png", "logo@shop.example"} {
		if bytes.Contains(buf.Bytes(), []byte(name)) {
			t.Errorf("raw forward contains %q in the clear", name)
		}
	}
	parsed, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("mail.ReadMessage() error = %v", err)
	}
	parts := leafParts(t, parsed.Header.Get("Content-Type"), parsed.Header.Get("Content-Transfer-Encoding"), parsed.Body)

	if len(parts) != 3 {
		t.Fatalf("got %d parts, want the body and 2 attachments", len(parts))
	}

	body := decryptWith(t, key, parts[0].data, crypto.Armor)
	for _, line := range []string{"attachment-1.pdf.pgp: deals.pdf", "attachment-2.png.pgp: logo.png (inline image)"} {
		if !strings.Contains(body, line) {
			t.Errorf("decrypted body = %q, want %q", body, line)
		}
	}

	want := []struct{ name, data string }{
		{"attachment-1.pdf.pgp", "%PDF-1.4"},
		{"attachment-2.png.pgp", "\x89PNG\r\n\x1a\n"},
	}
	for i, w := range want {
		p := parts[i+1]
		if ct := p.header.Get("Content-Type"); ct != `application/octet-stream; name="`+w.name+`"` {
			t.Errorf("part %d Content-Type = %q", i+1, ct)
		}
		if cd := p.header.Get("Content-Disposition"); cd != `attachment; filename="`+w.name+`"` {
			t.Errorf("part %d Content-Disposition = %q", i+1, cd)
		}
		if cid := p.header.Get("Content-Id"); cid != "" {
			t.Errorf("part %d Content-ID = %q, want none", i+1, cid)
		}
		if got := decryptWith(t, key, p.data, crypto.Bytes); got != w.data {
			t.Errorf("part %d decrypts to %q, want %q", i+1, got, w.data)
		}
	}
}

func TestPGPInlineFiles(t *testing.T) {
	attached := func(name string) letters.AttachedFile {
		return letters.AttachedFile{ContentDisposition: letters.ContentDispositionHeader{Params: map[string]string{"filename": name}}}
	}
	email := letters.Email{
		AttachedFiles: []letters.AttachedFile{
			attached("Report.PDF"),
			attached("no-extension"),
			attached("archive.tar.gz"),
			attached("odd.ext with space"),
			attached("long.verylongext"),
			attached("bad\r\nattachment-9.pgp: fake.pdf"),
		},
		InlineFiles: []letters.InlineFile{{ContentType: letters.ContentTypeHeader{ContentType: "image/png"}}},
	}

	files := pgpInlineFiles(email)
	wantNames := []string{
		"attachment-1.pdf.pgp",
		"attachment-2.pgp",
		"attachment-3.gz.pgp",
		"attachment-4.pgp",
		"attachment-5.pgp",
		"attachment-6.pdf.pgp",
		"attachment-7.pgp",
	}
	if len(files) != len(wantNames) {
		t.Fatalf("got %d files, want %d", len(files), len(wantNames))
	}
	for i, f := range files {
		if f.name != wantNames[i] {
			t.Errorf("file %d name = %q, want %q", i, f.name, wantNames[i])
		}
	}

	note := pgpInlineFilesNote(files)
	if !strings.Contains(note, "attachment-6.pdf.pgp: bad attachment-9.pgp: fake.pdf\n") {
		t.Errorf("note = %q, want the crafted name on one line", note)
	}
	if !strings.Contains(note, "attachment-7.pgp: image/png (inline image)\n") {
		t.Errorf("note = %q, want the inline image's type when it has no name", note)
	}
	if pgpInlineFilesNote(nil) != "" {
		t.Error("note added without attachments")
	}
}

func TestForwardRules_PGPMIMEGetsNewSubjectAndHeader(t *testing.T) {
	rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: testPublicKey(t)}
	outcome := model.RuleOutcome{MatchedNames: []string{"Deals"}, SubjectValues: []string{"[Deals] {subject}"}}

	fwd := buildRulesForward(t, testRulesRawEmail, "Shop", rcp, model.Settings{}, outcome)

	if ct := fwd.header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/encrypted") {
		t.Fatalf("Content-Type = %q, want PGP/MIME", ct)
	}
	if got := decodedHeader(t, fwd.header, "Subject"); got != "[Deals] Weekly deals" {
		t.Errorf("outer Subject = %q", got)
	}
	if got := fwd.header.Get("X-Mailx-Rules-Applied"); got != "Deals" {
		t.Errorf("X-Mailx-Rules-Applied = %q", got)
	}
}

func TestReplyRules(t *testing.T) {
	const raw = "From: Me <me@example.com>\r\n" +
		"To: myalias+support=vendor.example@mailx.net\r\n" +
		"Subject: Re: Ticket 7\r\n" +
		"In-Reply-To: <ticket7@vendor.example>\r\n" +
		"\r\n" +
		"Thanks\r\n"

	m := Mailer{}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	outcome := model.RuleOutcome{MatchedNames: []string{"Support"}, SubjectValues: []string{"{subject} [Customer]"}}

	msg, _, err := m.buildReplyMessage(alias.Name, "Customer", model.Recipient{Email: "support@vendor.example"}, []byte(raw), alias, outcome)
	if err != nil {
		t.Fatalf("buildReplyMessage() error = %v", err)
	}
	reply := writeAndParse(t, msg)

	if got := decodedHeader(t, reply.header, "Subject"); got != "Re: Ticket 7 [Customer]" {
		t.Errorf("Subject = %q", got)
	}
	if from, _ := mail.ParseAddress(reply.header.Get("From")); from == nil || from.Name != "Customer" {
		t.Errorf("From = %q", reply.header.Get("From"))
	}
	// Replies go to third parties, who must not see the user's rule names
	if _, ok := reply.header["X-Mailx-Rules-Applied"]; ok {
		t.Error("X-Mailx-Rules-Applied set on a reply")
	}
}

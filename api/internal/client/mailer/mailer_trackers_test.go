package mailer

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/mnako/letters"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils/trackers"
)

const testAppURL = "https://app.mailx.net"

type forwarded struct {
	header mail.Header
	html   string
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	return data
}

// originalHTML is the fixture's HTML body as forwarded today, without tracker removal.
func originalHTML(t *testing.T, name string) string {
	t.Helper()

	email, err := letters.ParseEmail(bytes.NewReader(readFixture(t, name)))
	if err != nil {
		t.Fatalf("ParseEmail() error = %v", err)
	}

	return email.HTML
}

func forwardFixture(t *testing.T, name string, settings model.Settings, rcp model.Recipient) forwarded {
	t.Helper()

	m := Mailer{cfg: config.SMTPClientConfig{AppURL: testAppURL}}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	templateData := map[string]any{"alias": alias.Name, "from": "news@shop.example"}

	msg, _, err := m.buildForwardMessage(
		"myalias+news=shop.example@mailx.net", "Sender", alias.Name, rcp,
		readFixture(t, name), "header.tmpl", templateData, settings, alias, model.ListUnsubscribeHeaders{},
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

	html, _ := findHTMLPart(t, parsed.Header.Get("Content-Type"), parsed.Header.Get("Content-Transfer-Encoding"), parsed.Body)

	return forwarded{header: parsed.Header, html: html}
}

func findHTMLPart(t *testing.T, contentType string, encoding string, body io.Reader) (string, bool) {
	t.Helper()

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", false
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err != nil {
				return "", false
			}
			// multipart.Part decodes quoted-printable itself and hides the header
			if html, ok := findHTMLPart(t, part.Header.Get("Content-Type"), part.Header.Get("Content-Transfer-Encoding"), part); ok {
				return html, true
			}
		}
	}
	if mediaType != "text/html" {
		return "", false
	}

	switch strings.ToLower(encoding) {
	case "quoted-printable":
		body = quotedprintable.NewReader(body)
	case "base64":
		body = base64.NewDecoder(base64.StdEncoding, body)
	}
	b, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read html part: %v", err)
	}

	// The SMTP transport writes CRLF line endings, with or without tracker removal
	return strings.ReplaceAll(string(b), "\r\n", "\n"), true
}

func reportFromHTML(t *testing.T, html string) (int, []trackers.Tracker) {
	t.Helper()

	prefix := `href="` + testAppURL + TrackersReportPath + "#"
	_, rest, ok := strings.Cut(html, prefix)
	if !ok {
		t.Fatalf("no report link in HTML: %s", html)
	}
	encoded, _, _ := strings.Cut(rest, `"`)

	total, removed, err := trackers.DecodeReport(encoded)
	if err != nil {
		t.Fatalf("DecodeReport(%q) error = %v", encoded, err)
	}

	return total, removed
}

func TestForwardTrackers_RemovedWhenOn(t *testing.T) {
	tests := []struct {
		fixture string
		gone    string
		keep    []string
		want    trackers.Tracker
	}{
		{
			fixture: "trackers_mailchimp.eml",
			gone:    "list-manage.com",
			keep: []string{
				`<img src="https://mcusercontent.com/123/images/logo.png" width="200" height="80" alt="Example Shop"/>`,
				`<img src="cid:banner@shop.example" width="600"/>`,
				`.hero{background:#fff url(https://mcusercontent.com/bg.png) no-repeat}`,
				`<body style="margin:0">`,
				"Grüße!",
			},
			want: trackers.Tracker{Domain: "shop.us1.list-manage.com", Element: trackers.ElementImg, Reason: trackers.ReasonKnownDomain},
		},
		{
			fixture: "trackers_sendgrid.eml",
			gone:    "ct.sendgrid.net",
			keep: []string{
				`<img src="https://cdn.mcauto-images-production.sendgrid.net/abc/logo.png" width="120"/>`,
				"<p>Welcome aboard!</p>",
			},
			want: trackers.Tracker{Domain: "u123456.ct.sendgrid.net", Element: trackers.ElementImg, Reason: trackers.ReasonKnownDomain},
		},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			fwd := forwardFixture(t, tt.fixture, model.Settings{BlockTrackers: true}, model.Recipient{Email: "real@domain.com"})

			if strings.Contains(fwd.html, tt.gone) {
				t.Errorf("HTML still contains the tracker %q", tt.gone)
			}
			for _, s := range tt.keep {
				if !strings.Contains(fwd.html, s) {
					t.Errorf("HTML lost %q:\n%s", s, fwd.html)
				}
			}
			if got := fwd.header.Get("X-Mailx-Trackers-Removed"); got != "1" {
				t.Errorf("X-Mailx-Trackers-Removed = %q, want 1", got)
			}
			if !strings.Contains(fwd.html, "1 tracker removed") {
				t.Errorf("banner doesn't say 1 tracker removed:\n%s", fwd.html)
			}
			if !strings.Contains(fwd.html, ">View report</a>") {
				t.Errorf("banner has no View report link")
			}

			total, removed := reportFromHTML(t, fwd.html)
			if total != 1 || len(removed) != 1 || removed[0] != tt.want {
				t.Errorf("report = %d %+v, want 1 [%+v]", total, removed, tt.want)
			}
		})
	}
}

func TestForwardTrackers_KeptWhenOff(t *testing.T) {
	for _, fixture := range []string{"trackers_mailchimp.eml", "trackers_sendgrid.eml"} {
		t.Run(fixture, func(t *testing.T) {
			fwd := forwardFixture(t, fixture, model.Settings{BlockTrackers: false}, model.Recipient{Email: "real@domain.com"})

			// Byte-for-byte what was forwarded before tracker blocking, after the banner
			if want := originalHTML(t, fixture); !strings.HasSuffix(fwd.html, want) {
				t.Errorf("HTML = %s\nwant it to end with the original HTML %s", fwd.html, want)
			}
			assertNoTrackerFeedback(t, fwd)
		})
	}
}

func TestForwardTrackers_NoneFoundLeavesHTMLUnchanged(t *testing.T) {
	fwd := forwardFixture(t, "trackers_none.eml", model.Settings{BlockTrackers: true}, model.Recipient{Email: "real@domain.com"})

	if want := originalHTML(t, "trackers_none.eml"); !strings.HasSuffix(fwd.html, want) {
		t.Errorf("HTML = %s\nwant it to end with the original HTML %s", fwd.html, want)
	}
	assertNoTrackerFeedback(t, fwd)
}

func TestForwardTrackers_AlreadyEncryptedUntouched(t *testing.T) {
	fwd := forwardFixture(t, "trackers_pgp_inline.eml", model.Settings{BlockTrackers: true}, model.Recipient{Email: "real@domain.com"})

	if want := originalHTML(t, "trackers_pgp_inline.eml"); !strings.HasSuffix(fwd.html, want) {
		t.Errorf("HTML = %s\nwant the encrypted email's HTML unchanged %s", fwd.html, want)
	}
	assertNoTrackerFeedback(t, fwd)

	fwd = forwardFixture(t, "trackers_pgp_mime.eml", model.Settings{BlockTrackers: true}, model.Recipient{Email: "real@domain.com"})
	assertNoTrackerFeedback(t, fwd)
}

func TestForwardTrackers_RemoveHeaderKeepsMailHeader(t *testing.T) {
	fwd := forwardFixture(t, "trackers_mailchimp.eml", model.Settings{BlockTrackers: true, RemoveHeader: true}, model.Recipient{Email: "real@domain.com"})

	if strings.Contains(fwd.html, "list-manage.com") {
		t.Error("HTML still contains the tracker")
	}
	if strings.Contains(fwd.html, "tracker removed") || strings.Contains(fwd.html, "View report") {
		t.Error("banner shown although the Mailx header is off")
	}
	if got := fwd.header.Get("X-Mailx-Trackers-Removed"); got != "1" {
		t.Errorf("X-Mailx-Trackers-Removed = %q, want 1", got)
	}
}

func TestForwardTrackers_PGPMIMEKeepsMailHeader(t *testing.T) {
	key, err := crypto.PGP().KeyGeneration().AddUserId("Test", "real@domain.com").New().GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}

	rcp := model.Recipient{Email: "real@domain.com", PGPEnabled: true, PGPKey: publicKey}
	fwd := forwardFixture(t, "trackers_mailchimp.eml", model.Settings{BlockTrackers: true}, rcp)

	if ct := fwd.header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/encrypted") {
		t.Fatalf("Content-Type = %q, want a PGP/MIME message", ct)
	}
	if got := fwd.header.Get("X-Mailx-Trackers-Removed"); got != "1" {
		t.Errorf("X-Mailx-Trackers-Removed = %q, want 1", got)
	}
}

func TestWithTrackersReport(t *testing.T) {
	removed := []trackers.Tracker{
		{Domain: "mailtrack.io", Element: trackers.ElementImg, Reason: trackers.ReasonKnownDomain},
		{Domain: "news.example.com", Element: trackers.ElementImg, Reason: trackers.ReasonTinyImage},
	}
	templateData := map[string]any{"alias": "a@mailx.net"}

	data := Mailer{}.withTrackersReport(templateData, removed)
	if data["trackersRemoved"] != 2 {
		t.Errorf("trackersRemoved = %v, want 2", data["trackersRemoved"])
	}
	if _, ok := data["trackersReportURL"]; ok {
		t.Error("trackersReportURL set without an app URL")
	}
	if _, ok := templateData["trackersRemoved"]; ok {
		t.Error("withTrackersReport() modified the caller's template data")
	}

	data = Mailer{cfg: config.SMTPClientConfig{AppURL: testAppURL}}.withTrackersReport(templateData, removed)
	if url, _ := data["trackersReportURL"].(string); !strings.HasPrefix(url, testAppURL+TrackersReportPath+"#") {
		t.Errorf("trackersReportURL = %q, want the report page with the report in the fragment", url)
	}
}

func assertNoTrackerFeedback(t *testing.T, fwd forwarded) {
	t.Helper()

	if got := fwd.header.Get("X-Mailx-Trackers-Removed"); got != "" {
		t.Errorf("X-Mailx-Trackers-Removed = %q, want no header", got)
	}
	if strings.Contains(fwd.html, "removed") || strings.Contains(fwd.html, "View report") {
		t.Errorf("banner mentions removed trackers:\n%s", fwd.html)
	}
}

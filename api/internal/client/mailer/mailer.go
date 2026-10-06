package mailer

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"maps"
	"mime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mnako/letters"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
	"ivpn.net/email/api/internal/utils/gomail.v2"
	"ivpn.net/email/api/internal/utils/trackers"
)

// TrackersReportPath is the app page that renders the removed trackers encoded in its fragment.
const TrackersReportPath = "/trackers-report"

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	cfg    config.SMTPClientConfig
}

// #nosec G104
func New(cfg config.SMTPClientConfig) Mailer {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Println("Invalid SMTP port:", cfg.Port)
		return Mailer{
			dialer: nil,
			cfg:    cfg,
		}
	}

	hosts := strings.Split(cfg.Host, ",")

	var dialer *gomail.Dialer
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if cfg.User == "" || cfg.Password == "" {
			dialer = &gomail.Dialer{Host: host, Port: port}
		} else {
			dialer = gomail.NewDialer(host, port, cfg.User, cfg.Password)
			if cfg.StartTLS == "opportunistic" {
				dialer.StartTLSPolicy = gomail.OpportunisticStartTLS
			}
		}

		conn, err := dialer.Dial()
		if err != nil {
			log.Printf("Failed to connect to SMTP host: %s, trying next host if available. Error: %v\n", host, err)
			continue
		} else {
			conn.Close()
			break
		}
	}

	if dialer == nil {
		return Mailer{
			dialer: nil,
			cfg:    cfg,
		}
	}

	return Mailer{
		dialer: dialer,
		cfg:    cfg,
	}
}

func (mailer Mailer) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailer.cfg.Sender, mailer.cfg.SenderName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.AddAlternative("text/html", body)

	err := mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func (mailer Mailer) Reply(from string, name string, rcp model.Recipient, data []byte, alias model.Alias, outcome model.RuleOutcome) error {
	m, messageID, err := mailer.buildReplyMessage(from, name, rcp, data, alias, outcome)
	if err != nil {
		return err
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Printf("Email reply sent successfully, %s", messageID)

	return nil
}

// buildReplyMessage builds the outgoing reply/send message without dialing the network, so it
// can be unit tested in isolation.
func (mailer Mailer) buildReplyMessage(from string, name string, rcp model.Recipient, data []byte, alias model.Alias, outcome model.RuleOutcome) (*gomail.Message, string, error) {
	email, err := parseEmail(data)
	if err != nil {
		return nil, "", err
	}

	email.Text = utils.RemoveHeader(email.Text)
	email.HTML = utils.RemoveHtmlHeader(email.HTML)

	if email.HTML == "" {
		email.HTML = model.PlainTextToHTML(email.Text)
	}

	// Decode the subject: letters may return an undecoded RFC 2047 value when
	// it cannot look up the charset (e.g. an empty-charset encoded-word that
	// slipped through preprocessing).  DecodeHeaderWithCharset is a no-op when
	// the subject is already plain text.
	decodedSubject := outcome.Subject(utils.DecodeHeaderWithCharset(email.Headers.Subject))

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", rcp.Email)
	m.SetHeader("Subject", decodedSubject)
	m.SetBody("text/plain", email.Text)
	m.AddAlternative("text/html", email.HTML)

	// Preserve original metadata
	if len(email.Headers.InReplyTo) > 0 {
		m.SetHeader("In-Reply-To", string(email.Headers.InReplyTo[0]))
	}
	m.SetHeader("X-Complaints-To", mailer.cfg.Report)
	m.SetHeader("X-Report-Abuse", mailer.cfg.Report)
	m.SetHeader("X-Report-Abuse-To", mailer.cfg.Report)
	m.SetHeader("Feedback-ID", fmt.Sprintf("mailx:%x:reply", sha256.Sum256([]byte(mailer.cfg.TokenSecret+alias.ID))))

	for _, a := range email.AttachedFiles {
		m.Attach(a.ContentDisposition.Params["filename"], gomail.SetCopyFunc(func(w io.Writer) error {
			_, err = w.Write(a.Data)
			return err
		}))
	}

	for _, f := range email.InlineFiles {
		m.Embed(
			f.ContentID,
			gomail.SetHeader(map[string][]string{
				"Content-ID":          {f.ContentID},
				"Content-Type":        {f.ContentType.ContentType},
				"Content-Disposition": {f.ContentDisposition.Params["type"] + "; filename=\"" + f.ContentDisposition.Params["filename"] + "\""},
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err = w.Write(f.Data)
				return err
			}),
		)
	}

	return m, string(email.Headers.MessageID), nil
}

// buildForwardMessage builds the outgoing forward message without dialing the network,
// so the To: header logic can be unit tested in isolation.
func (mailer Mailer) buildForwardMessage(from string, name string, to string, rcp model.Recipient, data []byte, templateFile string, templateData map[string]any, settings model.Settings, alias model.Alias, listUnsubscribe model.ListUnsubscribeHeaders, outcome model.RuleOutcome) (*gomail.Message, string, error) {
	if outcome.DisableEncryption {
		rcp.PGPEnabled = false
	}

	email, err := parseEmail(data)
	if err != nil {
		return nil, "", err
	}

	// The body of an email the sender encrypted can't be changed
	if outcome.RemoveAttachments && !isAlreadyEncrypted(email) {
		email.AttachedFiles = nil
	}

	removedTrackers := removeTrackers(&email, settings, rcp)
	templateData = mailer.withTrackersReport(templateData, removedTrackers)

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return nil, "", err
	}

	header := new(bytes.Buffer)
	if !settings.RemoveHeader {
		err = tmpl.ExecuteTemplate(header, "header", templateData)
		if err != nil {
			return nil, "", err
		}
	}

	headerHtml := new(bytes.Buffer)
	if !settings.RemoveHeader {
		err = tmpl.ExecuteTemplate(headerHtml, "headerHtml", templateData)
		if err != nil {
			return nil, "", err
		}
	}

	if email.HTML == "" {
		email.HTML = model.PlainTextToHTML(email.Text)
	}

	// Decode the subject: letters may return an undecoded RFC 2047 value when
	// it cannot look up the charset (e.g. an empty-charset encoded-word that
	// slipped through preprocessing).  DecodeHeaderWithCharset is a no-op when
	// the subject is already plain text.
	decodedSubject := outcome.Subject(utils.DecodeHeaderWithCharset(email.Headers.Subject))

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", decodedSubject)
	m.SetBody("text/plain", header.String()+email.Text)

	// Preserve original metadata
	originalSender := from
	if envFrom, ok := email.Headers.ExtraHeaders["Return-Path"]; ok && len(envFrom) > 0 {
		originalSender = envFrom[0]
	}
	m.SetHeader("X-Mailx-Original-Envelope-From", originalSender)
	m.SetHeader("X-Mailx-Original-Sender", originalSender)
	m.SetHeader("X-Mailx-Original-To", rcp.Email)
	if len(email.Headers.From) > 0 {
		m.SetHeader("X-Mailx-Original-From-Header", email.Headers.From[0].String())
	}
	if email.Headers.MessageID != "" {
		m.SetHeader("Message-ID", string(email.Headers.MessageID))
	}
	if authResults, ok := email.Headers.ExtraHeaders["Authentication-Results"]; ok && len(authResults) > 0 {
		m.SetHeader("X-Mailx-Authentication-Results", authResults...)
	}
	if len(email.Headers.InReplyTo) > 0 {
		m.SetHeader("In-Reply-To", string(email.Headers.InReplyTo[0]))
	}
	if refs, ok := email.Headers.ExtraHeaders["References"]; ok && len(refs) > 0 {
		m.SetHeader("References", refs...)
	}
	if received, ok := email.Headers.ExtraHeaders["Received"]; ok && len(received) > 0 {
		ownHop := fmt.Sprintf(
			"from unknown by %s with ESMTPS; %s",
			mailer.cfg.Host,
			time.Now().UTC().Format(time.RFC1123Z),
		)
		allReceived := append([]string{ownHop}, received...)
		m.SetHeader("Received", allReceived...)
	}
	arcHeaders := []string{
		"ARC-Seal",
		"ARC-Message-Signature",
		"ARC-Authentication-Results",
	}
	for _, h := range arcHeaders {
		if vals, ok := email.Headers.ExtraHeaders[h]; ok && len(vals) > 0 {
			m.SetHeader(h, vals...)
		}
	}
	m.SetHeader("X-Complaints-To", mailer.cfg.Report)
	m.SetHeader("X-Report-Abuse", mailer.cfg.Report)
	m.SetHeader("X-Report-Abuse-To", mailer.cfg.Report)
	m.SetHeader("Feedback-ID", fmt.Sprintf("mailx:%x:forward", sha256.Sum256([]byte(mailer.cfg.TokenSecret+alias.ID))))
	setListUnsubscribe(m, listUnsubscribe)
	setTrackersRemoved(m, removedTrackers)
	setRulesApplied(m, outcome.MatchedNames)

	// PGP/Inline encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && rcp.PGPInline {
		armored, err := utils.EncryptWithPGPInline(email.Text, rcp.PGPKey)
		if err != nil {
			return nil, "", err
		}
		m.SetHeader("Content-Type", "text/plain")
		m.SetBody("text/plain", armored)
	} else {
		m.AddAlternative("text/html", headerHtml.String()+email.HTML)
	}

	for _, a := range email.AttachedFiles {
		m.Attach(a.ContentDisposition.Params["filename"], gomail.SetCopyFunc(func(w io.Writer) error {
			_, err = w.Write(a.Data)
			return err
		}))
	}

	for _, f := range email.InlineFiles {
		m.Embed(
			f.ContentID,
			gomail.SetHeader(map[string][]string{
				"Content-ID":          {f.ContentID},
				"Content-Type":        {f.ContentType.ContentType},
				"Content-Disposition": {f.ContentDisposition.Params["type"] + "; filename=\"" + f.ContentDisposition.Params["filename"] + "\""},
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err = w.Write(f.Data)
				return err
			}),
		)
	}

	// PGP/MIME encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && !rcp.PGPInline {
		em, err := utils.EncryptWithPGPMIME(m, from, name, decodedSubject, to, rcp.PGPKey)
		if err != nil {
			return nil, "", err
		}
		// Kept on the outer message so mail clients can still offer unsubscribe
		setListUnsubscribe(em, listUnsubscribe)
		setTrackersRemoved(em, removedTrackers)
		setRulesApplied(em, outcome.MatchedNames)

		return em, string(email.Headers.MessageID), nil
	}

	return m, string(email.Headers.MessageID), nil
}

// setListUnsubscribe sets the List-Unsubscribe headers built for the forward. The original
// message's headers are never copied over, so they can't appear twice.
func setListUnsubscribe(m *gomail.Message, listUnsubscribe model.ListUnsubscribeHeaders) {
	if listUnsubscribe.Value == "" {
		return
	}

	m.SetHeader("List-Unsubscribe", listUnsubscribe.Value)
	if listUnsubscribe.Post != "" {
		m.SetHeader("List-Unsubscribe-Post", listUnsubscribe.Post)
	}
}

// removeTrackers strips tracking pixels from the HTML body in place when the user has tracker
// blocking on, and returns what was removed. On failure the HTML is forwarded unchanged.
func removeTrackers(email *letters.Email, settings model.Settings, rcp model.Recipient) []trackers.Tracker {
	// PGP/Inline forwards only the text body
	pgpInline := rcp.PGPEnabled && rcp.PGPKey != "" && rcp.PGPInline
	if !settings.BlockTrackers || email.HTML == "" || pgpInline || isAlreadyEncrypted(*email) {
		return nil
	}

	res, err := trackers.Remove(email.HTML)
	if err != nil {
		log.Printf("Warning: forwarding HTML without removing trackers: %v", err)
		return nil
	}
	email.HTML = res.HTML

	return res.Removed
}

// isAlreadyEncrypted reports whether the sender encrypted the message, so its body can't be changed.
func isAlreadyEncrypted(email letters.Email) bool {
	switch strings.ToLower(email.Headers.ContentType.ContentType) {
	case "multipart/encrypted", "application/pkcs7-mime", "application/x-pkcs7-mime":
		return true
	}

	const armor = "-----BEGIN PGP MESSAGE-----"
	return strings.Contains(email.Text, armor) || strings.Contains(email.HTML, armor)
}

// withTrackersReport adds the removed tracker count and report link for the header template.
func (mailer Mailer) withTrackersReport(templateData map[string]any, removed []trackers.Tracker) map[string]any {
	if len(removed) == 0 {
		return templateData
	}

	data := maps.Clone(templateData)
	if data == nil {
		data = map[string]any{}
	}
	data["trackersRemoved"] = len(removed)
	if mailer.cfg.AppURL != "" {
		// The fragment is never sent to the server, so the report stays in the email
		data["trackersReportURL"] = mailer.cfg.AppURL + TrackersReportPath + "#" + trackers.EncodeReport(removed)
	}

	return data
}

func setTrackersRemoved(m *gomail.Message, removed []trackers.Tracker) {
	if len(removed) == 0 {
		return
	}

	m.SetHeader("X-Mailx-Trackers-Removed", strconv.Itoa(len(removed)))
}

// setRulesApplied lists the rules that matched a forward, so users can filter on them in their
// mail client. It's never set on replies or sends, which go to third parties.
func setRulesApplied(m *gomail.Message, names []string) {
	if len(names) == 0 {
		return
	}

	encoded := make([]string, len(names))
	for i, name := range names {
		encoded[i] = encodeRuleName(name)
	}

	// Each value is joined with ", " and folded by the writer
	m.SetHeader("X-Mailx-Rules-Applied", encoded...)
}

// ruleNameSpecials are characters that are encoded in a rule name, so a name can't read as
// several names, or as markup, to a mail client.
const ruleNameSpecials = `,;<>"'()[]\=?`

// encodeRuleName RFC 2047-encodes a rule name that has non-ASCII text or ruleNameSpecials, as
// base64 encoded-words of at most 75 characters split on rune boundaries. Other names are used
// as they are.
func encodeRuleName(name string) string {
	if !strings.ContainsAny(name, ruleNameSpecials) && mime.QEncoding.Encode("UTF-8", name) == name {
		return name
	}

	// "=?UTF-8?b?" + "?=" leaves 63 characters, which is 45 bytes of base64 input
	const maxChunk = 45

	var words []string
	for len(name) > 0 {
		n := 0
		for n < len(name) {
			_, size := utf8.DecodeRuneInString(name[n:])
			if n > 0 && n+size > maxChunk {
				break
			}
			n += size
		}
		words = append(words, "=?UTF-8?b?"+base64.StdEncoding.EncodeToString([]byte(name[:n]))+"?=")
		name = name[n:]
	}

	return strings.Join(words, " ")
}

func (mailer Mailer) Forward(from string, name string, to string, rcp model.Recipient, data []byte, templateFile string, templateData map[string]any, settings model.Settings, alias model.Alias, listUnsubscribe model.ListUnsubscribeHeaders, outcome model.RuleOutcome) error {
	m, messageID, err := mailer.buildForwardMessage(from, name, to, rcp, data, templateFile, templateData, settings, alias, listUnsubscribe, outcome)
	if err != nil {
		return err
	}

	sc, err := mailer.dialer.Dial()
	if err != nil {
		return err
	}
	defer sc.Close()

	// Envelope recipient stays the real mailbox even though To: shows the alias.
	if err := sc.Send(from, []string{rcp.Email}, m); err != nil {
		return err
	}

	log.Printf("Email forward sent successfully, %s", messageID)
	return nil
}

func (mailer Mailer) SendTemplate(to string, subject string, templateFile string, templateData any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", templateData)
	if err != nil {
		return err
	}

	bodyHtml := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(bodyHtml, "bodyHtml", templateData)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailer.cfg.Sender, mailer.cfg.SenderName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body.String())
	m.AddAlternative("text/html", bodyHtml.String())

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email template sent successfully")
	return nil
}

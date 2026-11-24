package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"regexp"
	"strings"

	"ivpn.net/email/api/internal/utils"
)

var (
	ErrExtractOriginalFrom = fmt.Errorf("error extracting original From from bounce")
	replySubjectRE         = regexp.MustCompile(`(?i)^\s*(re|aw|antw|sv|rif|回复|回覆)\s*:\s*`)
	receivedForRE          = regexp.MustCompile(`(?i)\bfor\s+<?([^\s<>;,]+@[^\s<>;,]+)>?`)
)

type Msg struct {
	From                string
	FromName            string
	To                  []string
	Subject             string
	Body                string
	Type                MessageType
	EnvelopeRecipient   string
	ListUnsubscribe     string
	ListUnsubscribePost string
	ListID              string
}

// IsMailingList reports whether the message came from a mailing list or newsletter.
func (m Msg) IsMailingList() bool {
	return m.ListUnsubscribe != "" || m.ListID != ""
}

func ParseMsg(data []byte) (Msg, error) {
	// Preprocess email data to decode RFC 2047 encoded headers
	processedData, err := utils.PreprocessEmailData(data)
	if err != nil {
		log.Printf("Warning: failed to preprocess email data: %v", err)
		processedData = data // Fallback to original data
	}

	msg, err := mail.ReadMessage(bytes.NewReader(processedData))
	if err != nil {
		return Msg{}, err
	}

	subject := utils.DecodeHeaderWithCharset(msg.Header.Get("Subject"))

	to, err := parseToHeader(utils.DecodeHeaderWithCharset(msg.Header.Get("To")))
	if err != nil {
		return Msg{}, fmt.Errorf("error parsing To header: %w", err)
	}

	deliveredTo := ""
	for _, h := range []string{"Delivered-To", "X-Original-To"} {
		raw := strings.TrimSpace(msg.Header.Get(h))
		if raw == "" {
			continue
		}

		if addr, err := mail.ParseAddress(raw); err == nil {
			deliveredTo = addr.Address
		} else {
			deliveredTo = raw
		}
		break
	}
	if deliveredTo == "" {
		// Some pipe transports don't stamp Delivered-To/X-Original-To; fall back to
		// the "for <address>" clause Postfix adds to its own (topmost) Received
		// header when that hop had a single envelope recipient.
		if m := receivedForRE.FindStringSubmatch(msg.Header.Get("Received")); len(m) == 2 {
			deliveredTo = m[1]
		}
	}

	from, err := mail.ParseAddress(utils.DecodeHeaderWithCharset(msg.Header.Get("From")))
	if err != nil {
		return Msg{To: to, EnvelopeRecipient: deliveredTo}, fmt.Errorf("error parsing From header: %w", err)
	}
	fromAddress := from.Address

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(msg.Body)
	if err != nil {
		return Msg{To: to, EnvelopeRecipient: deliveredTo}, err
	}
	body := buf.String()
	msgType := Send

	if isReply(msg) {
		msgType = Reply
	}

	if isBounce(msg) {
		msgType = FailBounce
		fromAddress, err = ExtractOriginalFrom(processedData)
		if err != nil {
			log.Println("error extracting original From from bounce:", err)
			return Msg{To: to, EnvelopeRecipient: deliveredTo}, ErrExtractOriginalFrom
		}
	}

	return Msg{
		From:                fromAddress,
		FromName:            from.Name,
		To:                  to,
		Subject:             subject,
		Body:                body,
		Type:                msgType,
		EnvelopeRecipient:   deliveredTo,
		ListUnsubscribe:     strings.TrimSpace(msg.Header.Get("List-Unsubscribe")),
		ListUnsubscribePost: strings.TrimSpace(msg.Header.Get("List-Unsubscribe-Post")),
		ListID:              strings.TrimSpace(msg.Header.Get("List-Id")),
	}, nil
}

// parseToHeader returns the addresses in a decoded To header. A missing or blank
// header, or an empty group such as "undisclosed-recipients:;" (Bcc-only sends,
// some newsletters), is valid mail and yields no addresses, leaving the envelope
// recipient to pick the alias; net/mail would otherwise reject the blank header
// with "mail: no address". The header is parsed as-is first so group syntax and
// quoted semicolons survive, then retried with Outlook-style ";" separators
// turned into commas.
func parseToHeader(header string) ([]string, error) {
	to := make([]string, 0)
	if strings.TrimSpace(header) == "" {
		return to, nil
	}

	addresses, err := mail.ParseAddressList(header)
	if err != nil {
		normalized := utils.NormalizeAddressSeparators(header)
		if normalized == "" {
			return to, nil
		}
		addresses, err = mail.ParseAddressList(normalized)
		if err != nil {
			return nil, err
		}
	}

	for _, address := range addresses {
		to = append(to, address.Address)
	}
	return to, nil
}

// isReply checks whether the given email is a reply.
func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	subject := m.Header.Get("Subject")
	if replySubjectRE.MatchString(subject) {
		return true
	}

	return false
}

// isBounce checks whether the given email is a bounce (DSN).
func isBounce(m *mail.Message) bool {
	// 1. Check Return-Path header
	if rp := m.Header.Get("Return-Path"); strings.TrimSpace(rp) == "<>" {
		return true
	}

	// 2. Check Content-Type header
	ct := m.Header.Get("Content-Type")
	mediatype, params, err := mime.ParseMediaType(ct)
	if err == nil && strings.EqualFold(mediatype, "multipart/report") {
		if strings.EqualFold(params["report-type"], "delivery-status") {
			return true
		}
	}

	// 3. Optional: check Auto-Submitted
	if strings.EqualFold(m.Header.Get("Auto-Submitted"), "auto-replied") {
		return true
	}

	return false
}

// originalMessageTypes are the part types a bounce (DSN) uses to return the
// original message: the whole message (RFC 2046, RFC 6532) or only its header
// (RFC 6522, RFC 6533), e.g. when the sender asked for RET=HDRS.
var originalMessageTypes = map[string]bool{
	"message/rfc822":         true,
	"message/global":         true,
	"text/rfc822-headers":    true,
	"message/global-headers": true,
}

var errNoOriginalMessage = errors.New("no message/rfc822 or text/rfc822-headers part found")

// ExtractOriginalFrom parses a bounce (DSN) and returns the From address of
// the original message it returns. The original message, or only its header,
// can be at any depth up to maxMIMEDepth, e.g. a multipart/report inside a
// multipart/mixed. The first one found is used.
func ExtractOriginalFrom(data []byte) (string, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("read message: %w", err)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("parse content-type: %w", err)
	}

	// Must be multipart with a boundary.
	if !strings.HasPrefix(mediaType, "multipart/") || params["boundary"] == "" {
		return "", fmt.Errorf("not a multipart/report message")
	}

	original, err := findOriginalMessage(msg.Body, params["boundary"], 0)
	if err != nil {
		return "", err
	}

	// A header-only part may end without the blank line that ends a header,
	// which ReadMessage accepts, but must not start with one.
	innerMsg, err := mail.ReadMessage(bytes.NewReader(bytes.TrimLeft(original, "\r\n")))
	if err != nil {
		return "", fmt.Errorf("read inner message: %w", err)
	}

	from, err := mail.ParseAddress(innerMsg.Header.Get("From"))
	if err != nil {
		return "", fmt.Errorf("parse original From: %w", err)
	}
	return from.Address, nil
}

// findOriginalMessage walks the multipart body depth first and returns the
// content of the first part in originalMessageTypes. It doesn't look inside
// that part, so a message attached to the original message is never used.
// A nested multipart that can't be read is skipped.
func findOriginalMessage(body io.Reader, boundary string, depth int) ([]byte, error) {
	if boundary == "" || depth > maxMIMEDepth {
		return nil, errNoOriginalMessage
	}

	mr := multipart.NewReader(body, boundary)
	for {
		// NextPart decodes quoted-printable parts.
		part, err := mr.NextPart()
		if err == io.EOF {
			return nil, errNoOriginalMessage
		}
		if err != nil {
			return nil, fmt.Errorf("next part: %w", err)
		}

		// ParseMediaType returns the media type even if a parameter is invalid.
		mediaType, params, _ := mime.ParseMediaType(part.Header.Get("Content-Type"))

		switch {
		case originalMessageTypes[mediaType]:
			var r io.Reader = part
			if strings.EqualFold(strings.TrimSpace(part.Header.Get("Content-Transfer-Encoding")), "base64") {
				r = base64.NewDecoder(base64.StdEncoding, part)
			}
			content, err := io.ReadAll(r)
			if err != nil {
				return nil, fmt.Errorf("read inner part: %w", err)
			}
			return content, nil
		case strings.HasPrefix(mediaType, "multipart/"):
			if content, err := findOriginalMessage(part, params["boundary"], depth+1); err == nil {
				return content, nil
			}
		}
	}
}

package model

import (
	"fmt"
	"strings"
	"testing"
)

// FuzzParseMsg fuzzes the top-level entry point for inbound mail: every byte
// that Postfix pipes to POST /email is handed to ParseMsg (see
// internal/service/processor.go ProcessMessage). Any panic here is remotely
// triggerable by any external sender relaying mail to a domain hosted on
// this service, without authentication.
func FuzzParseMsg(f *testing.F) {
	seeds := []string{
		// Plain minimal message.
		"From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody.",
		// Reply, via In-Reply-To.
		"From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Re: Test\r\nIn-Reply-To: <id@example.com>\r\n\r\nBody.",
		// Multiple To addresses, semicolon separated (Outlook style).
		"From: sender@example.com\r\nTo: a@example.com; b@example.com;\r\nSubject: Test\r\n\r\nBody.",
		// Delivered-To / X-Original-To fallback.
		"From: sender@example.com\r\nTo: a@example.com\r\nDelivered-To: b@example.com\r\nSubject: Test\r\n\r\nBody.",
		"From: sender@example.com\r\nTo: a@example.com\r\nX-Original-To: <b@example.com>\r\nSubject: Test\r\n\r\nBody.",
		// Received "for <addr>" fallback.
		"From: sender@example.com\r\nTo: a@example.com\r\nReceived: from x by y for <c@example.com>; Mon, 1 Jan 2024 00:00:00 +0000\r\nSubject: Test\r\n\r\nBody.",
		// Bounce via empty Return-Path, multipart/report with embedded message/rfc822.
		"From: mailer-daemon@example.com\r\nTo: a@example.com\r\nReturn-Path: <>\r\nContent-Type: multipart/report; report-type=delivery-status; boundary=BND\r\n\r\n--BND\r\nContent-Type: text/plain\r\n\r\nFailure notice.\r\n\r\n--BND\r\nContent-Type: message/rfc822\r\n\r\nFrom: original@example.com\r\nTo: a@example.com\r\nSubject: Original\r\n\r\nOriginal body.\r\n--BND--\r\n",
		// Bounce via Auto-Submitted.
		"From: mailer-daemon@example.com\r\nTo: a@example.com\r\nAuto-Submitted: auto-replied\r\nSubject: Bounce\r\n\r\nBody.",
		// Malformed To header (should surface parseErr path).
		"From: sender@example.com\r\nTo: not-an-address\r\nSubject: Test\r\n\r\nBody.",
		// Missing From header entirely.
		"To: a@example.com\r\nSubject: Test\r\n\r\nBody.",
		// RFC 2047 encoded subject/from, several charsets.
		"From: =?UTF-8?B?SsO6bGl1cw==?= <sender@example.com>\r\nTo: a@example.com\r\nSubject: =?ISO-8859-1?Q?Caf=E9?=\r\n\r\nBody.",
		// Empty-charset encoded word (should be fixed up by PreprocessEmailData).
		"From: =??Q?broken?= <sender@example.com>\r\nTo: a@example.com\r\nSubject: =??B?dGVzdA==?=\r\n\r\nBody.",
		// No body / no blank line separator at all.
		"From: sender@example.com\r\nTo: a@example.com\r\nSubject: Test",
		// CRLF-less (LF only) headers.
		"From: sender@example.com\nTo: a@example.com\nSubject: Test\n\nBody.",
		// Non-UTF8 bytes in header value and body.
		"From: sender@example.com\r\nTo: a@example.com\r\nSubject: \xff\xfe\x00bad\r\n\r\n\x80\x81\x82",
		// Duplicate headers.
		"From: a@example.com\r\nFrom: b@example.com\r\nTo: c@example.com\r\nTo: d@example.com\r\nSubject: dup\r\n\r\nBody.",
		// International subject prefixes for isReply.
		"From: sender@example.com\r\nTo: a@example.com\r\nSubject: 回复: hi\r\n\r\nBody.",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	// Deeply nested multipart/report bounce, to probe recursive descent in
	// ExtractOriginalFrom / multipart.Reader boundary handling.
	nested := "From: original@example.com\r\nTo: a@example.com\r\nSubject: inner\r\n\r\nleaf"
	for i := 0; i < 30; i++ {
		boundary := fmt.Sprintf("B%d", i)
		nested = fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n--%s\r\n%s\r\n--%s--\r\n", boundary, boundary, nested, boundary)
	}
	bounceNested := "From: mailer-daemon@example.com\r\nTo: a@example.com\r\nReturn-Path: <>\r\nContent-Type: multipart/report; report-type=delivery-status; boundary=OUTER\r\n\r\n--OUTER\r\nContent-Type: text/plain\r\n\r\nnotice\r\n\r\n--OUTER\r\nContent-Type: message/rfc822\r\n\r\n" + nested + "\r\n--OUTER--\r\n"
	f.Add([]byte(bounceNested))

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("ParseMsg panicked on input (len=%d): %v\ninput=%q", len(data), r, truncate(data, 500))
			}
		}()
		_, _ = ParseMsg(data)
	})
}

// FuzzExtractOriginalFrom fuzzes the bounce/DSN embedded-message extraction
// path, which walks a multipart/report structure and recursively re-parses
// an inner message/rfc822 part taken directly from attacker-controlled bytes.
func FuzzExtractOriginalFrom(f *testing.F) {
	seeds := []string{
		"Content-Type: multipart/report; report-type=delivery-status; boundary=BND\r\n\r\n--BND\r\nContent-Type: text/plain\r\n\r\nnotice\r\n\r\n--BND\r\nContent-Type: message/rfc822\r\n\r\nFrom: original@example.com\r\nTo: a@example.com\r\n\r\nBody\r\n--BND--\r\n",
		// No boundary param.
		"Content-Type: multipart/report; report-type=delivery-status\r\n\r\nbody",
		// Not multipart at all.
		"Content-Type: text/plain\r\n\r\nbody",
		// Boundary present but body never uses it (multipart.Reader must not hang).
		"Content-Type: multipart/mixed; boundary=X\r\n\r\nno boundary markers here at all just plain text that keeps going and going",
		// message/rfc822 part with malformed inner From.
		"Content-Type: multipart/report; boundary=BND\r\n\r\n--BND\r\nContent-Type: message/rfc822\r\n\r\nFrom: not-an-address\r\n\r\nBody\r\n--BND--\r\n",
		// Empty message.
		"",
		// Header only, no body/boundary.
		"Content-Type: multipart/report; boundary=\r\n\r\n",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("ExtractOriginalFrom panicked on input (len=%d): %v\ninput=%q", len(data), r, truncate(data, 500))
			}
		}()
		_, _ = ExtractOriginalFrom(data)
	})
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) <= n {
		return s
	}
	return strings.ToValidUTF8(s[:n], "�") + "...(truncated)"
}

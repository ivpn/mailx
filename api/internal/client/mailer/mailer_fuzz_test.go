package mailer

import (
	"bytes"
	"fmt"
	"net/mail"
	"testing"

	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/utils"
)

// newForwardParser mirrors the letters.EmailParser construction used in
// production by Mailer.Reply and Mailer.buildForwardMessage, so the fuzz
// target exercises exactly the code path attacker-controlled inbound
// message bytes flow through.
func newForwardParser() *letters.EmailParser {
	return letters.NewEmailParser(
		letters.WithToHeaderParser(
			func(header mail.Header, s string) ([]*mail.Address, error) {
				return nil, nil
			},
		),
		letters.WithReplyToHeaderParser(
			func(header mail.Header, s string) ([]*mail.Address, error) {
				return nil, nil
			},
		),
	)
}

// FuzzLettersParseEmail fuzzes the vendored third-party MIME parser
// (github.com/mnako/letters) exactly as it is invoked on raw, unauthenticated
// inbound message bytes in Mailer.Reply / Mailer.buildForwardMessage — i.e.
// on every message Postfix pipes to POST /email once SPF/DKIM/DMARC checks
// pass. A crash here is remotely triggerable and, since the API server has
// no panic-recovery middleware (see internal/transport/api/server.go), a
// panic or unrecoverable runtime crash in this path takes down the whole
// process, not just the one request.
func FuzzLettersParseEmail(f *testing.F) {
	plain := "From: a@example.com\r\nTo: b@example.com\r\nSubject: Test\r\n\r\nBody."
	multipartMixed := "From: a@example.com\r\nTo: b@example.com\r\nSubject: Test\r\n" +
		"Content-Type: multipart/mixed; boundary=BND\r\n\r\n" +
		"--BND\r\nContent-Type: text/plain\r\n\r\nplain part\r\n" +
		"--BND\r\nContent-Type: text/html\r\n\r\n<p>html part</p>\r\n" +
		"--BND\r\nContent-Type: application/octet-stream; name=a.bin\r\nContent-Disposition: attachment; filename=a.bin\r\nContent-Transfer-Encoding: base64\r\n\r\nQUJDRA==\r\n" +
		"--BND--\r\n"
	inlineImage := "From: a@example.com\r\nTo: b@example.com\r\nSubject: Test\r\n" +
		"Content-Type: multipart/related; boundary=BND\r\n\r\n" +
		"--BND\r\nContent-Type: text/html\r\n\r\n<img src=cid:img1>\r\n" +
		"--BND\r\nContent-Type: image/png\r\nContent-ID: <img1>\r\nContent-Disposition: inline; filename=a.png\r\nContent-Transfer-Encoding: base64\r\n\r\niVBORw0KGgo=\r\n" +
		"--BND--\r\n"
	badContentType := "From: a@example.com\r\nTo: b@example.com\r\nContent-Type: multipart/mixed\r\n\r\nno boundary param, body follows"
	badCTE := "From: a@example.com\r\nTo: b@example.com\r\nContent-Transfer-Encoding: not-a-real-encoding\r\n\r\nBody."
	malformedBase64 := "From: a@example.com\r\nTo: b@example.com\r\nContent-Transfer-Encoding: base64\r\n\r\n!!!not-valid-base64!!!"
	rfc2047EmptyCharset := "From: =??Q?broken?= <a@example.com>\r\nTo: b@example.com\r\nSubject: =??B?dGVzdA==?=\r\n\r\nBody."
	nonUTF8 := "From: a@example.com\r\nTo: b@example.com\r\nSubject: \xff\xfe\r\n\r\n\x80\x81\x82"
	emptyBoundaryLoop := "From: a@example.com\r\nTo: b@example.com\r\nContent-Type: multipart/mixed; boundary=\r\n\r\nbody"

	for _, s := range []string{
		plain, multipartMixed, inlineImage, badContentType, badCTE,
		malformedBase64, rfc2047EmptyCharset, nonUTF8, emptyBoundaryLoop,
		"", "garbage, not an email at all",
	} {
		f.Add([]byte(s))
	}

	// Deeply nested multipart/mixed structure ("MIME bomb" shape): each
	// level wraps the previous one in a new multipart/mixed part. The
	// letters.parsePart function recurses once per nesting level with no
	// depth limit. Parse time scales at least quadratically with depth
	// (measured: depth=1600 -> ~1s, depth=3200 -> ~4.2s, depth=5000 -> hangs
	// the Go fuzz engine's baseline-coverage watchdog outright — see
	// docs/fuzzing-findings.md finding F-1). Only a small depth is kept
	// here so this seed corpus entry stays fast for plain `go test`; the
	// quadratic blowup is reproduced separately in the findings report.
	for _, depth := range []int{50} {
		leaf := "From: a@example.com\r\nTo: b@example.com\r\nSubject: leaf\r\n\r\nleaf body"
		nested := leaf
		for i := 0; i < depth; i++ {
			boundary := fmt.Sprintf("B%d", i)
			nested = fmt.Sprintf(
				"Content-Type: multipart/mixed; boundary=%s\r\n\r\n--%s\r\n%s\r\n--%s--\r\n",
				boundary, boundary, nested, boundary,
			)
		}
		full := "From: a@example.com\r\nTo: b@example.com\r\nSubject: nested\r\n" + nested
		f.Add([]byte(full))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("letters.ParseEmail panicked on input (len=%d): %v", len(data), r)
			}
		}()

		processed, err := utils.PreprocessEmailData(data)
		if err != nil {
			processed = data
		}

		parser := newForwardParser()
		_, _ = parser.Parse(bytes.NewReader(processed))
	})
}

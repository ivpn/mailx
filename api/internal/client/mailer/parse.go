package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/mail"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/utils"
)

const (
	// maxMIMEDepth is the deepest multipart nesting parsed. Real messages nest a few levels,
	// while letters parses each level in time that grows with the depth.
	maxMIMEDepth = 30

	// parseTimeout bounds how long letters may parse a message, in case some other shape of
	// message is also slow to parse.
	parseTimeout = 10 * time.Second
)

var (
	ErrMIMETooDeep  = errors.New("message has too deeply nested MIME parts")
	ErrParseTimeout = errors.New("timed out parsing message")
)

// parseEmail parses inbound message bytes for a reply or forward. Parsing is bounded, so a
// crafted message can't hold up delivery or make letters crash the server.
func parseEmail(data []byte) (letters.Email, error) {
	// Preprocess email data to decode RFC 2047 encoded headers
	processedData, err := utils.PreprocessEmailData(data)
	if err != nil {
		log.Printf("Warning: failed to preprocess email data: %v", err)
		processedData = data // Fallback to original data
	}

	if mimeDepthExceeds(processedData, maxMIMEDepth) {
		return letters.Email{}, ErrMIMETooDeep
	}

	parser := letters.NewEmailParser(
		// Malformed headers can cause the parser to fail,
		// so we skip parsing of headers that are not needed for forwarding
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

	return parseWithTimeout(func() (letters.Email, error) {
		return parser.Parse(bytes.NewReader(processedData))
	}, parseTimeout)
}

// parseWithTimeout runs parse, returning ErrParseTimeout if it takes longer than timeout and an
// error if it panics. A parse that times out can't be stopped, so it finishes in the background.
func parseWithTimeout(parse func() (letters.Email, error), timeout time.Duration) (letters.Email, error) {
	type result struct {
		email letters.Email
		err   error
	}

	// Buffered, so the goroutine can always send and exit after a timeout
	done := make(chan result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic parsing message: %v\n%s", r, debug.Stack())
				done <- result{err: fmt.Errorf("panic parsing message: %v", r)}
			}
		}()

		email, err := parse()
		done <- result{email, err}
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case res := <-done:
		return res.email, res.err
	case <-timer.C:
		return letters.Email{}, ErrParseTimeout
	}
}

// mimeDepthExceeds reports whether the message has multipart parts nested deeper than max. It
// walks the parts as letters does, but stops past max, so it's fast for any message. A message
// that can't be read is left for letters to reject.
func mimeDepthExceeds(data []byte, max int) bool {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return false
	}

	return partDepthExceeds(msg.Body, msg.Header.Get("Content-Type"), 1, max)
}

func partDepthExceeds(r io.Reader, contentType string, depth int, max int) bool {
	ct, err := letters.ParseContentTypeHeader(contentType)
	if err != nil || !strings.HasPrefix(ct.ContentType, "multipart/") {
		return false
	}
	if depth > max {
		return true
	}

	mr := multipart.NewReader(r, ct.Params["boundary"])
	for {
		part, err := mr.NextPart()
		if err != nil {
			return false
		}
		if partDepthExceeds(part, part.Header.Get("Content-Type"), depth+1, max) {
			return true
		}
	}
}

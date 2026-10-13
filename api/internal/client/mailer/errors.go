package mailer

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"regexp"

	"ivpn.net/email/api/internal/utils/gomail.v2"
)

var (
	// smtpReply matches an SMTP reply code and optional enhanced status code, e.g. "550 5.1.1",
	// at the start of an error or after a wrapping prefix ("gomail: could not send email 1: ").
	// textproto.Error quotes the reply text, hence the optional quote.
	smtpReply    = regexp.MustCompile(`(?:^|: )([2-5]\d\d)(?:[ -]"?([2-5]\.\d{1,3}\.\d{1,3}))?\b`)
	enhancedCode = regexp.MustCompile(`^([2-5]\.\d{1,3}\.\d{1,3})\b`)
)

// ErrorClass describes an error from Forward, Reply or Send for the API log without its text.
// An SMTP server's reply often echoes the recipient's address ("550 5.1.1 <user@example.com>:
// Recipient address rejected"), so only the status codes are kept.
func ErrorClass(err error) string {
	if err == nil {
		return ""
	}

	var tpErr *textproto.Error
	if errors.As(err, &tpErr) {
		if m := enhancedCode.FindStringSubmatch(tpErr.Msg); m != nil {
			return fmt.Sprintf("smtp %03d %s", tpErr.Code, m[1])
		}
		return fmt.Sprintf("smtp %03d", tpErr.Code)
	}

	// gomail.Send wraps with %v, so the reply is only left in the text.
	if m := smtpReply.FindStringSubmatch(err.Error()); m != nil {
		if m[2] != "" {
			return "smtp " + m[1] + " " + m[2]
		}
		return "smtp " + m[1]
	}

	switch {
	case errors.Is(err, ErrMIMETooDeep), errors.Is(err, ErrParseTimeout), errors.Is(err, gomail.ErrStartTLSUnsupported):
		return err.Error()
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "network timeout"
		}
		return "network error"
	}

	return "unclassified error"
}

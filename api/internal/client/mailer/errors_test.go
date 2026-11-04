package mailer

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"testing"

	"ivpn.net/email/api/internal/utils/gomail.v2"
)

func TestErrorClass(t *testing.T) {
	rejected := &textproto.Error{Code: 550, Msg: "5.1.1 <user@real-inbox.example>: Recipient address rejected"}

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, ""},
		{"smtp error", rejected, "smtp 550 5.1.1"},
		{"wrapped smtp error", fmt.Errorf("send: %w", rejected), "smtp 550 5.1.1"},
		{"smtp error without enhanced code", &textproto.Error{Code: 421, Msg: "Service not available"}, "smtp 421"},
		{"gomail send text", fmt.Errorf("gomail: could not send email 1: %v", rejected), "smtp 550 5.1.1"},
		{"gomail send text multiline", errors.New("gomail: could not send email 1: 452-4.2.2 Mailbox full <user@real-inbox.example>"), "smtp 452 4.2.2"},
		{"gomail send text without enhanced code", errors.New("gomail: could not send email 1: 554 Transaction failed"), "smtp 554"},
		{"mime too deep", ErrMIMETooDeep, ErrMIMETooDeep.Error()},
		{"parse timeout", ErrParseTimeout, ErrParseTimeout.Error()},
		{"starttls unsupported", gomail.ErrStartTLSUnsupported, gomail.ErrStartTLSUnsupported.Error()},
		{"network error", &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}, "network error"},
		{"network timeout", &net.OpError{Op: "dial", Net: "tcp", Err: context.DeadlineExceeded}, "network timeout"},
		{"other", errors.New("cannot parse address <user@real-inbox.example>"), "unclassified error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorClass(tt.err)
			if got != tt.want {
				t.Errorf("ErrorClass() = %q, want %q", got, tt.want)
			}
			if strings.Contains(got, "@") {
				t.Errorf("ErrorClass() = %q leaks an address", got)
			}
		})
	}
}

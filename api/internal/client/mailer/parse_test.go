package mailer

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/model"
)

// nestedEmail builds a message with depth levels of multipart/mixed, each wrapping the last.
func nestedEmail(depth int) []byte {
	nested := "Content-Type: text/plain\r\n\r\nleaf body"
	for i := range depth {
		b := fmt.Sprintf("B%d", i)
		nested = fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n--%s\r\n%s\r\n--%s--\r\n", b, b, nested, b)
	}

	return []byte("From: a@example.com\r\nTo: b@example.com\r\nSubject: nested\r\n" + nested)
}

func TestMIMEDepthExceeds(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"plain text", []byte(testForwardRawEmail), false},
		{"one level", nestedEmail(1), false},
		{"at the limit", nestedEmail(maxMIMEDepth), false},
		{"past the limit", nestedEmail(maxMIMEDepth + 1), true},
		{"far past the limit", nestedEmail(5000), true},
		{"not an email", []byte("garbage"), false},
		{"no boundary", []byte("From: a@example.com\r\nContent-Type: multipart/mixed\r\n\r\nbody"), false},
		{"malformed content type", []byte("From: a@example.com\r\nContent-Type: multipart/mixed; boundary=\"\r\n\r\nbody"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mimeDepthExceeds(tt.data, maxMIMEDepth); got != tt.want {
				t.Errorf("mimeDepthExceeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMIMEDepthExceeds_DeepSiblingAfterShallowPart(t *testing.T) {
	// A deep part after a shallow sibling must still be found
	deep := string(nestedEmail(maxMIMEDepth))
	deep = deep[strings.Index(deep, "Content-Type:"):]
	data := "From: a@example.com\r\nContent-Type: multipart/mixed; boundary=TOP\r\n\r\n" +
		"--TOP\r\nContent-Type: text/plain\r\n\r\nshallow\r\n" +
		"--TOP\r\n" + deep + "\r\n" +
		"--TOP--\r\n"

	if !mimeDepthExceeds([]byte(data), maxMIMEDepth) {
		t.Error("mimeDepthExceeds() = false, want true for a deep second part")
	}
}

func TestParseEmail(t *testing.T) {
	t.Run("parses nesting at the limit", func(t *testing.T) {
		email, err := parseEmail(nestedEmail(maxMIMEDepth))
		if err != nil {
			t.Fatalf("parseEmail() error = %v", err)
		}
		if !strings.Contains(email.Text, "leaf body") {
			t.Errorf("Text = %q, want the leaf body", email.Text)
		}
	})

	t.Run("rejects deep nesting quickly", func(t *testing.T) {
		// Built before timing, as building the message is itself slow
		data := nestedEmail(5000)

		start := time.Now()
		_, err := parseEmail(data)
		if !errors.Is(err, ErrMIMETooDeep) {
			t.Fatalf("parseEmail() error = %v, want ErrMIMETooDeep", err)
		}
		if d := time.Since(start); d > time.Second {
			t.Errorf("parseEmail() took %v, want it rejected quickly", d)
		}
	})
}

func TestParseWithTimeout(t *testing.T) {
	t.Run("returns the parsed email", func(t *testing.T) {
		want := letters.Email{Text: "body"}
		got, err := parseWithTimeout(func() (letters.Email, error) { return want, nil }, time.Second)
		if err != nil || got.Text != want.Text {
			t.Errorf("parseWithTimeout() = %q, %v, want %q, nil", got.Text, err, want.Text)
		}
	})

	t.Run("returns the parse error", func(t *testing.T) {
		parseErr := errors.New("bad message")
		_, err := parseWithTimeout(func() (letters.Email, error) { return letters.Email{}, parseErr }, time.Second)
		if !errors.Is(err, parseErr) {
			t.Errorf("parseWithTimeout() error = %v, want %v", err, parseErr)
		}
	})

	t.Run("times out", func(t *testing.T) {
		release := make(chan struct{})
		defer close(release)

		_, err := parseWithTimeout(func() (letters.Email, error) {
			<-release
			return letters.Email{}, nil
		}, 10*time.Millisecond)
		if !errors.Is(err, ErrParseTimeout) {
			t.Errorf("parseWithTimeout() error = %v, want ErrParseTimeout", err)
		}
	})

	t.Run("recovers a panic", func(t *testing.T) {
		_, err := parseWithTimeout(func() (letters.Email, error) {
			panic("boom")
		}, time.Second)
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Errorf("parseWithTimeout() error = %v, want the panic as an error", err)
		}
	})
}

func TestBuildMessages_RejectDeepNesting(t *testing.T) {
	m := Mailer{}
	rcp := model.Recipient{Email: "real.recipient@domain.com"}
	alias := model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: "myalias@mailx.net"}
	data := nestedEmail(5000)

	_, _, err := m.buildForwardMessage(
		"sender@example.com", "Sender", alias.Name, rcp,
		data, "header.tmpl", map[string]any{}, model.Settings{}, alias, model.ListUnsubscribeHeaders{}, model.RuleOutcome{},
	)
	if !errors.Is(err, ErrMIMETooDeep) {
		t.Errorf("buildForwardMessage() error = %v, want ErrMIMETooDeep", err)
	}

	_, _, err = m.buildReplyMessage(alias.Name, "", rcp, data, alias, model.RuleOutcome{})
	if !errors.Is(err, ErrMIMETooDeep) {
		t.Errorf("buildReplyMessage() error = %v, want ErrMIMETooDeep", err)
	}
}

package signedurl

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func newTestSigner(now time.Time) *Signer {
	s := New("test-secret", "https://api.example.com")
	s.now = func() time.Time { return now }
	return s
}

func splitURL(t *testing.T, raw string) (string, string) {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", raw, err)
	}
	return u.Path, u.RawQuery
}

func TestNew_RequiresSecretAndBaseURL(t *testing.T) {
	if New("", "https://api.example.com") != nil {
		t.Error("New() with empty secret should return nil")
	}
	if New("secret", "") != nil {
		t.Error("New() with empty base URL should return nil")
	}
}

func TestSignVerify(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := newTestSigner(now)

	signed := s.Sign("/v1/block-email-one-click/alias-1", url.Values{"email": {"news+x@example.com"}}, time.Hour)
	if !strings.HasPrefix(signed, "https://api.example.com/v1/block-email-one-click/alias-1?") {
		t.Fatalf("Sign() = %q, want base URL + path prefix", signed)
	}

	path, query := splitURL(t, signed)
	if !s.Verify(path, query) {
		t.Fatal("Verify() = false for a freshly signed URL")
	}

	q, _ := url.ParseQuery(query)
	if got := q.Get("email"); got != "news+x@example.com" {
		t.Errorf("email param = %q, want it to survive encoding", got)
	}
}

func TestVerify_Rejects(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := newTestSigner(now)
	signed := s.Sign("/v1/block-email-one-click/alias-1", url.Values{"email": {"a@example.com"}}, time.Hour)
	path, query := splitURL(t, signed)
	q, _ := url.ParseQuery(query)

	with := func(key, value string) string {
		c := url.Values{}
		for k, v := range q {
			c[k] = append([]string(nil), v...)
		}
		if value == "" {
			c.Del(key)
		} else {
			c.Set(key, value)
		}
		return c.Encode()
	}

	tests := []struct {
		name  string
		path  string
		query string
	}{
		{"tampered email", path, with("email", "b@example.com")},
		{"added param", path, query + "&domain=example.com"},
		{"tampered alias path", "/v1/block-email-one-click/alias-2", query},
		{"tampered expires", path, with("expires", "9999999999")},
		{"missing signature", path, with("signature", "")},
		{"garbage signature", path, with("signature", "zz")},
		{"missing expires", path, with("expires", "")},
		{"malformed query", path, "%zz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if s.Verify(tt.path, tt.query) {
				t.Errorf("Verify(%q, %q) = true, want false", tt.path, tt.query)
			}
		})
	}
}

func TestVerify_Expired(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := newTestSigner(now)
	path, query := splitURL(t, s.Sign("/v1/deactivate-one-click/alias-1", nil, time.Hour))

	s.now = func() time.Time { return now.Add(time.Hour + time.Second) }
	if s.Verify(path, query) {
		t.Error("Verify() = true for an expired URL")
	}
}

func TestVerify_DifferentSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	path, query := splitURL(t, newTestSigner(now).Sign("/v1/deactivate-one-click/alias-1", nil, time.Hour))

	other := New("other-secret", "https://api.example.com")
	other.now = func() time.Time { return now }
	if other.Verify(path, query) {
		t.Error("Verify() = true for a URL signed with another secret")
	}
}
